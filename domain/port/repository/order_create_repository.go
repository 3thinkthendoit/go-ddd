package repository

import (
	"think.com/go-ddd/domain/model/aggregate/create"
	"think.com/go-ddd/domain/model/aggregate/shippingcallback"
)

// OrderCreateRepository 订单创建资源库
type OrderCreateRepository interface {
	// Save 持久化订单创建聚合
	Save(aggregate *create.OrderCreateAggregate) error
	// Update 根据发货回传聚合更新订单
	Update(aggregate *shippingcallback.ShippingCallbackAggregate) error
}
