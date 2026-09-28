package request

// RiskCheckSkuInfo 风控校验用到的 sku 快照
type RiskCheckSkuInfo struct {
	SkuCode   string
	SkuAmount int
	SkuPrice  int64
}

// RiskCheckRequest 风控校验请求
type RiskCheckRequest struct {
	// 下单用户
	Username string
	// 下单地址
	Address string
	// 下单手机号
	PhoneNo string
	// 下单商品
	SkuInfos []RiskCheckSkuInfo
}
