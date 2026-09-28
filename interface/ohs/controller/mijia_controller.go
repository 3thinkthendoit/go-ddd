// Package controller HTTP 接口（北向网关）。
//
// 只做「协议解析 → 调本地服务 → 组织返回」，不写业务逻辑。
package controller

import (
	"github.com/go-spring/spring-core/gs"
	"github.com/go-spring/spring-core/web"

	"think.com/go-ddd/infrastructure/acl/api/mijia"
	"think.com/go-ddd/interface/local"
	"think.com/go-ddd/interface/ohs/dto/resp"
	"think.com/go-ddd/interface/ohs/httpx"
)

func init() {
	gs.Object(new(MiJiaController)).Init(func(c *MiJiaController) {
		gs.PostMapping("/api/order/mijia/create", c.Create)
		gs.PostMapping("/api/order/mijia/query", c.Query)
	})
}

// MiJiaController 米家订单接口（小米主动推送）
type MiJiaController struct {
	orderLocalService *local.OrderLocalService `autowire:""`
	miJiaClient       *mijia.Client            `autowire:""`
}

// Create 接收米家订单
func (c *MiJiaController) Create(ctx web.Context) {
	body, err := ctx.RequestBody()
	if err != nil {
		httpx.WriteFailure(ctx, err)
		return
	}
	// 解析米家协议 → 创建订单命令
	cmd, err := c.miJiaClient.ParseCreateOrder(body)
	if err != nil {
		httpx.WriteFailure(ctx, err)
		return
	}
	if err := c.orderLocalService.CreateOrder(cmd); err != nil {
		httpx.WriteFailure(ctx, err)
		return
	}
	httpx.WriteSuccess(ctx, map[string]interface{}{
		"externalOrderNo": cmd.ExternalOrderNo,
		"msg":             "订单已接入，后续由 OrderCreatedEvent 驱动分仓拆单",
	})
}

// Query 订单查询
func (c *MiJiaController) Query(ctx web.Context) {
	body, err := ctx.RequestBody()
	if err != nil {
		httpx.WriteFailure(ctx, err)
		return
	}
	q, err := c.miJiaClient.ParseQuery(body)
	if err != nil {
		httpx.WriteFailure(ctx, err)
		return
	}
	orders, err := c.orderLocalService.Query(q)
	if err != nil {
		httpx.WriteFailure(ctx, err)
		return
	}
	httpx.WriteSuccess(ctx, resp.NewOrderInfoQueryResps(orders))
}
