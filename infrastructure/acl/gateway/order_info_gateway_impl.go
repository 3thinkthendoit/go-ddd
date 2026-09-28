package gateway

import (
	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/domain/model/constant"
	domainpl "think.com/go-ddd/domain/pl"
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
	"think.com/go-ddd/domain/port/gateway"
	aclpl "think.com/go-ddd/infrastructure/acl/pl"
	"think.com/go-ddd/infrastructure/core/orm"
	"think.com/go-ddd/infrastructure/core/orm/po"
)

func init() {
	gs.Object(new(OrderInfoGatewayImpl)).Export((*gateway.OrderInfoGateway)(nil))
}

// OrderInfoGatewayImpl 订单信息南向网关实现。
//
// 查询订单主表 + sku item 表并转换成对外统一订单视图；真实项目可换成
// 查询订单中心 / 其他域服务。
type OrderInfoGatewayImpl struct {
	store *orm.Store `autowire:""`
}

// Query 按订单号或「外部订单号 + 来源」查询订单
func (g *OrderInfoGatewayImpl) Query(req *request.OrderQueryRequest) (*response.OrderQueryResponse, error) {
	orders := g.store.SelectOrderBaseInfos(func(v *po.OrderBaseInfo) bool {
		if req.OrderNo != "" {
			return v.OrderNo == req.OrderNo
		}
		if req.ExternalOrderNo != "" && req.OrderSource != constant.OrderSourceUnKnow {
			return v.ExternalOrderNo == req.ExternalOrderNo && v.OrderSource == req.OrderSource.Code()
		}
		return false
	})

	result := make([]*domainpl.OrderInfo, 0, len(orders))
	for _, order := range orders {
		skuItemInfos := g.store.SelectOrderSkuItemInfos(func(v *po.OrderSkuItemInfo) bool {
			return v.OrderNo == order.OrderNo
		})
		result = append(result, aclpl.ToOrderInfo(order, skuItemInfos))
	}
	return &response.OrderQueryResponse{Orders: result}, nil
}
