package job

import (
	"log"
	"time"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/infrastructure/acl/api/pdd"
	"think.com/go-ddd/infrastructure/core/scheduler"
	"think.com/go-ddd/interface/local"
)

func init() {
	gs.Object(new(PinDuoDuoOrderTask)).Init(func(t *PinDuoDuoOrderTask) {
		t.scheduler.Every("PinDuoDuoOrderTask", 5*time.Minute, t.PullOrder)
	})
}

// PinDuoDuoOrderTask 拼多多订单定时拉取任务
type PinDuoDuoOrderTask struct {
	orderLocalService *local.OrderLocalService `autowire:""`
	pddClient         *pdd.Client              `autowire:""`
	scheduler         *scheduler.Scheduler     `autowire:""`
}

// PullOrder 拉取最近一天的订单
func (t *PinDuoDuoOrderTask) PullOrder() {
	begin := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	cmd, err := t.pddClient.PullOrder(begin, end)
	if err != nil {
		log.Printf("[job] 拼多多订单拉取失败: %v", err)
		return
	}
	if cmd == nil {
		return
	}
	if err := t.orderLocalService.CreateOrder(cmd); err != nil {
		log.Printf("[job] 拼多多订单接入失败 externalOrderNo=%s: %v", cmd.ExternalOrderNo, err)
	}
}
