package command

import (
	"think.com/go-ddd/domain/common/assert"
)

// ShippingInfo 单条发货信息
type ShippingInfo struct {
	SkuCode     string
	SkuAmount   int
	ExpressCode string
	ExpressNo   string
}

// SkuShippingCommand 发货回传命令（WMS 发货后推送给 OMS）。
//
// 同时携带 OMS 订单号与 WMS 发货单号：
//   - OrderNo 用于定位订单聚合，发货回传的主体是订单而不是发货单；
//   - WmsOrderNo 用于对账与发货单维度的幂等控制。
type SkuShippingCommand struct {
	// OMS 订单号
	OrderNo string
	// WMS 发货单号（拆单后的父单号）
	WmsOrderNo string
	// 发货明细
	ShippingInfos []ShippingInfo
}

// Validate 校验命令必填项
func (c *SkuShippingCommand) Validate() error {
	if err := assert.NotBlank(c.OrderNo, "orderNo is blank"); err != nil {
		return err
	}
	return assert.NotBlank(c.WmsOrderNo, "wmsOrderNo is blank")
}
