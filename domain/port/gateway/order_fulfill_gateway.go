package gateway

import (
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
)

// OrderFulfillGateway 履约（WMS）南向网关
type OrderFulfillGateway interface {
	// Fulfill 推送发货单给 WMS
	Fulfill(request *request.OrderFulfillRequest) (*response.OrderFulfillResponse, error)
	// Query 查询 WMS 发货信息
	Query(request *request.ShippingQueryRequest) (*response.ShippingQueryResponse, error)
}
