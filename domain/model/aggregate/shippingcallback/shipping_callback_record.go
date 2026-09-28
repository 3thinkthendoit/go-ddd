package shippingcallback

import "think.com/go-ddd/domain/model/constant"

// ShippingCallbackRecord 发货回传记录实体（每次回传留痕，便于对账与重试）。
type ShippingCallbackRecord struct {
	orderNo         string
	externalOrderNo string
	orderSource     constant.OrderSource
	skuItems        string
	callStatus      int
	result          string
}

// NewShippingCallbackRecord 创建回传记录
func NewShippingCallbackRecord(orderNo, externalOrderNo string, orderSource constant.OrderSource,
	skuItems string, callStatus int, result string) *ShippingCallbackRecord {
	return &ShippingCallbackRecord{
		orderNo:         orderNo,
		externalOrderNo: externalOrderNo,
		orderSource:     orderSource,
		skuItems:        skuItems,
		callStatus:      callStatus,
		result:          result,
	}
}

// OrderNo OMS 订单号
func (r *ShippingCallbackRecord) OrderNo() string { return r.orderNo }

// ExternalOrderNo 外部订单号
func (r *ShippingCallbackRecord) ExternalOrderNo() string { return r.externalOrderNo }

// OrderSource 订单来源
func (r *ShippingCallbackRecord) OrderSource() constant.OrderSource { return r.orderSource }

// SkuItems 回传的 sku 明细快照（JSON）
func (r *ShippingCallbackRecord) SkuItems() string { return r.skuItems }

// CallStatus 回传状态
func (r *ShippingCallbackRecord) CallStatus() int { return r.callStatus }

// Result 回传返回内容
func (r *ShippingCallbackRecord) Result() string { return r.result }
