package service

import "think.com/go-ddd/domain/station/model/command"

type StationDomainService struct {
}

/**
 * 获取所有车站领域服务
 * 返回车站信息
**/
func (s *StationDomainService) GetStationInfos(dataSource string) []*command.CreateStationCmd {
	if "12306" == dataSource {

	} else if "9zhouchuxing" == dataSource {

	}
	cmds := []*command.CreateStationCmd{}
	return cmds
}
