package create

import (
	"fmt"

	"think.com/go-ddd/domain/model/constant"
)

// OrderSkuItem sku item 级下单实体。
//
// sku 可能被拆成多个 item（组合商品拆子商品、金额按权重分摊），
// 后续拆单、发货、售后都以 item 为单位流转。
type OrderSkuItem struct {
	skuCode        string
	skuAmount      int
	stockAmount    int
	payPrice       int64
	priority       int
	feeAmountInfos map[constant.FeeType]int64
}

// NewOrderSkuItem 创建 sku item
func NewOrderSkuItem(skuCode string, skuAmount int, payPrice int64,
	feeAmountInfos map[constant.FeeType]int64) *OrderSkuItem {
	if feeAmountInfos == nil {
		feeAmountInfos = map[constant.FeeType]int64{}
	}
	return &OrderSkuItem{
		skuCode:        skuCode,
		skuAmount:      skuAmount,
		payPrice:       payPrice,
		feeAmountInfos: feeAmountInfos,
	}
}

// PriorityProcessing 设置发货优先级（领域行为）
func (i *OrderSkuItem) PriorityProcessing(priority int) {
	i.priority = priority
}

// SetStockAmount 设置占用库存数量
func (i *OrderSkuItem) SetStockAmount(stockAmount int) {
	i.stockAmount = stockAmount
}

// SkuCode 内部 skuCode
func (i *OrderSkuItem) SkuCode() string { return i.skuCode }

// SkuAmount 下单数量
func (i *OrderSkuItem) SkuAmount() int { return i.skuAmount }

// StockAmount 占用库存数量
func (i *OrderSkuItem) StockAmount() int { return i.stockAmount }

// PayPrice 实付金额（分）
func (i *OrderSkuItem) PayPrice() int64 { return i.payPrice }

// Priority 发货优先级，0 表示无需发货，数字越大优先级越高
func (i *OrderSkuItem) Priority() int { return i.priority }

// FeeAmountInfos 分摊到该 item 的附加费用
func (i *OrderSkuItem) FeeAmountInfos() map[constant.FeeType]int64 { return i.feeAmountInfos }

func (i *OrderSkuItem) String() string {
	return fmt.Sprintf("OrderSkuItem{skuCode=%s, amount=%d, payPrice=%d, priority=%d}",
		i.skuCode, i.skuAmount, i.payPrice, i.priority)
}
