package request

// SkuInfoQueryRequest 商品信息查询请求
type SkuInfoQueryRequest struct {
	SkuIds           []string
	ExternalSkuIds   []string
	SkuCodes         []string
	ExternalSkuCodes []string
}
