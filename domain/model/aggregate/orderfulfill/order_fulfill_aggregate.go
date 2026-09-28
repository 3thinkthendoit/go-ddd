package orderfulfill

import (
	"errors"
	"fmt"
	"sort"

	"think.com/go-ddd/domain/common/assert"
	"think.com/go-ddd/domain/model/dp"
	"think.com/go-ddd/domain/pl"
)

// OrderFulfillAggregate 订单履约聚合根。
//
// 负责「分仓 → 拆单」：先按 店铺指定仓 > sku 指定仓 > 最优仓计算 三级策略确定
// 每个 sku 的发货仓库，再按仓库维度把订单拆成若干发货单。
type OrderFulfillAggregate struct {
	orderNo   string
	storeCode string

	// 订单下单 sku item
	fulfillSkuItems []*FulfillSkuItem
	// sku 仓库库存信息
	fulfillWarehouses []*FulfillWarehouse
	// sku → 仓库 映射关系
	skuMappingWarehouseMap map[string]string
	// 店铺 → 仓库 映射关系
	storeMappingWarehouseMap map[string]string
	// 拆单结果：key 为仓库编码
	splitOrders map[string][]*OrderSplitResult
}

// NewOrderFulfillAggregate 由订单号 + sku item 信息还原履约聚合根
func NewOrderFulfillAggregate(orderNo, storeCode string, skuItemInfos []*pl.SkuItemInfo) (*OrderFulfillAggregate, error) {
	if err := assert.NotBlank(orderNo, "orderNo is blank"); err != nil {
		return nil, err
	}
	if err := assert.NotBlank(storeCode, "storeCode is blank"); err != nil {
		return nil, err
	}
	if err := assert.NotEmpty(len(skuItemInfos), "skuItemInfos is empty"); err != nil {
		return nil, err
	}

	aggregate := &OrderFulfillAggregate{
		orderNo:                  orderNo,
		storeCode:                storeCode,
		fulfillSkuItems:          make([]*FulfillSkuItem, 0, len(skuItemInfos)),
		fulfillWarehouses:        make([]*FulfillWarehouse, 0),
		skuMappingWarehouseMap:   map[string]string{},
		storeMappingWarehouseMap: map[string]string{},
		splitOrders:              map[string][]*OrderSplitResult{},
	}
	for _, skuItemInfo := range skuItemInfos {
		item := NewFulfillSkuItem(skuItemInfo)
		if item == nil {
			return nil, errors.New("skuItemInfo is invalid")
		}
		aggregate.fulfillSkuItems = append(aggregate.fulfillSkuItems, item)
	}
	return aggregate, nil
}

// InitBaseInfo 由领域服务调用南向网关后回填仓库、映射关系等基础信息
func (a *OrderFulfillAggregate) InitBaseInfo(warehouseInfos []*pl.WarehouseInfo,
	skuMappingWarehouseMap, storeMappingWarehouseMap map[string]string) error {
	if skuMappingWarehouseMap != nil {
		a.skuMappingWarehouseMap = skuMappingWarehouseMap
	}
	if storeMappingWarehouseMap != nil {
		a.storeMappingWarehouseMap = storeMappingWarehouseMap
	}
	warehouses := make([]*FulfillWarehouse, 0, len(warehouseInfos))
	for _, info := range warehouseInfos {
		warehouse := NewFulfillWarehouse(info)
		if warehouse == nil {
			return fmt.Errorf("warehouseInfo is invalid, warehouseCode=%v", info)
		}
		warehouses = append(warehouses, warehouse)
	}
	a.fulfillWarehouses = warehouses
	return nil
}

// Check 业务检查
func (a *OrderFulfillAggregate) Check() error {
	return assert.NotEmpty(len(a.fulfillSkuItems), "fulfillSkuItems is empty")
}

// Dispatch 分仓：逐个 sku 计算发货仓库
func (a *OrderFulfillAggregate) Dispatch() error {
	for _, skuItem := range a.fulfillSkuItems {
		if err := a.doDispatch(skuItem); err != nil {
			return err
		}
	}
	return nil
}

