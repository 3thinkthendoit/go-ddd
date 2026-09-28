package constant

import "fmt"

// OrderStatus 订单状态
type OrderStatus int

const (
	OrderStatusUnPay     OrderStatus = -1 // 未支付
	OrderStatusPayed     OrderStatus = 1  // 已支付
	OrderStatusCancel    OrderStatus = 2  // 取消
	OrderStatusFulfilled OrderStatus = 3  // 已发货
	OrderStatusHangUp    OrderStatus = 4  // 挂起
)

var orderStatusDesc = map[OrderStatus]string{
	OrderStatusUnPay:     "未支付",
	OrderStatusPayed:     "已支付",
	OrderStatusCancel:    "取消",
	OrderStatusFulfilled: "已发货",
	OrderStatusHangUp:    "挂起",
}

// Code 返回状态编码
func (s OrderStatus) Code() int { return int(s) }

// Desc 返回状态描述
func (s OrderStatus) Desc() string { return orderStatusDesc[s] }

func (s OrderStatus) String() string {
	if desc, ok := orderStatusDesc[s]; ok {
		return desc
	}
	return fmt.Sprintf("OrderStatus(%d)", int(s))
}

// OrderStatusOf 按编码查找订单状态，不存在时返回 false
func OrderStatusOf(code int) (OrderStatus, bool) {
	s := OrderStatus(code)
	_, ok := orderStatusDesc[s]
	return s, ok
}
