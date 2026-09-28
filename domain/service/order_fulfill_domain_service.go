package service

import (
	"think.com/go-ddd/domain/model/aggregate/orderfulfill"
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/port/gateway"
)

// OrderFulfillDomainService 订单履约领域服务。
//
// 负责从仓库域取回分仓所需的外部数据，再交给聚合根完成分仓与拆单。
type OrderFulfillDomainService struct {
	warehouseQueryGateway gateway.WarehouseQueryGateway `autowire:""`
}

// InitBaseInfo 回填仓库信息与指定仓映射
func (s *OrderFulfillDomainService) InitBaseInfo(aggregate *orderfulfill.OrderFulfillAggregate) error {
	skuCodes := make([]string, 0, len(aggregate.FulfillSkuItems()))
	for _, item := range aggregate.FulfillSkuItems() {
		skuCodes = append(skuCodes, item.SkuCode())
	}
	resp, err := s.warehouseQueryGateway.Query(&request.WarehouseQueryRequest{
		SkuCodes:  skuCodes,
		StoreCode: aggregate.StoreCode(),
	})
	if err != nil {
		return err
	}
	if resp == nil {
		return aggregate.InitBaseInfo(nil, nil, nil)
	}
	return aggregate.InitBaseInfo(resp.WarehouseInfos, resp.SkuMappingWarehouseMap, resp.StoreMappingWarehouseMap)
}
