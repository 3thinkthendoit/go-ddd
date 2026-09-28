package event

// OrderFulfillEvent 订单分仓拆单完成事件。
//
// 拆单结果落库后发布，驱动后续推送 WMS 履约。
type OrderFulfillEvent struct {
	OrderOperationEvent
}

// NewOrderFulfillEvent 创建订单履约事件（默认本地事件）
func NewOrderFulfillEvent(orderNo string) *OrderFulfillEvent {
	return &OrderFulfillEvent{
		OrderOperationEvent: OrderOperationEvent{
			PublishType: PublishTypeLocal,
			OrderNo:     orderNo,
		},
	}
}

// EventName 事件名
func (e *OrderFulfillEvent) EventName() string { return EventNameOrderFulfill }
