package gateway

import (
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
)

// WarehouseQueryGateway 仓库域南向网关。
//
// 履约分仓需要仓库库存、店铺/sku 指定仓等外域数据，统一从这里获取，
// 让 OrderFulfillAggregate 只关心分仓算法本身。
type WarehouseQueryGateway interface {
	// Query 查询参与分仓的仓库与指定仓映射
	Query(request *request.WarehouseQueryRequest) (*response.WarehouseQueryResponse, error)
}
