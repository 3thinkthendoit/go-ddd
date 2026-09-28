package publisher

import (
	"github.com/go-spring/spring-core/gs"

	pe "think.com/go-ddd/domain/pl/event"
	"think.com/go-ddd/domain/port/publisher"
	"think.com/go-ddd/infrastructure/common/util"
	"think.com/go-ddd/infrastructure/core/event"
	"think.com/go-ddd/infrastructure/core/mq"
)

func init() {
	gs.Object(new(OrderEventPublisherImpl)).Export((*publisher.OrderEventPublisher)(nil))
}

// OrderEventPublisherImpl 领域事件发布实现。
//
// 本地事件走进程内事件总线（等价 Spring ApplicationEvent），
// MQ 事件走消息队列，发布方式由事件自身声明。
type OrderEventPublisherImpl struct {
	bus      *event.Bus         `autowire:""`
	mqClient *mq.RocketMqClient `autowire:""`

	// MQ 事件默认 topic
	Topic string `value:"${oms.event.topic:=think-oms}"`
}

// Publish 发布订单事件
func (p *OrderEventPublisherImpl) Publish(evt pe.OrderEvent) error {
	if evt.GetPublishType() == pe.PublishTypeLocal {
		p.bus.Publish(evt)
		return nil
	}
	return p.mqClient.Send(p.Topic, util.ToJson(evt))
}
