package response

// RiskCheckResponse 风控校验结果
type RiskCheckResponse struct {
	// 批量下单
	IsBatchBuy bool
	// 错误收货信息
	IsIllegalAddress bool
	// 价格欺诈
	IsIllegalSkuPrice bool
	// 支付欺诈
	IsIllegalPay bool
	// 错误信息
	Desc string
}

// Rejected 是否存在任意一项风控命中
func (r *RiskCheckResponse) Rejected() bool {
	return r.IsBatchBuy || r.IsIllegalAddress || r.IsIllegalSkuPrice || r.IsIllegalPay
}
