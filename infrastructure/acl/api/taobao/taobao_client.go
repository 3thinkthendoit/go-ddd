// Package taobao 淘宝天猫订单协议实现。
package taobao

import (
	"fmt"
	"time"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/domain/model/constant"
	"think.com/go-ddd/domain/pl"
	"think.com/go-ddd/domain/pl/command"
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
	"think.com/go-ddd/infrastructure/common/demo"
)

func init() {
	gs.Object(NewClient())
}

// Client 淘宝开放平台客户端
type Client struct {
	// 是否使用内置示例数据；接入真实环境后置为 false 并实现 API 调用
	MockEnabled bool `value:"${oms.platform-mock.taobao:=true}"`
}

// NewClient 创建客户端
func NewClient() *Client {
	return &Client{}
}

// PullOrder 按更新时间拉取订单并转成创建订单命令
func (c *Client) PullOrder(begin, end time.Time) (*command.OrderCreateCommand, error) {
	if !c.MockEnabled {
		return nil, fmt.Errorf("淘宝订单拉取未接入真实 API, begin=%s end=%s", begin, end)
	}
	return demo.BuildOrderCommand("TB", constant.OrderSourceTaoBao,
		constant.PayTypeAlipay, "STORE_TB_001", "淘宝天猫订单"), nil
}

// ShippingCallback 发货回传
func (c *Client) ShippingCallback(req *request.ShippingCallbackRequest) (*response.ShippingCallbackResponse, error) {
	// 真实实现：按淘宝开放平台协议签名后调用接口
	return &response.ShippingCallbackResponse{
		CallbackResult: &pl.ShippingCallbackResult{
			CallStatus: 0,
			Result:     fmt.Sprintf("taobao callback success, externalOrderNo=%s", req.ExternalOrderNo),
		},
	}, nil
}
