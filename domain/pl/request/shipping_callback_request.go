package request

import "think.com/go-ddd/domain/model/constant"

// ShippingCallbackInfo 单条回传明细
type ShippingCallbackInfo struct {
	ExternalSkuCode string
	SkuAmount       int
	ShippingAmount  int
	ExpressCode     string
	ExpressNo       string
}

// ShippingCallbackRequest 向上游平台回传发货信息请求
type ShippingCallbackRequest struct {
	ExternalOrderNo string
	OrderSource     constant.OrderSource
	ShippingInfos   []ShippingCallbackInfo
}
