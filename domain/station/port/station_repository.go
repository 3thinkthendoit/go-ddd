package port

import "think.com/go-ddd/domain/station/model/aggregate"

//资源库接口
type StationRepository interface {
	//SaveAllStation([]*aggregate.StationAggregate)
	Save(a *aggregate.StationAggregate)
}
