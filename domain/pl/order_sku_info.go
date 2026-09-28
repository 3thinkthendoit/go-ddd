package pl

// OrderSkuInfo 外部平台传入的 sku 下单信息
type OrderSkuInfo struct {
	ExternalSkuId   string
	ExternalSkuCode string
	SkuName         string
	SkuPayPrice     int64
	SkuBuyAmount    int
}
