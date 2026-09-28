package response

import "think.com/go-ddd/domain/pl"

// OrderQueryResponse 订单查询结果
type OrderQueryResponse struct {
	Orders []*pl.OrderInfo
}
