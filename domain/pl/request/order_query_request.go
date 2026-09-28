package request

import "think.com/go-ddd/domain/model/constant"

// OrderQueryRequest 订单查询请求
type OrderQueryRequest struct {
	ExternalOrderNo string
	OrderSource     constant.OrderSource
	OrderNo         string
}
