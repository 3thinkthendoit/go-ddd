package po

// OrderSkuInfo 订单 sku 表（order_sku_info）
type OrderSkuInfo struct {
	Id              int64
	OrderNo         string
	ExternalSkuId   string
	ExternalSkuCode string
	SkuCode         string
	SkuName         string
	SkuPayPrice     int64
	SkuBuyAmount    int
	SkuType         int
	SkuCategory     int
}
