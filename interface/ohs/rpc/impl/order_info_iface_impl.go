// Package impl 微服务接口实现。
package impl

import (
	"github.com/go-spring/spring-core/gs"
	"github.com/go-spring/spring-core/web"

	"think.com/go-ddd/domain/pl/query"
	"think.com/go-ddd/interface/local"
	"think.com/go-ddd/interface/ohs/dto/req"
	"think.com/go-ddd/interface/ohs/dto/resp"
	"think.com/go-ddd/interface/ohs/httpx"
	"think.com/go-ddd/interface/ohs/rpc"
)

func init() {
	gs.Object(new(OrderInfoIfaceImpl)).Export((*rpc.OrderInfoIface)(nil)).Init(func(i *OrderInfoIfaceImpl) {
		gs.PostMapping("/orderInfoIface/query", i.Query)
	})
}

// OrderInfoIfaceImpl 订单微服务接口实现
type OrderInfoIfaceImpl struct {
	orderLocalService *local.OrderLocalService `autowire:""`
}

// Query 订单查询：DTO → 领域查询对象 → 领域订单视图 → DTO
func (i *OrderInfoIfaceImpl) Query(ctx web.Context) {
	var reqDto req.OrderInfoQueryReq
	if err := httpx.BindJSON(ctx, &reqDto); err != nil {
		httpx.WriteFailure(ctx, err)
		return
	}

	orders, err := i.orderLocalService.Query(&query.OrderInfoQuery{
		OrderNo:         reqDto.OrderNo,
		ExternalOrderNo: reqDto.ExternalOrderNo,
		OrderSource:     reqDto.OrderSource,
	})
	if err != nil {
		httpx.WriteFailure(ctx, err)
		return
	}
	httpx.WriteSuccess(ctx, resp.NewOrderInfoQueryResps(orders))
}
