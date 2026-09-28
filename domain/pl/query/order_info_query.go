package query

// OrderInfoQuery 订单信息查询条件（查询侧，可绕过领域层直接走南向网关）
type OrderInfoQuery struct {
	OrderNo         string
	ExternalOrderNo string
	// 订单来源编码，0 表示不限
	OrderSource int
}
