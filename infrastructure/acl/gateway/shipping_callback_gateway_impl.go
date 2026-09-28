package gateway

import (
	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/domain/model/constant"
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
	"think.com/go-ddd/domain/port/gateway"
	"think.com/go-ddd/infrastructure/acl/api/douyin"
	"think.com/go-ddd/infrastructure/acl/api/mijia"
	"think.com/go-ddd/infrastructure/acl/api/pdd"
	"think.com/go-ddd/infrastructure/acl/api/taobao"
)

func init() {
	gs.Object(new(ShippingCallbackGatewayImpl)).Export((*gateway.ShippingCallbackGateway)(nil))
}

// ShippingCallbackGatewayImpl 发货回传南向网关实现。
//
// 按订单来源路由到对应平台的协议客户端，新增平台时只需要在这里加一个分支。
type ShippingCallbackGatewayImpl struct {
	taoBaoClient *taobao.Client `autowire:""`
	duoYinClient *douyin.Client `autowire:""`
	pddClient    *pdd.Client    `autowire:""`
	miJiaClient  *mijia.Client  `autowire:""`
}

// Callback 把发货信息回传给订单来源平台
func (g *ShippingCallbackGatewayImpl) Callback(req *request.ShippingCallbackRequest) (*response.ShippingCallbackResponse, error) {
	switch req.OrderSource {
	case constant.OrderSourceTaoBao:
		return g.taoBaoClient.ShippingCallback(req)
	case constant.OrderSourceDouYin:
		return g.duoYinClient.ShippingCallback(req)
	case constant.OrderSourcePdd:
		return g.pddClient.ShippingCallback(req)
	case constant.OrderSourceMiJia:
		return g.miJiaClient.ShippingCallback(req)
	default:
		return &response.ShippingCallbackResponse{}, nil
	}
}
