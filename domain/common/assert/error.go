package assert

import "fmt"

// fmtErrorf 独立出来便于后续替换为带错误码的错误类型。
func fmtErrorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
