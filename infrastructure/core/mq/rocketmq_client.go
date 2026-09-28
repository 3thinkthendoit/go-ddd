// Package mq 消息队列客户端。
//
// 这里给出的是进程内实现（按 topic 分发），用于让「发布 MQ 事件 → 消费者处理」
// 这条链路可以本地跑通；接入真实 RocketMQ 时只需替换本文件，端口与调用方不变。
package mq

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-spring/spring-core/gs"
)

func init() {
	gs.Object(NewRocketMqClient())
}

// Handler 消息处理器
type Handler func(topic string, message string) error

type consumer struct {
	group   string
	handler Handler
}

// RocketMqClient MQ 客户端
type RocketMqClient struct {
	mu        sync.RWMutex
	consumers map[string][]consumer
}

// NewRocketMqClient 创建 MQ 客户端
func NewRocketMqClient() *RocketMqClient {
	return &RocketMqClient{consumers: map[string][]consumer{}}
}

// Send 发送消息到指定 topic
func (c *RocketMqClient) Send(topic, message string) error {
	c.mu.RLock()
	consumers := make([]consumer, len(c.consumers[topic]))
	copy(consumers, c.consumers[topic])
	c.mu.RUnlock()

	for _, item := range consumers {
		go func(item consumer) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[mq] topic=%s group=%s 消费 panic: %v", topic, item.group, r)
				}
			}()
			if err := item.handler(topic, message); err != nil {
				log.Printf("[mq] topic=%s group=%s 消费失败: %v", topic, item.group, err)
			}
		}(item)
	}
	log.Printf("[mq] send topic=%s message=%s", topic, message)
	return nil
}

// Subscribe 订阅 topic
func (c *RocketMqClient) Subscribe(topic, consumerGroup string, handler Handler) error {
	if handler == nil {
		return fmt.Errorf("handler is nil, topic=%s", topic)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.consumers[topic] = append(c.consumers[topic], consumer{group: consumerGroup, handler: handler})
	return nil
}
