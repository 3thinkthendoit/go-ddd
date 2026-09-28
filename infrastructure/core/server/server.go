// Package server HTTP 服务端组件。
//
// 引入 go-spring 的 gin starter，为容器提供 web.Server 与 web.Router bean；
// 具体的路由由接口层的 controller 各自注册（等价于 Spring 的 @RequestMapping 扫描）。
package server

import (
	_ "github.com/go-spring/starter-gin"
)
