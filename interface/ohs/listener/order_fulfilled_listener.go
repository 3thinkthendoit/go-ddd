package listener

import (
	"log"

	"github.com/go-spring/spring-core/gs"

	pe "think.com/go-ddd/domain/pl/event"
	infraevent "think.com/go-ddd/infrastructure/core/event"
	"think.com/go-ddd/interface/local"
)

func init() {
	gs.Object(new(OrderFulfilledListener)).Init(func(l *OrderFulfilledListener) {
		l.bus.SubscribeAsync(pe.EventNameOrderFulfill, l.OnEvent)
	})
}

// OrderFulfilledListener 订单履约事件监听
type OrderFulfilledListener struct {
	orderLocalService *local.OrderLocalService `autowire:""`
	bus               *infraevent.Bus          `autowire:""`
}

// OnEvent 收到履约事件后推送 WMS
func (l *OrderFulfilledListener) OnEvent(e infraevent.Event) {
	evt, ok := e.(*pe.OrderFulfillEvent)
	if !ok {
		return
	}
	log.Printf("[listener] 收到 OrderFulfillEvent: orderNo=%s", evt.GetOrderNo())
	if err := l.orderLocalService.FulfillOrder(evt.GetOrderNo()); err != nil {
		log.Printf("[listener] 推送 WMS 失败 orderNo=%s: %v", evt.GetOrderNo(), err)
	}
}
