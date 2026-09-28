// Package util 通用工具类
package util

import (
	"encoding/json"
	"log"
)

// ToJson 对象转 JSON 字符串，失败时返回空串并打日志
func ToJson(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("[util] json marshal failed: %v", err)
		return ""
	}
	return string(data)
}

// FromJson JSON 字符串转对象
func FromJson(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}
