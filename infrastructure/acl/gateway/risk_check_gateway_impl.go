package gateway

import (
	"strings"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
	"think.com/go-ddd/domain/port/gateway"
)

func init() {
	gs.Object(new(RiskCheckGatewayImpl)).Export((*gateway.RiskCheckGateway)(nil))
}

// RiskCheckGatewayImpl 风控南向网关实现。
//
// 真实实现应调用风控系统；这里内置一组规则，便于本地验证「命中风控 → 挂起订单」
// 这条业务补偿链路。
type RiskCheckGatewayImpl struct {
	// 单笔订单最多允许的 sku 数，超过视为批量下单
	MaxSkuKinds int `value:"${oms.risk.max-sku-kinds:=20}"`
}

// Check 订单风控校验
func (g *RiskCheckGatewayImpl) Check(req *request.RiskCheckRequest) (*response.RiskCheckResponse, error) {
	resp := &response.RiskCheckResponse{}

	if req.Address == "" {
		resp.IsIllegalAddress = true
		resp.Desc = "收货地址为空"
	}
	if req.PhoneNo == "" {
		resp.IsIllegalAddress = true
		resp.Desc = "收货手机号为空"
	}
	// 演示规则：收件人 / 地址命中黑名单关键字
	if strings.Contains(req.Address, "黑名单") || strings.Contains(req.Username, "黑名单") {
		resp.IsIllegalAddress = true
		resp.Desc = "收货信息命中风控黑名单"
	}
	if len(req.SkuInfos) > g.MaxSkuKinds {
		resp.IsBatchBuy = true
		resp.Desc = "疑似批量下单"
	}
	for _, sku := range req.SkuInfos {
		if sku.SkuPrice <= 0 {
			resp.IsIllegalSkuPrice = true
			resp.Desc = "sku 价格异常"
		}
	}
	return resp, nil
}
