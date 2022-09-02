package service

import (
	"think.com/go-ddd/domain/station/model/aggregate"
	"think.com/go-ddd/domain/station/port"
	"think.com/go-ddd/domain/station/service"
)

func init() {

}

type StationAppService struct {
	stationDomainService *service.StationDomainService `autowire:"?"`
	stationRepository    port.StationRepository        `autowire:"?"`
}

/**
	业务编排,通用新增所有车站信息
**/
func (s *StationAppService) CreateAllStations(dataSource string) {
	//通过domian service从南向网关获取数据
	cmds := s.stationDomainService.GetStationInfos(dataSource)
	for i := 0; i < len(cmds); i++ {
		//批量还原聚合根
		stationAggregate := aggregate.CreateStationAggregate(cmds[i])
		//聚合领域方法
		stationAggregate.Enable()
		//资源库持久化
		s.stationRepository.Save(stationAggregate)
	}

}
