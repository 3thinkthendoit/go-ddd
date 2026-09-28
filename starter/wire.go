package main

import (
	"github.com/go-spring/spring-core/gs"

	domainservice "think.com/go-ddd/domain/service"
)

// init 领域层的装配。
//
// 领域层刻意不引入任何框架依赖（不 import gs、不打 autowire 之外的框架代码），
// 所以领域服务的注册集中放在启动类里；应用层、基础设施层、接口层各自在包内
// init() 中注册自己，等价于 Spring 的 @Component / @Service 扫描。
func init() {
	gs.Object(new(domainservice.OrderCreateDomainService))
	gs.Object(new(domainservice.OrderFulfillDomainService))
	gs.Object(new(domainservice.OrderShippingDomainService))
}
