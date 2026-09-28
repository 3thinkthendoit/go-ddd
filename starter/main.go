// Package main 启动类。
//
// 通过 import 触发各层包的 init()，完成 bean 注册 —— 等价于 Spring 的
// @ComponentScan；领域层不引入框架依赖，其装配见 wire.go。
package main

import (
	"github.com/go-spring/spring-core/gs"

	// 应用层
	_ "think.com/go-ddd/application/service"

	// 基础设施层：南向网关实现
	_ "think.com/go-ddd/infrastructure/acl/api/douyin"
	_ "think.com/go-ddd/infrastructure/acl/api/mijia"
	_ "think.com/go-ddd/infrastructure/acl/api/pdd"
	_ "think.com/go-ddd/infrastructure/acl/api/taobao"
	_ "think.com/go-ddd/infrastructure/acl/gateway"
	_ "think.com/go-ddd/infrastructure/acl/publisher"
	_ "think.com/go-ddd/infrastructure/acl/repository"

	// 基础设施层：技术组件
	_ "think.com/go-ddd/infrastructure/core/event"
	_ "think.com/go-ddd/infrastructure/core/http"
	_ "think.com/go-ddd/infrastructure/core/mq"
	_ "think.com/go-ddd/infrastructure/core/orm"
	_ "think.com/go-ddd/infrastructure/core/redis"
	_ "think.com/go-ddd/infrastructure/core/scheduler"
	_ "think.com/go-ddd/infrastructure/core/server"

	// 接口层：北向网关
	_ "think.com/go-ddd/interface/local"
	_ "think.com/go-ddd/interface/ohs/controller"
	_ "think.com/go-ddd/interface/ohs/job"
	_ "think.com/go-ddd/interface/ohs/listener"
	_ "think.com/go-ddd/interface/ohs/mq"
	_ "think.com/go-ddd/interface/ohs/rpc/impl"
)

func init() {
	gs.Setenv("GS_SPRING_CONFIG_LOCATIONS", "infrastructure/config/")
}

func main() {
	if err := gs.Run(); err != nil {
		panic(err)
	}
}
