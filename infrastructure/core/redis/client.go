// Package redis 缓存组件。
//
// 给出最小可用的接口形态，真实项目替换为 go-redis 即可，调用方不感知。
package redis

import (
	"sync"
	"time"

	"github.com/go-spring/spring-core/gs"
)

func init() {
	gs.Object(NewClient())
}

// Client 缓存客户端
type Client struct {
	mu   sync.RWMutex
	data map[string]entry
}

type entry struct {
	value     string
	expiredAt time.Time
}

// NewClient 创建缓存客户端
func NewClient() *Client {
	return &Client{data: map[string]entry{}}
}

// Set 写入缓存
func (c *Client) Set(key, value string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var expiredAt time.Time
	if ttl > 0 {
		expiredAt = time.Now().Add(ttl)
	}
	c.data[key] = entry{value: value, expiredAt: expiredAt}
	return nil
}

// Get 读取缓存，第二个返回值表示是否命中
func (c *Client) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.data[key]
	if !ok {
		return "", false
	}
	if !e.expiredAt.IsZero() && time.Now().After(e.expiredAt) {
		return "", false
	}
	return e.value, true
}

// Del 删除缓存
func (c *Client) Del(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}
