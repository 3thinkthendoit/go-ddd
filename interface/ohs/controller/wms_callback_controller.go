package controller

import (
	"github.com/go-spring/spring-core/gs"
	"github.com/go-spring/spring-core/web"

	inframq "think.com/go-ddd/infrastructure/core/mq"
	"think.com/go-ddd/interface/ohs/httpx"
)

func init() {
	gs.Object(new(WmsCallbackController)).Init(func(c *WmsCallbackController) {
		gs.PostMapping("/api/wms/callback", c.Callback)
	})
}

// WmsCallbackController 模拟 WMS 推送发货消息的入口。
//
// 真实环境 WMS 直接把消息投递到 MQ，由 mq.WmsConsumer 消费；这里提供一个 HTTP
// 入口把报文投递到同一个 topic，方便本地一条命令跑通「发货回传」链路。
type WmsCallbackController struct {
	mqClient *inframq.RocketMqClient `autowire:""`

	// WMS 发货消息 topic
	Topic string `value:"${oms.wms.callback-topic:=order-fulfillment-center}"`
}

// Callback 把 WMS 报文投递到 MQ
func (c *WmsCallbackController) Callback(ctx web.Context) {
	body, err := ctx.RequestBody()
	if err != nil {
		httpx.WriteFailure(ctx, err)
		return
	}
	if err := c.mqClient.Send(c.Topic, string(body)); err != nil {
		httpx.WriteFailure(ctx, err)
		return
	}
	httpx.WriteSuccess(ctx, map[string]interface{}{
		"msg":   "发货消息已投递到 MQ",
		"topic": c.Topic,
	})
}
