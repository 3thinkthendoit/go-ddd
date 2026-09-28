package response

import "think.com/go-ddd/domain/pl"

// ShippingCallbackResponse 上游平台发货回传结果
type ShippingCallbackResponse struct {
	CallbackResult *pl.ShippingCallbackResult
}
