// Package rpc 微服务接口（北向网关）。
package rpc

import "github.com/go-spring/spring-core/web"

// OrderInfoIface 订单微服务接口，供其他微服务通过 HTTP 调用
type OrderInfoIface interface {
	// Query 订单查询
	Query(ctx web.Context)
}
