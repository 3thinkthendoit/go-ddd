package response

// ShippingInfo WMS 返回的单条发货信息
type ShippingInfo struct {
	SkuCode        string
	ShippingAmount int
	ExpressCode    string
	ExpressNo      string
}

// ShippingQueryResponse WMS 发货信息查询结果
type ShippingQueryResponse struct {
	ShippingInfos []ShippingInfo
}
