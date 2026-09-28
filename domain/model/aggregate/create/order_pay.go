package create

import (
	"fmt"

	"think.com/go-ddd/domain/model/constant"
)

// OrderPay 支付信息（值对象）
type OrderPay struct {
	currency       constant.Currency
	payType        constant.PayType
	payAmount      int64
	discountAmount int64
	feeAmountMap   map[constant.FeeType]int64
}

// NewOrderPay 创建支付信息
func NewOrderPay(currency constant.Currency, payType constant.PayType, payAmount,
	discountAmount int64, feeAmountMap map[constant.FeeType]int64) *OrderPay {
	if feeAmountMap == nil {
		feeAmountMap = map[constant.FeeType]int64{}
	}
	return &OrderPay{
		currency:       currency,
		payType:        payType,
		payAmount:      payAmount,
		discountAmount: discountAmount,
		feeAmountMap:   feeAmountMap,
	}
}

// Currency 币种
func (p *OrderPay) Currency() constant.Currency { return p.currency }

// PayType 支付方式
func (p *OrderPay) PayType() constant.PayType { return p.payType }

// PayAmount 实际支付金额（分）
func (p *OrderPay) PayAmount() int64 { return p.payAmount }

// DiscountAmount 优惠金额（分）
func (p *OrderPay) DiscountAmount() int64 { return p.discountAmount }

// FeeAmountMap 附加费用 / 优惠费用
func (p *OrderPay) FeeAmountMap() map[constant.FeeType]int64 { return p.feeAmountMap }

func (p *OrderPay) String() string {
	return fmt.Sprintf("OrderPay{payAmount=%d, discountAmount=%d, payType=%s}",
		p.payAmount, p.discountAmount, p.payType)
}
