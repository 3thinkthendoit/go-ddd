package controller

import (
	"github.com/go-spring/spring-core/web"
	"think.com/go-ddd/application/service"
)

type StationController struct {
	stationAppService *service.StationAppService `autowire:"?"`
}

/**
   按产品定制-从12306获取车站信息
**/
func (s *StationController) CreateAllStationsBy12306(ctx web.Context) {
	s.stationAppService.CreateAllStations("12306")
	ctx.JSON(web.SUCCESS)
}
