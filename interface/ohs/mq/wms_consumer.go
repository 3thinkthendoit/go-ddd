// Package mq 远程事件监听（北向网关）。
package mq

import (
	"encoding/json"
	"log"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/domain/pl/command"
	inframq "think.com/go-ddd/infrastructure/core/mq"
	"think.com/go-ddd/interface/local"
)

func init() {
	gs.Object(new(WmsConsumer)).Init(func(c *WmsConsumer) {
		if err := c.mqClient.Subscribe(c.Topic, c.ConsumerGroup, c.OnMessage); err != nil {
			log.Printf("[wms-consumer] 订阅失败: %v", err)
		}
	})
}

// wmsMessage WMS 发货消息报文
type wmsMessage struct {
	// OMS 订单号
	OrderNo string `json:"orderNo"`
	// WMS 发货单号（拆单后的父单号）
	OmsOrderNo    string         `json:"omsOrderNo"`
	ShippingInfos []shippingInfo `json:"shippingInfos"`
}

type shippingInfo struct {
	SkuCode     string `json:"skuCode"`
	SkuAmount   int    `json:"skuAmount"`
	ExpressCode string `json:"expressCode"`
	ExpressNo   string `json:"expressNo"`
}

// WmsConsumer 监听履约系统（WMS）的发货消息
type WmsConsumer struct {
	mqClient          *inframq.RocketMqClient  `autowire:""`
	orderLocalService *local.OrderLocalService `autowire:""`

	// 消费的 topic 与消费组
	Topic         string `value:"${oms.wms.callback-topic:=order-fulfillment-center}"`
	ConsumerGroup string `value:"${oms.wms.consumer-group:=order-fulfillment-center-group}"`
}

// OnMessage 消费发货消息
func (c *WmsConsumer) OnMessage(_ string, message string) error {
	log.Printf("[wms-consumer] 收到 WMS 发货信息 msg=%s", message)

	var msg wmsMessage
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return err
	}

	shippingInfos := make([]command.ShippingInfo, 0, len(msg.ShippingInfos))
	for _, item := range msg.ShippingInfos {
		shippingInfos = append(shippingInfos, command.ShippingInfo{
			SkuCode:     item.SkuCode,
			SkuAmount:   item.SkuAmount,
			ExpressCode: item.ExpressCode,
			ExpressNo:   item.ExpressNo,
		})
	}

	return c.orderLocalService.ShippingCallback(&command.SkuShippingCommand{
		OrderNo:       msg.OrderNo,
		WmsOrderNo:    msg.OmsOrderNo,
		ShippingInfos: shippingInfos,
	})
}
