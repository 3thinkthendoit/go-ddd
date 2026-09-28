package repository

import (
	"think.com/go-ddd/domain/model/aggregate/orderfulfill"
	"think.com/go-ddd/domain/pl"
)

// OrderFulfillRepository 订单履约资源库
type OrderFulfillRepository interface {
	// OfByOrderNo 按订单号还原履约聚合根
	OfByOrderNo(orderNo string) (*orderfulfill.OrderFulfillAggregate, error)
	// Save 持久化履约聚合（拆单结果）
	Save(aggregate *orderfulfill.OrderFulfillAggregate) error
	// UpdateOrderFulfill 更新发货单履约状态
	UpdateOrderFulfill(omsOrderNo string) error
	// QueryFulfillOrderInfos 查询待推送 WMS 的发货单
	QueryFulfillOrderInfos(omsOrderNo string) ([]*pl.FulfillOrderInfo, error)
}
