// Package http HTTP 客户端组件。
//
// 南向网关调用外部系统（商品中心、WMS、风控等）统一走这里，
// 便于集中管理超时、重试、日志与鉴权。
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-spring/spring-core/gs"
)

func init() {
	// 配置绑定完成后初始化底层 http.Client
	gs.Object(new(Client)).Init(func(c *Client) { c.initClient() })
}

// Client HTTP 客户端
type Client struct {
	client  *http.Client
	Timeout int `value:"${oms.http.timeout-ms:=5000}"`
}

// initClient 用配置好的超时初始化底层客户端
func (c *Client) initClient() {
	timeout := time.Duration(c.Timeout) * time.Millisecond
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	c.client = &http.Client{Timeout: timeout}
}

// Get 发起 GET 请求
func (c *Client) Get(url string) (string, error) {
	return c.do(http.MethodGet, url, nil)
}

// Post 发起 JSON POST 请求
func (c *Client) Post(url string, params map[string]interface{}) (string, error) {
	return c.do(http.MethodPost, url, params)
}

func (c *Client) do(method, url string, params map[string]interface{}) (string, error) {
	if c.client == nil {
		c.initClient()
	}
	var body io.Reader
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return "", err
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}
	if params != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return string(data), fmt.Errorf("http %s %s 返回状态码 %d", method, url, resp.StatusCode)
	}
	return string(data), nil
}
