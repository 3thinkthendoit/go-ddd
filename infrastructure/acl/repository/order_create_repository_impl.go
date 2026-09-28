package repository

import (
	"log"
	"time"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/domain/model/aggregate/create"
	"think.com/go-ddd/domain/model/aggregate/shippingcallback"
	"think.com/go-ddd/domain/model/constant"
	"think.com/go-ddd/domain/port/repository"
	aclpl "think.com/go-ddd/infrastructure/acl/pl"
	"think.com/go-ddd/infrastructure/core/orm"
	"think.com/go-ddd/infrastructure/core/orm/po"
)

func init() {
	gs.Object(new(OrderCreateRepositoryImpl)).Export((*repository.OrderCreateRepository)(nil))
}

// OrderCreateRepositoryImpl 订单创建资源库实现。
type OrderCreateRepositoryImpl struct {
	store *orm.Store `autowire:""`
}

// Save 持久化订单创建聚合。
//
// 主表 + sku + sku item 三张表一次写入，任意一步失败直接返回 error，
// 真实数据库实现应放在同一个事务里。
func (r *OrderCreateRepositoryImpl) Save(aggregate *create.OrderCreateAggregate) error {
	baseInfo := aclpl.ToOrderBaseInfoPO(aggregate)
	skuInfos := aclpl.ToOrderSkuPOs(aggregate)
	skuItemInfos := aclpl.ToOrderSkuItemPOs(aggregate)

	if err := r.store.InsertOrderBaseInfo(baseInfo); err != nil {
		return err
	}
	if err := r.store.InsertOrderSkuInfos(skuInfos); err != nil {
		return err
	}
	if err := r.store.InsertOrderSkuItemInfos(skuItemInfos); err != nil {
		return err
	}
	log.Printf("[repository] 订单创建聚合已落库 orderNo=%s, sku=%d, skuItem=%d, status=%s",
		baseInfo.OrderNo, len(skuInfos), len(skuItemInfos), aggregate.OrderStatus())
	return nil
}

// Update 根据发货回传聚合更新订单主表状态。
//
// 规则：所有需要发货的 sku item 都已发完，则主订单置为「已发货」。
func (r *OrderCreateRepositoryImpl) Update(aggregate *shippingcallback.ShippingCallbackAggregate) error {
	orderNo := aggregate.OrderId().OrderNo()
	skuItemInfos := r.store.SelectOrderSkuItemInfos(func(v *po.OrderSkuItemInfo) bool {
		return v.OrderNo == orderNo
	})
	if len(skuItemInfos) == 0 {
		return nil
	}

	allShipped := true
	for _, item := range skuItemInfos {
		// 优先级为 0 表示虚拟商品，无需发货
		if item.Priority == 0 {
			continue
		}
		if item.ShippingAmount < item.SkuAmount {
			allShipped = false
			break
		}
	}
	if !allShipped {
		return nil
	}

	affected := r.store.UpdateOrderBaseInfo(
		func(v *po.OrderBaseInfo) bool {
			return v.OrderNo == orderNo && v.OrderStatus != constant.OrderStatusFulfilled.Code()
		},
		func(v *po.OrderBaseInfo) {
			v.OrderStatus = constant.OrderStatusFulfilled.Code()
			v.UpdateTime = time.Now()
		})
	if affected > 0 {
		log.Printf("[repository] 订单已全部发货，主订单状态更新为已发货 orderNo=%s", orderNo)
	}
	return nil
}
