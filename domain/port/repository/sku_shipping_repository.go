package repository

import "think.com/go-ddd/domain/model/aggregate/shippingcallback"

// SkuShippingRepository 发货回传资源库
type SkuShippingRepository interface {
	// Save 持久化发货回传聚合
	Save(aggregate *shippingcallback.ShippingCallbackAggregate) error
}
