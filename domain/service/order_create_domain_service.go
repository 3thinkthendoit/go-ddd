package service

import (
	"fmt"

	"think.com/go-ddd/domain/model/aggregate/create"
	"think.com/go-ddd/domain/pl"
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/port/gateway"
)

// OrderCreateDomainService 订单创建领域服务。
//
// 领域服务只承载「需要与其他域协作」的领域行为，纯业务规则仍然放在聚合根里。
type OrderCreateDomainService struct {
	orderInfoGateway    gateway.OrderInfoGateway    `autowire:""`
	skuInfoQueryGateway gateway.SkuInfoQueryGateway `autowire:""`
	riskCheckGateway    gateway.RiskCheckGateway    `autowire:""`
}

// IsExist 判断订单是否已经接入过（幂等保护）
func (s *OrderCreateDomainService) IsExist(aggregate *create.OrderCreateAggregate) error {
	orderId := aggregate.OrderId()
	resp, err := s.orderInfoGateway.Query(&request.OrderQueryRequest{
		ExternalOrderNo: orderId.ExternalOrderNo(),
		OrderSource:     orderId.OrderSource(),
	})
	if err != nil {
		return err
	}
	if resp != nil && len(resp.Orders) > 0 {
		return fmt.Errorf("订单 externalOrderNo=%s 已经存在", orderId.ExternalOrderNo())
	}
	return nil
}

// InitBaseInfo 完善订单基础信息（店铺、商品、发票、地址、用户）
func (s *OrderCreateDomainService) InitBaseInfo(aggregate *create.OrderCreateAggregate) error {
	if err := s.initStoreInfo(aggregate); err != nil {
		return err
	}
	if err := s.initSkuInfo(aggregate); err != nil {
		return err
	}
	if err := s.initInvoiceInfo(aggregate); err != nil {
		return err
	}
	if err := s.initShippingAddress(aggregate); err != nil {
		return err
	}
	return s.initBuyer(aggregate)
}

// Audit 订单审核：聚合自身规则 + 外域风控
func (s *OrderCreateDomainService) Audit(aggregate *create.OrderCreateAggregate) error {
	if err := aggregate.Check(); err != nil {
		return err
	}
	return s.riskCheck(aggregate)
}

// initSkuInfo 调用商品域，把外部 sku 换成内部 sku
func (s *OrderCreateDomainService) initSkuInfo(aggregate *create.OrderCreateAggregate) error {
	externalSkuIds := make([]string, 0, len(aggregate.SkuInfos()))
	for _, orderSku := range aggregate.SkuInfos() {
		externalSkuIds = append(externalSkuIds, orderSku.SkuInfo().ExternalSkuId())
	}
	resp, err := s.skuInfoQueryGateway.Query(&request.SkuInfoQueryRequest{ExternalSkuIds: externalSkuIds})
	if err != nil {
		return err
	}
	skuInfoMap := make(map[string]*pl.SkuFullInfo, len(externalSkuIds))
	if resp != nil {
		for _, sku := range resp.SkuInfos {
			skuInfoMap[sku.ExternalSkuId] = sku
		}
	}
	return aggregate.ModifyOrderSku(skuInfoMap)
}

// initInvoiceInfo 查询发票域完善发票信息
func (s *OrderCreateDomainService) initInvoiceInfo(aggregate *create.OrderCreateAggregate) error {
	// 需要时调用发票域 gateway 补全 OrderInvoice
	return nil
}

// initShippingAddress 查询基础信息域完善下单地址（补标准地址编码）
func (s *OrderCreateDomainService) initShippingAddress(aggregate *create.OrderCreateAggregate) error {
	// 需要时调用基础信息域 gateway，回填 addressCode
	return nil
}

// initBuyer 查询用户域完善用户信息
func (s *OrderCreateDomainService) initBuyer(aggregate *create.OrderCreateAggregate) error {
	return nil
}

// initStoreInfo 查询店铺域完善店铺信息
func (s *OrderCreateDomainService) initStoreInfo(aggregate *create.OrderCreateAggregate) error {
	return nil
}

// riskCheck 风控处理：命中规则不直接失败，而是挂起订单走业务补偿
func (s *OrderCreateDomainService) riskCheck(aggregate *create.OrderCreateAggregate) error {
	address := aggregate.ShippingAddress()
	skuInfos := make([]request.RiskCheckSkuInfo, 0, len(aggregate.SkuInfos()))
	for _, orderSku := range aggregate.SkuInfos() {
		skuInfos = append(skuInfos, request.RiskCheckSkuInfo{
			SkuCode:   orderSku.SkuInfo().SkuCode(),
			SkuAmount: orderSku.SkuBuyAmount(),
			SkuPrice:  orderSku.SkuPayPrice(),
		})
	}
	resp, err := s.riskCheckGateway.Check(&request.RiskCheckRequest{
		Username: aggregate.Buyer().Username(),
		Address:  address.Address(),
		PhoneNo:  address.ContactInfo(),
		SkuInfos: skuInfos,
	})
	if err != nil {
		return err
	}
	if resp != nil && resp.Rejected() {
		aggregate.Hangup(resp.Desc)
	}
	return nil
}
