package constant

import "fmt"

// OrderSkuStatus sku 级别的发货/售后状态
type OrderSkuStatus int

const (
	OrderSkuStatusNotShipped  OrderSkuStatus = 1 // 未发货
	OrderSkuStatusPartShipped OrderSkuStatus = 2 // 部分发货
	OrderSkuStatusShipped     OrderSkuStatus = 3 // 全部发货
	OrderSkuStatusReturn      OrderSkuStatus = 4 // 全部退货
	OrderSkuStatusPartReturn  OrderSkuStatus = 5 // 部分退货
	OrderSkuStatusCheckedIn   OrderSkuStatus = 6 // 已签收
)

var orderSkuStatusDesc = map[OrderSkuStatus]string{
	OrderSkuStatusNotShipped:  "未发货",
	OrderSkuStatusPartShipped: "部分发货",
	OrderSkuStatusShipped:     "全部发货",
	OrderSkuStatusReturn:      "全部退货",
	OrderSkuStatusPartReturn:  "部分退货",
	OrderSkuStatusCheckedIn:   "已签收",
}

// Code 返回状态编码
func (s OrderSkuStatus) Code() int { return int(s) }

// Desc 返回状态描述
func (s OrderSkuStatus) Desc() string { return orderSkuStatusDesc[s] }

func (s OrderSkuStatus) String() string {
	if desc, ok := orderSkuStatusDesc[s]; ok {
		return desc
	}
	return fmt.Sprintf("OrderSkuStatus(%d)", int(s))
}

// OrderSkuStatusOf 按编码查找 sku 状态
func OrderSkuStatusOf(code int) (OrderSkuStatus, bool) {
	s := OrderSkuStatus(code)
	_, ok := orderSkuStatusDesc[s]
	return s, ok
}
