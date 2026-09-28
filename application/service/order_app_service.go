package service

import (
	"log"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/domain/model/aggregate/create"
	"think.com/go-ddd/domain/model/aggregate/shippingcallback"
	"think.com/go-ddd/domain/model/constant"
	domainpl "think.com/go-ddd/domain/pl"
	"think.com/go-ddd/domain/pl/command"
	"think.com/go-ddd/domain/pl/event"
	"think.com/go-ddd/domain/pl/query"
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/port/gateway"
	"think.com/go-ddd/domain/port/publisher"
	"think.com/go-ddd/domain/port/repository"
	domainservice "think.com/go-ddd/domain/service"
)

func init() {
	gs.Object(new(OrderAppService))
}

// OrderAppService 订单应用服务（CQRS：命令走领域模型，查询直接走南向网关）。
//
// 应用层只做用例编排，不写业务规则：业务规则在聚合根，跨域协作在领域服务。
type OrderAppService struct {
	orderCreateDomainService   *domainservice.OrderCreateDomainService   `autowire:""`
	orderFulfillDomainService  *domainservice.OrderFulfillDomainService  `autowire:""`
	orderShippingDomainService *domainservice.OrderShippingDomainService `autowire:""`

	orderCreateRepository  repository.OrderCreateRepository  `autowire:""`
	orderFulfillRepository repository.OrderFulfillRepository `autowire:""`
	skuShippingRepository  repository.SkuShippingRepository  `autowire:""`

	orderEventPublisher publisher.OrderEventPublisher `autowire:""`

	orderInfoGateway    gateway.OrderInfoGateway    `autowire:""`
	orderFulfillGateway gateway.OrderFulfillGateway `autowire:""`
}

// CreateOrder 统一创建订单（流程编排-低耦合）
func (s *OrderAppService) CreateOrder(cmd *command.OrderCreateCommand) error {
	aggregate, err := create.CreateOrderCreateAggregate(cmd)
	if err != nil {
		return err
	}
	if err := s.orderCreateDomainService.IsExist(aggregate); err != nil {
		return err
	}
	if err := s.orderCreateDomainService.InitBaseInfo(aggregate); err != nil {
		return err
	}
	if err := s.orderCreateDomainService.Audit(aggregate); err != nil {
		return err
	}
	if err := aggregate.PriceCalculate(); err != nil {
		return err
	}
	if err := aggregate.PriorityProcessing(); err != nil {
		return err
	}
	if err := s.orderCreateRepository.Save(aggregate); err != nil {
		return err
	}
	if aggregate.OrderStatus() != constant.OrderStatusPayed {
		// 风控挂起等非正常状态：等待人工审核，不触发后续履约用例
		log.Printf("[app] 订单状态=%s，跳过履约流程 orderNo=%s, desc=%s",
			aggregate.OrderStatus(), aggregate.OrderId().OrderNo(), aggregate.Desc())
		return nil
	}
	return s.orderEventPublisher.Publish(event.NewOrderCreatedEvent(aggregate.OrderId().OrderNo()))
}

// DispatchOrder 订单分仓、拆单
func (s *OrderAppService) DispatchOrder(orderNo string) error {
	aggregate, err := s.orderFulfillRepository.OfByOrderNo(orderNo)
	if err != nil {
		return err
	}
	if err := s.orderFulfillDomainService.InitBaseInfo(aggregate); err != nil {
		return err
	}
	if err := aggregate.Check(); err != nil {
		return err
	}
	if err := aggregate.Dispatch(); err != nil {
		return err
	}
	if err := aggregate.Split(); err != nil {
		return err
	}
	if err := s.orderFulfillRepository.Save(aggregate); err != nil {
		return err
	}
	return s.orderEventPublisher.Publish(event.NewOrderFulfillEvent(orderNo))
}

// FulfillOrder 推送 WMS 履约（可做业务补偿）。
//
// 这一步业务逻辑很薄，直接绕过领域层操作南向网关。
func (s *OrderAppService) FulfillOrder(orderNo string) error {
	fulfillOrderInfos, err := s.orderFulfillRepository.QueryFulfillOrderInfos(orderNo)
	if err != nil {
		return err
	}
	for _, fulfillOrderInfo := range fulfillOrderInfos {
		skus := make([]request.OrderFulfillSku, 0, len(fulfillOrderInfo.Skus))
		for _, sku := range fulfillOrderInfo.Skus {
			skus = append(skus, request.OrderFulfillSku{SkuCode: sku.SkuCode, SkuAmount: sku.SkuAmount})
		}
		resp, err := s.orderFulfillGateway.Fulfill(&request.OrderFulfillRequest{
			OmsOrderNo:    fulfillOrderInfo.OmsOrderNo,
			OrderNo:       fulfillOrderInfo.OrderNo,
			WarehouseCode: fulfillOrderInfo.WarehouseCode,
			Skus:          skus,
		})
		if err != nil {
			return err
		}
		if resp != nil && resp.IsFulfill {
			if err := s.orderFulfillRepository.UpdateOrderFulfill(fulfillOrderInfo.OmsOrderNo); err != nil {
				return err
			}
		}
	}
	return nil
}

// ShippingCallback 订单发货回传处理
func (s *OrderAppService) ShippingCallback(cmd *command.SkuShippingCommand) error {
	aggregate, err := shippingcallback.CreateShippingCallbackAggregate(cmd)
	if err != nil {
		return err
	}
	if err := s.orderShippingDomainService.InitBaseInfo(aggregate); err != nil {
		return err
	}
	if err := aggregate.Check(); err != nil {
		return err
	}
	// 回传上游平台
	if err := s.orderShippingDomainService.ShippingCallback(aggregate); err != nil {
		return err
	}
	// 发货数量与回传记录一起落库
	if err := s.skuShippingRepository.Save(aggregate); err != nil {
		return err
	}
	return s.orderCreateRepository.Update(aggregate)
}

// OrderAfterSaleService 订单售后
func (s *OrderAppService) OrderAfterSaleService(_ *command.OrderAssCommand) error {
	// 售后链路：创建售后聚合 → 校验 → 落库 → 通知履约/结算域
	return nil
}

// Query 查询订单信息。
//
// 查询类用例可以直接绕过领域层，调用南向网关返回读模型。
func (s *OrderAppService) Query(q *query.OrderInfoQuery) ([]*domainpl.OrderInfo, error) {
	orderSource := constant.OrderSourceUnKnow
	if q.OrderSource != 0 {
		orderSource = constant.OrderSourceOfByCode(q.OrderSource)
	}
	resp, err := s.orderInfoGateway.Query(&request.OrderQueryRequest{
		OrderNo:         q.OrderNo,
		ExternalOrderNo: q.ExternalOrderNo,
		OrderSource:     orderSource,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	return resp.Orders, nil
}
