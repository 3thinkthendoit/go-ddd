package event

// PublishType 事件发布方式
type PublishType int

const (
	// PublishTypeLocal 本地事件（进程内，等价于 Spring 的 ApplicationEvent）
	PublishTypeLocal PublishType = 1
	// PublishTypeMq MQ 事件（跨进程）
	PublishTypeMq PublishType = 2
)

// 事件名常量：发布方与订阅方共用，避免两边写死字符串不一致。
const (
	// EventNameOrderCreated 订单创建完成
	EventNameOrderCreated = "OrderCreatedEvent"
	// EventNameOrderFulfill 订单分仓拆单完成
	EventNameOrderFulfill = "OrderFulfillEvent"
)

// OrderEvent 订单事件抽象。
//
// 领域服务只依赖该接口，具体走本地事件总线还是 MQ 由基础设施层的
// OrderEventPublisher 实现决定。
type OrderEvent interface {
	// EventName 事件名，用于事件总线订阅
	EventName() string
	// GetOrderNo 关联的订单号
	GetOrderNo() string
	// GetPublishType 发布方式
	GetPublishType() PublishType
}

// OrderOperationEvent 订单操作事件基类
type OrderOperationEvent struct {
	PublishType PublishType
	OrderNo     string
}

// EventName 默认事件名，子类可覆盖
func (e *OrderOperationEvent) EventName() string { return "OrderOperationEvent" }

// GetOrderNo 关联的订单号
func (e *OrderOperationEvent) GetOrderNo() string { return e.OrderNo }

// GetPublishType 发布方式
func (e *OrderOperationEvent) GetPublishType() PublishType { return e.PublishType }
