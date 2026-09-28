package repository

import (
	"log"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/domain/model/aggregate/shippingcallback"
	"think.com/go-ddd/domain/model/constant"
	"think.com/go-ddd/domain/port/repository"
	aclpl "think.com/go-ddd/infrastructure/acl/pl"
	"think.com/go-ddd/infrastructure/core/orm"
	"think.com/go-ddd/infrastructure/core/orm/po"
)

func init() {
	gs.Object(new(SkuShippingRepositoryImpl)).Export((*repository.SkuShippingRepository)(nil))
}

// SkuShippingRepositoryImpl 发货回传资源库实现。
type SkuShippingRepositoryImpl struct {
	store *orm.Store `autowire:""`
}

// Save 持久化发货回传结果：更新 sku item 发货数量 + 新增回传记录
func (r *SkuShippingRepositoryImpl) Save(aggregate *shippingcallback.ShippingCallbackAggregate) error {
	orderNo := aggregate.OrderId().OrderNo()
	for skuCode, item := range aggregate.OrderSkuItems() {
		r.store.UpdateOrderSkuItemInfos(
			func(v *po.OrderSkuItemInfo) bool {
				return v.OrderNo == orderNo && v.SkuCode == skuCode
			},
			func(v *po.OrderSkuItemInfo) {
				v.ShippingAmount = item.ShippingAmount()
				v.SkuStatus = skuStatusOf(item).Code()
			})
	}

	record := aclpl.ToShippingCallbackRecordPO(aggregate)
	if record == nil {
		return nil
	}
	if err := r.store.InsertShippingCallbackRecord(record); err != nil {
		return err
	}
	log.Printf("[repository] 发货回传已落库 orderNo=%s, callStatus=%d", orderNo, record.CallStatus)
	return nil
}

// skuStatusOf 由发货数量推导 sku 发货状态
func skuStatusOf(item *shippingcallback.ShippingSkuItem) constant.OrderSkuStatus {
	switch {
	case item.ShippingAmount() <= 0:
		return constant.OrderSkuStatusNotShipped
	case item.ShippingAmount() < item.SkuAmount():
		return constant.OrderSkuStatusPartShipped
	default:
		return constant.OrderSkuStatusShipped
	}
}
