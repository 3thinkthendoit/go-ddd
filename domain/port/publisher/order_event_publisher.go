package publisher

import (
	pe "think.com/go-ddd/domain/pl/event"
)

// OrderEventPublisher 领域事件发布端口。
//
// 领域层只声明「要发事件」，至于发到本地事件总线还是 MQ，
// 由基础设施层的实现决定，领域层不感知。
type OrderEventPublisher interface {
	// Publish 发布订单事件
	Publish(event pe.OrderEvent) error
}
