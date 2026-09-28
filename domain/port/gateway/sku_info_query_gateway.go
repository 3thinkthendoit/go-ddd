package gateway

import (
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
)

// SkuInfoQueryGateway 商品信息南向网关
type SkuInfoQueryGateway interface {
	// Query 批量查询商品完整信息
	Query(request *request.SkuInfoQueryRequest) (*response.SkuInfoQueryResponse, error)
}
