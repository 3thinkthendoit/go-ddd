package request

// ShippingQueryRequest WMS 发货信息查询请求
type ShippingQueryRequest struct {
	// 发货单号（拆单后的父单号）
	OmsOrderNos []string
}
