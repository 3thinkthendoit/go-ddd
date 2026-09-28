package gateway

import (
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
)

// RiskCheckGateway 风控南向网关
type RiskCheckGateway interface {
	// Check 订单风控校验
	Check(request *request.RiskCheckRequest) (*response.RiskCheckResponse, error)
}
