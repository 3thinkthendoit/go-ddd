package orderfulfill

import "fmt"

// OrderSplitResult 拆单结果：一个发货单（父单号）下的某个 sku 在某仓库的发货明细。
type OrderSplitResult struct {
	parentOrderNo string
	orderNo       string
	skuCode       string
	skuAmount     int
	warehouseCode string
}

// NewOrderSplitResult 创建拆单结果
func NewOrderSplitResult(orderNo, parentOrderNo, skuCode string, skuAmount int, warehouseCode string) *OrderSplitResult {
	return &OrderSplitResult{
		orderNo:       orderNo,
		parentOrderNo: parentOrderNo,
		skuCode:       skuCode,
		skuAmount:     skuAmount,
		warehouseCode: warehouseCode,
	}
}

// ParentOrderNo 父单号（发货单号），唯一
func (r *OrderSplitResult) ParentOrderNo() string { return r.parentOrderNo }

// OrderNo 原子单号
func (r *OrderSplitResult) OrderNo() string { return r.orderNo }

// SkuCode 内部 skuCode
func (r *OrderSplitResult) SkuCode() string { return r.skuCode }

// SkuAmount 拆单数量
func (r *OrderSplitResult) SkuAmount() int { return r.skuAmount }

// WarehouseCode 分仓编码
func (r *OrderSplitResult) WarehouseCode() string { return r.warehouseCode }

func (r *OrderSplitResult) String() string {
	return fmt.Sprintf("OrderSplitResult{parentOrderNo=%s, orderNo=%s, skuCode=%s, skuAmount=%d, warehouseCode=%s}",
		r.parentOrderNo, r.orderNo, r.skuCode, r.skuAmount, r.warehouseCode)
}
