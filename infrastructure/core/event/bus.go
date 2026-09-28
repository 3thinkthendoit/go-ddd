// Package event 进程内事件总线。
//
// 等价于 Spring 的 ApplicationEventPublisher + @EventListener + @Async：
// 发布方只管 Publish，订阅方决定同步或异步消费，双方通过事件名解耦。
package event

import (
	"context"
	"log"
	"sync"

	"github.com/go-spring/spring-core/gs"
)

// Event 可被总线分发的事件
type Event interface {
	// EventName 事件名，订阅方按事件名匹配
	EventName() string
}

// Handler 事件处理器
type Handler func(Event)

type entry struct {
	handler Handler
	async   bool
}

func init() {
	// 事件总线同时作为应用生命周期事件，停机时等待异步处理器跑完
	gs.Object(NewBus()).Export((*gs.AppEvent)(nil))
}

// Bus 进程内事件总线
type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]entry
	wg       sync.WaitGroup
}

// NewBus 创建事件总线
func NewBus() *Bus {
	return &Bus{handlers: map[string][]entry{}}
}

// OnAppStart 应用启动事件（事件总线无需额外处理）
func (b *Bus) OnAppStart(_ gs.Context) {}

// OnAppStop 容器停止时等待所有异步处理器执行完成
func (b *Bus) OnAppStop(_ context.Context) {
	b.Wait()
}

// Subscribe 同步订阅：事件在当前 goroutine 内处理完才返回
func (b *Bus) Subscribe(eventName string, handler Handler) {
	b.add(eventName, entry{handler: handler, async: false})
}

// SubscribeAsync 异步订阅：对应 Spring 的 @Async，事件在独立 goroutine 中处理
func (b *Bus) SubscribeAsync(eventName string, handler Handler) {
	b.add(eventName, entry{handler: handler, async: true})
}

func (b *Bus) add(eventName string, e entry) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.handlers == nil {
		b.handlers = map[string][]entry{}
	}
	b.handlers[eventName] = append(b.handlers[eventName], e)
}

// Publish 发布事件。
//
// 同步订阅者立即执行，异步订阅者交给 goroutine，单个处理器 panic 不影响其他处理器。
func (b *Bus) Publish(evt Event) {
	b.mu.RLock()
	handlers := make([]entry, len(b.handlers[evt.EventName()]))
	copy(handlers, b.handlers[evt.EventName()])
	b.mu.RUnlock()

	for _, e := range handlers {
		if e.async {
			b.wg.Add(1)
			go func(h Handler) {
				defer b.wg.Done()
				safeHandle(h, evt)
			}(e.handler)
			continue
		}
		safeHandle(e.handler, evt)
	}
}

// Wait 等待所有异步处理器执行完成（优雅停机 / 测试用）
func (b *Bus) Wait() {
	b.wg.Wait()
}

func safeHandle(h Handler, evt Event) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[event] 处理事件 %s 时发生 panic: %v", evt.EventName(), r)
		}
	}()
	h(evt)
}
