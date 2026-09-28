package gateway

import (
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
)

// ShippingCallbackGateway 上游平台发货回传南向网关
type ShippingCallbackGateway interface {
	// Callback 把发货信息回传给订单来源平台
	Callback(request *request.ShippingCallbackRequest) (*response.ShippingCallbackResponse, error)
}
