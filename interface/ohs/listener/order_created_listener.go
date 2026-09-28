// Package listener 本地事件监听（北向网关）。
//
// 通过订阅领域事件把「创建订单」与「订单履约」两个用例解耦：
// 创建用例只管发事件，谁来消费它并不关心。
package listener

import (
	"log"

	"github.com/go-spring/spring-core/gs"

	pe "think.com/go-ddd/domain/pl/event"
	infraevent "think.com/go-ddd/infrastructure/core/event"
	"think.com/go-ddd/interface/local"
)

func init() {
	gs.Object(new(OrderCreatedListener)).Init(func(l *OrderCreatedListener) {
		// 对应 Spring 的 @EventListener + @Async
		l.bus.SubscribeAsync(pe.EventNameOrderCreated, l.OnEvent)
	})
}

// OrderCreatedListener 订单创建事件监听
type OrderCreatedListener struct {
	orderLocalService *local.OrderLocalService `autowire:""`
	bus               *infraevent.Bus          `autowire:""`
}

// OnEvent 收到订单创建事件后触发分仓拆单
func (l *OrderCreatedListener) OnEvent(e infraevent.Event) {
	evt, ok := e.(*pe.OrderCreatedEvent)
	if !ok {
		return
	}
	log.Printf("[listener] 收到 OrderCreatedEvent: orderNo=%s", evt.GetOrderNo())
	if err := l.orderLocalService.DispatchOrder(evt.GetOrderNo()); err != nil {
		log.Printf("[listener] 分仓拆单失败 orderNo=%s: %v", evt.GetOrderNo(), err)
	}
}
