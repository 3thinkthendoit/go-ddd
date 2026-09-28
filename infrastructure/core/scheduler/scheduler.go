// Package scheduler 定时任务调度组件。
//
// 等价于 Spring 的 @Scheduled(fixedRate=...)：接口层的 job 只声明
// 「多久跑一次 + 跑什么」，启动/停止由容器生命周期驱动。
package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-spring/spring-core/gs"
)

func init() {
	// 调度器跟随容器生命周期启动 / 停止
	gs.Object(new(Scheduler)).Export((*gs.AppEvent)(nil))
}

type task struct {
	name     string
	interval time.Duration
	fn       func()
}

// Scheduler 固定频率调度器
type Scheduler struct {
	Enabled bool `value:"${oms.scheduler.enabled:=false}"`

	mu      sync.Mutex
	tasks   []task
	stopCh  chan struct{}
	started bool
	wg      sync.WaitGroup
}

// NewScheduler 创建调度器
func NewScheduler() *Scheduler {
	return &Scheduler{}
}

// OnAppStart 容器启动完成后统一启动任务（等价 Spring 的 @EnableScheduling 生效时机）
func (s *Scheduler) OnAppStart(_ gs.Context) {
	s.Start()
}

// OnAppStop 容器停止时优雅关闭
func (s *Scheduler) OnAppStop(_ context.Context) {
	_ = s.Stop()
}

// Every 注册一个固定频率任务
func (s *Scheduler) Every(name string, interval time.Duration, fn func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = append(s.tasks, task{name: name, interval: interval, fn: fn})
}

// Start 启动所有已注册任务
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.started || !s.Enabled {
		s.mu.Unlock()
		if !s.Enabled {
			log.Printf("[scheduler] 未开启（oms.scheduler.enabled=false），已注册任务不会执行")
		}
		return
	}
	s.started = true
	s.stopCh = make(chan struct{})
	tasks := make([]task, len(s.tasks))
	copy(tasks, s.tasks)
	s.mu.Unlock()

	for _, t := range tasks {
		s.wg.Add(1)
		go func(t task) {
			defer s.wg.Done()
			ticker := time.NewTicker(t.interval)
			defer ticker.Stop()
			log.Printf("[scheduler] 启动任务 %s，周期 %s", t.name, t.interval)
			for {
				select {
				case <-s.stopCh:
					return
				case <-ticker.C:
					s.run(t)
				}
			}
		}(t)
	}
}

// Stop 停止调度器
func (s *Scheduler) Stop() error {
	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return nil
	}
	s.started = false
	close(s.stopCh)
	s.mu.Unlock()
	s.wg.Wait()
	log.Printf("[scheduler] 已停止")
	return nil
}

func (s *Scheduler) run(t task) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[scheduler] 任务 %s 执行 panic: %v", t.name, r)
		}
	}()
	t.fn()
}
