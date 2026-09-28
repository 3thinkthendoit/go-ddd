package response

import "think.com/go-ddd/domain/pl"

// SkuInfoQueryResponse 商品信息查询结果
type SkuInfoQueryResponse struct {
	SkuInfos []*pl.SkuFullInfo
}
