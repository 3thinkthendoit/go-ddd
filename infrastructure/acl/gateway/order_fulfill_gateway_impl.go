package gateway

import (
	"log"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
	"think.com/go-ddd/domain/port/gateway"
)

func init() {
	gs.Object(new(OrderFulfillGatewayImpl)).Export((*gateway.OrderFulfillGateway)(nil))
}

// OrderFulfillGatewayImpl 履约（WMS）南向网关实现。
type OrderFulfillGatewayImpl struct {
	// 是否使用内置模拟（真实环境关闭，改为调用 WMS 接口）
	MockEnabled bool `value:"${oms.wms.mock:=true}"`
	// WMS 地址
	Endpoint string `value:"${oms.wms.endpoint:=http://wms/api/fulfill}"`
}

// Fulfill 推送发货单给 WMS
func (g *OrderFulfillGatewayImpl) Fulfill(req *request.OrderFulfillRequest) (*response.OrderFulfillResponse, error) {
	log.Printf("[wms] 推送发货单 omsOrderNo=%s, warehouse=%s, skus=%v",
		req.OmsOrderNo, req.WarehouseCode, req.Skus)
	if !g.MockEnabled {
		// 真实实现：调用 WMS 接口，按返回结果判断是否发货成功
		return &response.OrderFulfillResponse{IsFulfill: false, Msg: "WMS 接口未接入"}, nil
	}
	return &response.OrderFulfillResponse{IsFulfill: true, Msg: "mock fulfill success"}, nil
}

// Query 查询 WMS 发货信息
func (g *OrderFulfillGatewayImpl) Query(req *request.ShippingQueryRequest) (*response.ShippingQueryResponse, error) {
	// 发货数量已随 WMS 回传命令传入，这里返回空结果
	return &response.ShippingQueryResponse{}, nil
}
