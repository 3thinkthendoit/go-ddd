// Package assert 提供领域层共用的前置条件 / 不变式断言。
//
// 领域模型里的校验统一走这里并返回 error，而不是 panic：
// 由应用层决定是中断用例、还是把订单挂起走业务补偿。
package assert

import (
	"errors"
	"reflect"
)

// NotNull 断言非空（兼容接口中包装的 typed nil 指针）
func NotNull(v interface{}, msg string) error {
	if isNilValue(v) {
		return errors.New(msg)
	}
	return nil
}

// IsNil 断言为空
func IsNil(v interface{}, msg string) error {
	if !isNilValue(v) {
		return errors.New(msg)
	}
	return nil
}

// NotBlank 断言字符串非空白
func NotBlank(v, msg string) error {
	if v == "" {
		return errors.New(msg)
	}
	return nil
}

// True 断言条件成立
func True(ok bool, msg string) error {
	if !ok {
		return errors.New(msg)
	}
	return nil
}

// Truef 断言条件成立，支持格式化消息
func Truef(ok bool, format string, args ...interface{}) error {
	if !ok {
		return fmtErrorf(format, args...)
	}
	return nil
}

// NotEmpty 断言集合非空
func NotEmpty(size int, msg string) error {
	if size == 0 {
		return errors.New(msg)
	}
	return nil
}

func isNilValue(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}
