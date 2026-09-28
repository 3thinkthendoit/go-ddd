package event

// OrderCreatedEvent 订单创建完成事件。
//
// 由订单创建用例在持久化之后发布，驱动后续的订单分仓 / 拆单，
// 使「创建」与「履约」两个用例解耦。
type OrderCreatedEvent struct {
	OrderOperationEvent
}

// NewOrderCreatedEvent 创建订单创建事件（默认本地事件）
func NewOrderCreatedEvent(orderNo string) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		OrderOperationEvent: OrderOperationEvent{
			PublishType: PublishTypeLocal,
			OrderNo:     orderNo,
		},
	}
}

// EventName 事件名
func (e *OrderCreatedEvent) EventName() string { return EventNameOrderCreated }
