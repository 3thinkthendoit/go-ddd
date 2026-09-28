package constant

import "fmt"

// OrderType 订单类型
type OrderType int

const (
	OrderTypeStandard OrderType = 1 // 标准订单
	OrderTypePreSale  OrderType = 2 // 预售订单
	OrderTypeTeamBuy  OrderType = 3 // 拼团订单
)

var orderTypeDesc = map[OrderType]string{
	OrderTypeStandard: "标准订单",
	OrderTypePreSale:  "预售订单",
	OrderTypeTeamBuy:  "拼团订单",
}

// Code 返回类型编码
func (t OrderType) Code() int { return int(t) }

// Desc 返回类型描述
func (t OrderType) Desc() string { return orderTypeDesc[t] }

func (t OrderType) String() string {
	if desc, ok := orderTypeDesc[t]; ok {
		return desc
	}
	return fmt.Sprintf("OrderType(%d)", int(t))
}

// OrderTypeOf 按编码查找订单类型
func OrderTypeOf(code int) (OrderType, bool) {
	t := OrderType(code)
	_, ok := orderTypeDesc[t]
	return t, ok
}
