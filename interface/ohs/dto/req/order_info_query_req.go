package req

// OrderInfoQueryReq 订单查询请求
type OrderInfoQueryReq struct {
	OrderNo         string `json:"orderNo"`
	ExternalOrderNo string `json:"externalOrderNo"`
	// 订单来源编码，0 表示不限
	OrderSource int `json:"orderSource"`
}
