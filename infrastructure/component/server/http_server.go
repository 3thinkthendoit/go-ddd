package server

import (
	"github.com/go-spring/spring-core/gs"
	_ "github.com/go-spring/starter-gin"
	"think.com/go-ddd/osh/adapter/http/controller"
)

func init() {
	//获取所有车站信息
	gs.Object(new(controller.StationController)).Init(func(c *controller.StationController) {
		gs.GetMapping("/CreateAllStationsBy12306", c.CreateAllStationsBy12306)
	})
	//获取所有车次信息
}
