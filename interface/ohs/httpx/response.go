// Package httpx 北向 HTTP 适配器的公共辅助。
//
// Spring 里由 @ResponseBody + 消息转换器自动完成的部分，在 Go 里显式封装一层，
// 让 controller / rpc 只关心业务参数与返回对象。
package httpx

import (
	"encoding/json"
	"log"

	"github.com/go-spring/spring-core/web"

	"think.com/go-ddd/interface/ohs/dto/resp"
)

// Write 输出 JSON 响应
func Write(ctx web.Context, body interface{}) {
	ctx.SetContentType(web.MIMEApplicationJSONCharsetUTF8)
	if err := json.NewEncoder(ctx.Response()).Encode(body); err != nil {
		log.Printf("[httpx] 写响应失败: %v", err)
	}
}

// WriteSuccess 输出成功响应
func WriteSuccess(ctx web.Context, data interface{}) {
	Write(ctx, resp.Success(data))
}

// WriteFailure 输出失败响应（HTTP 状态码保持 200，由业务码区分）
func WriteFailure(ctx web.Context, err error) {
	Write(ctx, resp.Failure(err.Error()))
}

// BindJSON 读取请求体并反序列化到目标对象
func BindJSON(ctx web.Context, target interface{}) error {
	body, err := ctx.RequestBody()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}
