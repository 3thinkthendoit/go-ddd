package gateway

import (
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
)

// OrderInfoGateway 订单信息南向网关（查询其他域的订单数据 / 本域历史数据）
type OrderInfoGateway interface {
	// Query 按订单号 / 外部订单号 + 来源查询订单
	Query(request *request.OrderQueryRequest) (*response.OrderQueryResponse, error)
}
