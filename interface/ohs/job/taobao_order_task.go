// Package job 定时任务（北向网关）。
//
// 对应 Spring 的 @Scheduled(fixedRate=...)：按固定频率从外部平台拉取订单，
// 转成创建订单命令后交给本地服务。
package job

import (
	"log"
	"time"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/infrastructure/acl/api/taobao"
	"think.com/go-ddd/infrastructure/core/scheduler"
	"think.com/go-ddd/interface/local"
)

func init() {
	gs.Object(new(TaoBaoOrderTask)).Init(func(t *TaoBaoOrderTask) {
		t.scheduler.Every("TaoBaoOrderTask", 5*time.Minute, t.PullOrder)
	})
}

// TaoBaoOrderTask 淘宝订单定时拉取任务
type TaoBaoOrderTask struct {
	orderLocalService *local.OrderLocalService `autowire:""`
	taoBaoClient      *taobao.Client           `autowire:""`
	scheduler         *scheduler.Scheduler     `autowire:""`
}

// PullOrder 拉取最近一天的订单
func (t *TaoBaoOrderTask) PullOrder() {
	begin := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	cmd, err := t.taoBaoClient.PullOrder(begin, end)
	if err != nil {
		log.Printf("[job] 淘宝订单拉取失败: %v", err)
		return
	}
	if cmd == nil {
		return
	}
	if err := t.orderLocalService.CreateOrder(cmd); err != nil {
		log.Printf("[job] 淘宝订单接入失败 externalOrderNo=%s: %v", cmd.ExternalOrderNo, err)
	}
}
