package aggregate

import (
	"think.com/go-ddd/domain/station/model/command"
	"think.com/go-ddd/domain/station/model/vo"
)

type StationAggregate struct {
	stationId *vo.StationId //Domain primitive 类型

}

/**
 *领域方法-启用车站
**/
func (s *StationAggregate) Enable() {
	//高内聚逻辑
}

//工作方法还原聚合根
func CreateStationAggregate(cmd *command.CreateStationCmd) *StationAggregate {

	return &StationAggregate{}
}
