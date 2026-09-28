package po

import "time"

// ShippingCallbackRecord 发货回传记录表（shipping_callback_record）
type ShippingCallbackRecord struct {
	Id              int64
	OrderNo         string
	ExternalOrderNo string
	OrderSource     int
	SkuItems        string
	CallStatus      int
	Result          string
	CreateTime      time.Time
}
