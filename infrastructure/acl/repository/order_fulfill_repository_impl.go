package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/domain/model/aggregate/orderfulfill"
	"think.com/go-ddd/domain/model/constant"
	domainpl "think.com/go-ddd/domain/pl"
	"think.com/go-ddd/domain/port/repository"
	aclpl "think.com/go-ddd/infrastructure/acl/pl"
	"think.com/go-ddd/infrastructure/core/orm"
	"think.com/go-ddd/infrastructure/core/orm/po"
)

func init() {
	gs.Object(new(OrderFulfillRepositoryImpl)).Export((*repository.OrderFulfillRepository)(nil))
}

// OrderFulfillRepositoryImpl 订单履约资源库实现。
type OrderFulfillRepositoryImpl struct {
	store *orm.Store `autowire:""`
}

// OfByOrderNo 按订单号还原履约聚合根（只还原待履约的已支付订单）
func (r *OrderFulfillRepositoryImpl) OfByOrderNo(orderNo string) (*orderfulfill.OrderFulfillAggregate, error) {
	orders := r.store.SelectOrderBaseInfos(func(v *po.OrderBaseInfo) bool {
		return v.OrderNo == orderNo && v.OrderStatus == constant.OrderStatusPayed.Code()
	})
	if len(orders) == 0 {
		return nil, fmt.Errorf("根据 orderNo=[%s] 查询不到需要履约的订单信息", orderNo)
	}
	skuItemInfos := r.store.SelectOrderSkuItemInfos(func(v *po.OrderSkuItemInfo) bool {
		return v.OrderNo == orderNo
	})
	return aclpl.ToOrderFulfillAggregate(orders[0], skuItemInfos)
}

// Save 持久化履约聚合的拆单结果
func (r *OrderFulfillRepositoryImpl) Save(aggregate *orderfulfill.OrderFulfillAggregate) error {
	splitResults := aclpl.ToOrderSplitResultPOs(aggregate)
	if len(splitResults) == 0 {
		return nil
	}
	if err := r.store.InsertOrderSplitResults(splitResults); err != nil {
		return err
	}
	log.Printf("[repository] 拆单结果已落库 orderNo=%s, 发货单数=%d, 明细数=%d",
		aggregate.OrderNo(), len(aggregate.SplitOrders()), len(splitResults))
	return nil
}

// UpdateOrderFulfill 标记发货单已推送 WMS
func (r *OrderFulfillRepositoryImpl) UpdateOrderFulfill(omsOrderNo string) error {
	affected := r.store.UpdateOrderSplitResults(
		func(v *po.OrderSplitResult) bool {
			return v.ParentOrderNo == omsOrderNo && v.FulfillStatus == 0
		},
		func(v *po.OrderSplitResult) {
			v.FulfillStatus = 1
			v.UpdateTime = time.Now()
		})
	if affected > 0 {
		log.Printf("[repository] 发货单已推送 WMS, omsOrderNo=%s", omsOrderNo)
	}
	return nil
}

// QueryFulfillOrderInfos 查询指定订单下待推送 WMS 的发货单
func (r *OrderFulfillRepositoryImpl) QueryFulfillOrderInfos(orderNo string) ([]*domainpl.FulfillOrderInfo, error) {
	splitResults := r.store.SelectOrderSplitResults(func(v *po.OrderSplitResult) bool {
		return v.OrderNo == orderNo && v.FulfillStatus == 0
	})
	return aclpl.ToFulfillOrderInfos(splitResults), nil
}
