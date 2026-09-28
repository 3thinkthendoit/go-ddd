package service

import (
	"log"

	"think.com/go-ddd/domain/model/aggregate/shippingcallback"
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/port/gateway"
)

// OrderShippingDomainService 发货回传领域服务。
//
// 负责补全回传所需的订单信息、发货信息，并在需要时通知上游平台。
type OrderShippingDomainService struct {
	orderInfoGateway        gateway.OrderInfoGateway        `autowire:""`
	shippingCallbackGateway gateway.ShippingCallbackGateway `autowire:""`
	orderFulfillGateway     gateway.OrderFulfillGateway     `autowire:""`
}

// InitBaseInfo 完善订单信息与 sku 发货信息
func (s *OrderShippingDomainService) InitBaseInfo(aggregate *shippingcallback.ShippingCallbackAggregate) error {
	if err := s.initOrderInfo(aggregate); err != nil {
		return err
	}
	return s.initSkuShippingInfo(aggregate)
}

// initOrderInfo 按订单号查询订单，补全订单标识与 sku 明细
func (s *OrderShippingDomainService) initOrderInfo(aggregate *shippingcallback.ShippingCallbackAggregate) error {
	orderNo := aggregate.OrderId().OrderNo()
	resp, err := s.orderInfoGateway.Query(&request.OrderQueryRequest{OrderNo: orderNo})
	if err != nil {
		return err
	}
	if resp == nil || len(resp.Orders) == 0 {
		return nil
	}
	orderInfo := resp.Orders[0]
	aggregate.CompleteOrderId(orderInfo.ExternalOrderNo, orderInfo.OrderSource)
	return aggregate.InitBaseInfo(orderInfo.SkuInfos)
}

// initSkuShippingInfo 调用 WMS 网关获取发货信息。
//
// WMS 的发货数量已随命令传入，这里仅保留按需回查 WMS 的扩展点；
// 注意查的是拆单后的**发货单号**（WmsOrderNo），不是 OMS 订单号。
func (s *OrderShippingDomainService) initSkuShippingInfo(aggregate *shippingcallback.ShippingCallbackAggregate) error {
	omsOrderNos := make([]string, 0, 1)
	if wmsOrderNo := aggregate.WmsOrderNo(); wmsOrderNo != "" {
		omsOrderNos = append(omsOrderNos, wmsOrderNo)
	}
	_, err := s.orderFulfillGateway.Query(&request.ShippingQueryRequest{
		OmsOrderNos: omsOrderNos,
	})
	return err
}

// ShippingCallback 通知上游平台发货回传
func (s *OrderShippingDomainService) ShippingCallback(aggregate *shippingcallback.ShippingCallbackAggregate) error {
	if !aggregate.Callback() {
		log.Printf("[shipping] orderNo=%s 无需发货回传", aggregate.OrderId().OrderNo())
		return nil
	}
	shippingInfos := make([]request.ShippingCallbackInfo, 0, len(aggregate.OrderSkuItems()))
	for _, item := range aggregate.OrderSkuItems() {
		shippingInfos = append(shippingInfos, request.ShippingCallbackInfo{
			ExternalSkuCode: item.ExternalSkuId(),
			SkuAmount:       item.SkuAmount(),
			ShippingAmount:  item.ShippingAmount(),
			ExpressCode:     item.ExpressCode(),
			ExpressNo:       item.ExpressNo(),
		})
	}
	resp, err := s.shippingCallbackGateway.Callback(&request.ShippingCallbackRequest{
		ExternalOrderNo: aggregate.OrderId().ExternalOrderNo(),
		OrderSource:     aggregate.OrderId().OrderSource(),
		ShippingInfos:   shippingInfos,
	})
	if err != nil {
		return err
	}
	return aggregate.HandleCallbackResult(resp.CallbackResult)
}