// doDispatch 三级分仓策略
func (a *OrderFulfillAggregate) doDispatch(skuItem *FulfillSkuItem) error {
	if !skuItem.NeedShip() {
		// 虚拟商品 / 已发完，无需分仓
		return nil
	}
	// 1. 店铺指定仓库发货
	if warehouseCode, ok := a.storeMappingWarehouseMap[a.storeCode]; ok && warehouseCode != "" {
		skuItem.Dispatch(warehouseCode, skuItem.PendingAmount())
		return nil
	}
	// 2. sku 指定了仓库发货
	if warehouseCode, ok := a.skuMappingWarehouseMap[skuItem.SkuCode()]; ok && warehouseCode != "" {
		skuItem.Dispatch(warehouseCode, skuItem.PendingAmount())
		return nil
	}
	// 3. 地理位置就近 + 库存余量 + 仓库处理能力 权重计算
	return a.dispatchByOptimalWarehouse(skuItem)
}

// dispatchByOptimalWarehouse 最优仓库选择。
//
// 先过滤有库存的仓库，按 Score（处理能力 - 距离）降序贪心分配，
// 一个仓库不够时继续用下一个仓库，全部仓库都不够则报错。
func (a *OrderFulfillAggregate) dispatchByOptimalWarehouse(skuItem *FulfillSkuItem) error {
	candidates := make([]*FulfillWarehouse, 0, len(a.fulfillWarehouses))
	for _, warehouse := range a.fulfillWarehouses {
		if warehouse.InventoryOf(skuItem.SkuCode()) > 0 {
			candidates = append(candidates, warehouse)
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Score() > candidates[j].Score()
	})

	remain := skuItem.PendingAmount()
	for _, warehouse := range candidates {
		if remain <= 0 {
			break
		}
		available := warehouse.InventoryOf(skuItem.SkuCode())
		take := remain
		if int64(take) > available {
			take = int(available)
		}
		skuItem.Dispatch(warehouse.WarehouseCode(), take)
		remain -= take
	}
	if remain > 0 {
		return fmt.Errorf("skuCode=%s 库存不足, 待发货=%d, 缺口=%d",
			skuItem.SkuCode(), skuItem.PendingAmount(), remain)
	}
	return nil
}

// Split 根据分仓结果拆单。
//
// 同一仓库的多个 sku 归到同一个发货单（父单号相同），不同仓库生成不同的父单号。
func (a *OrderFulfillAggregate) Split() error {
	for _, skuItem := range a.fulfillSkuItems {
		warehouseCodes := make([]string, 0, len(skuItem.DispatchInfo()))
		for warehouseCode := range skuItem.DispatchInfo() {
			warehouseCodes = append(warehouseCodes, warehouseCode)
		}
		sort.Strings(warehouseCodes)

		for _, warehouseCode := range warehouseCodes {
			skuAmount := skuItem.DispatchInfo()[warehouseCode]
			splitResults, ok := a.splitOrders[warehouseCode]
			var parentOrderNo string
			if !ok || len(splitResults) == 0 {
				parentOrderNo = dp.NewParentOrderNo()
				splitResults = make([]*OrderSplitResult, 0, 1)
			} else {
				// 复用同仓库已生成的父单号，保证一个仓库只对应一个发货单
				parentOrderNo = splitResults[0].ParentOrderNo()
			}
			splitResults = append(splitResults, NewOrderSplitResult(
				a.orderNo, parentOrderNo, skuItem.SkuCode(), skuAmount, warehouseCode))
			a.splitOrders[warehouseCode] = splitResults
		}
	}
	return nil
}

// OrderNo 订单号
func (a *OrderFulfillAggregate) OrderNo() string { return a.orderNo }

// StoreCode 店铺编码
func (a *OrderFulfillAggregate) StoreCode() string { return a.storeCode }

// FulfillSkuItems 履约 sku item
func (a *OrderFulfillAggregate) FulfillSkuItems() []*FulfillSkuItem { return a.fulfillSkuItems }

// FulfillWarehouses 参与分仓的仓库
func (a *OrderFulfillAggregate) FulfillWarehouses() []*FulfillWarehouse { return a.fulfillWarehouses }

// SplitOrders 拆单结果：key 为仓库编码
func (a *OrderFulfillAggregate) SplitOrders() map[string][]*OrderSplitResult { return a.splitOrders }
