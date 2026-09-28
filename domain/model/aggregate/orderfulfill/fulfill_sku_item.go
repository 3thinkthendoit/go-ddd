package orderfulfill

import (
	"fmt"

	"think.com/go-ddd/domain/pl"
)

// FulfillSkuItem 履约 sku item 实体。
//
// 持有下单数量、已发货数量、发货优先级，以及分仓结果（仓库 → 发货数量）。
type FulfillSkuItem struct {
	skuCode        string
	skuName        string
	skuAmount      int
	shippingAmount int
	priority       int
	dispatchInfo   map[string]int
}

// NewFulfillSkuItem 由 sku item 信息创建
func NewFulfillSkuItem(skuItemInfo *pl.SkuItemInfo) *FulfillSkuItem {
	if skuItemInfo == nil || skuItemInfo.SkuFullInfo == nil {
		return nil
	}
	return &FulfillSkuItem{
		skuCode:        skuItemInfo.SkuFullInfo.SkuCode,
		skuName:        skuItemInfo.SkuFullInfo.SkuName,
		skuAmount:      skuItemInfo.SkuAmount,
		shippingAmount: skuItemInfo.ShippingAmount,
		priority:       skuItemInfo.Priority,
		dispatchInfo:   map[string]int{},
	}
}

// Dispatch 分仓（领域方法）：把指定数量的 sku 分配到指定仓库
func (f *FulfillSkuItem) Dispatch(warehouseCode string, skuAmount int) {
	f.dispatchInfo[warehouseCode] += skuAmount
}

// PendingAmount 待发货数量
func (f *FulfillSkuItem) PendingAmount() int {
	pending := f.skuAmount - f.shippingAmount
	if pending < 0 {
		return 0
	}
	return pending
}

// NeedShip 是否需要发货：优先级为 0（虚拟商品）或已发完则无需发货
func (f *FulfillSkuItem) NeedShip() bool {
	return f.priority > 0 && f.PendingAmount() > 0
}

// SkuCode 内部 skuCode
func (f *FulfillSkuItem) SkuCode() string { return f.skuCode }

// SkuName 商品名称
func (f *FulfillSkuItem) SkuName() string { return f.skuName }

// SkuAmount 下单数量
func (f *FulfillSkuItem) SkuAmount() int { return f.skuAmount }

// ShippingAmount 已发货数量
func (f *FulfillSkuItem) ShippingAmount() int { return f.shippingAmount }

// Priority 发货优先级
func (f *FulfillSkuItem) Priority() int { return f.priority }

// DispatchInfo 分仓结果：key 为仓库编码，value 为发货数量
func (f *FulfillSkuItem) DispatchInfo() map[string]int { return f.dispatchInfo }

func (f *FulfillSkuItem) String() string {
	return fmt.Sprintf("FulfillSkuItem{skuCode=%s, skuAmount=%d, priority=%d, dispatch=%v}",
		f.skuCode, f.skuAmount, f.priority, f.dispatchInfo)
}
