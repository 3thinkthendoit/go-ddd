// Package local 本地网关。
//
// 供本进程内的北向网关（HTTP / Job / MQ / 本地事件监听）统一调用，
// 屏蔽应用服务的具体实现，也让接口层只依赖这一层。
package local

import (
	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/application/service"
	domainpl "think.com/go-ddd/domain/pl"
	"think.com/go-ddd/domain/pl/command"
	"think.com/go-ddd/domain/pl/query"
)

func init() {
	gs.Object(new(OrderLocalService))
}

// OrderLocalService 订单本地服务
type OrderLocalService struct {
	orderAppService *service.OrderAppService `autowire:""`
}

// CreateOrder 创建订单
func (s *OrderLocalService) CreateOrder(cmd *command.OrderCreateCommand) error {
	return s.orderAppService.CreateOrder(cmd)
}

// DispatchOrder 订单分仓拆单
func (s *OrderLocalService) DispatchOrder(orderNo string) error {
	return s.orderAppService.DispatchOrder(orderNo)
}

// FulfillOrder 推送 WMS 履约
func (s *OrderLocalService) FulfillOrder(orderNo string) error {
	return s.orderAppService.FulfillOrder(orderNo)
}

// ShippingCallback 发货回传
func (s *OrderLocalService) ShippingCallback(cmd *command.SkuShippingCommand) error {
	return s.orderAppService.ShippingCallback(cmd)
}

// OrderAfterSaleService 订单售后
func (s *OrderLocalService) OrderAfterSaleService(cmd *command.OrderAssCommand) error {
	return s.orderAppService.OrderAfterSaleService(cmd)
}

// Query 查询订单
func (s *OrderLocalService) Query(q *query.OrderInfoQuery) ([]*domainpl.OrderInfo, error) {
	return s.orderAppService.Query(q)
}
