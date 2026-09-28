// Package resp 北向网关返回对象（DTO）。
package resp

// Result 统一返回结构
type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功返回
func Success(data interface{}) *Result {
	return &Result{Code: 0, Msg: "success", Data: data}
}

// Failure 失败返回
func Failure(msg string) *Result {
	return &Result{Code: 1, Msg: msg}
}
