package constant

import "fmt"

// PayType 支付方式
type PayType int

const (
	PayTypeBalance  PayType = 1 // 余额支付
	PayTypeAlipay   PayType = 2 // 支付宝
	PayTypeWxPay    PayType = 3 // 微信支付
	PayTypeUnionPay PayType = 4 // 银联支付
)

var payTypeDesc = map[PayType]string{
	PayTypeBalance:  "余额支付",
	PayTypeAlipay:   "支付宝",
	PayTypeWxPay:    "微信支付",
	PayTypeUnionPay: "银联支付",
}

// Code 返回支付方式编码
func (p PayType) Code() int { return int(p) }

// Desc 返回支付方式描述
func (p PayType) Desc() string { return payTypeDesc[p] }

func (p PayType) String() string {
	if desc, ok := payTypeDesc[p]; ok {
		return desc
	}
	return fmt.Sprintf("PayType(%d)", int(p))
}

// PayTypeOf 按编码查找支付方式
func PayTypeOf(code int) (PayType, bool) {
	p := PayType(code)
	_, ok := payTypeDesc[p]
	return p, ok
}
