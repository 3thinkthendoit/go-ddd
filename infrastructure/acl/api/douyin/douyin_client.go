// Package douyin 抖音订单协议实现。
package douyin

import (
	"fmt"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/domain/pl"
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
)

func init() {
	gs.Object(NewClient())
}

// Client 抖音开放平台客户端
type Client struct {
	// 是否使用内置示例数据
	MockEnabled bool `value:"${oms.platform-mock.douyin:=true}"`
}

// NewClient 创建客户端
func NewClient() *Client {
	return &Client{}
}

// ShippingCallback 发货回传
//
// 抖音平台无需 OMS 主动回传发货信息，这里保留协议实现位。
func (c *Client) ShippingCallback(req *request.ShippingCallbackRequest) (*response.ShippingCallbackResponse, error) {
	return &response.ShippingCallbackResponse{
		CallbackResult: &pl.ShippingCallbackResult{
			CallStatus: 0,
			Result:     fmt.Sprintf("douyin callback skipped, externalOrderNo=%s", req.ExternalOrderNo),
		},
	}, nil
}
