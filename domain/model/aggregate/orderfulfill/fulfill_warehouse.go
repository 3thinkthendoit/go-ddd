package orderfulfill

import (
	"think.com/go-ddd/domain/pl"
)

// FulfillWarehouse 参与分仓计算的仓库实体。
type FulfillWarehouse struct {
	warehouseCode string
	areaCode      string
	distance      float64
	shippingLevel int
	inventoryMap  map[string]int64
}

// NewFulfillWarehouse 由仓库信息创建
func NewFulfillWarehouse(info *pl.WarehouseInfo) *FulfillWarehouse {
	if info == nil {
		return nil
	}
	return &FulfillWarehouse{
		warehouseCode: info.WarehouseCode,
		areaCode:      info.AreaCode,
		distance:      info.Distance,
		shippingLevel: info.ShippingLevel,
		inventoryMap:  info.InventoryMap,
	}
}

// WarehouseCode 仓库编码
func (w *FulfillWarehouse) WarehouseCode() string { return w.warehouseCode }

// AreaCode 区域编码
func (w *FulfillWarehouse) AreaCode() string { return w.areaCode }

// Distance 与收货地址的距离（km）
func (w *FulfillWarehouse) Distance() float64 { return w.distance }

// ShippingLevel 仓库发货处理水平
func (w *FulfillWarehouse) ShippingLevel() int { return w.shippingLevel }

// InventoryOf 查询某个 sku 在该仓库的库存余量
func (w *FulfillWarehouse) InventoryOf(skuCode string) int64 {
	if w.inventoryMap == nil {
		return 0
	}
	return w.inventoryMap[skuCode]
}

// Score 最优仓权重：处理能力越高、距离越近得分越高。
func (w *FulfillWarehouse) Score() float64 {
	return float64(w.shippingLevel) - w.distance
}
