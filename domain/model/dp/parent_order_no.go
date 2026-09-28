package dp

import (
	"fmt"
	"sync/atomic"
	"time"
)

// parentOrderNoSeq 同一毫秒内拆单父单号自增序列。
var parentOrderNoSeq uint32

// NewParentOrderNo 生成拆单后的父单号（发货单号）。
//
// 与订单号一样属于标识类概念，收敛在 DP 层统一生成，避免各处自行拼字符串。
func NewParentOrderNo() string {
	seq := atomic.AddUint32(&parentOrderNoSeq, 1) % 10000
	return fmt.Sprintf("10%d%04d", time.Now().UnixMilli(), seq)
}
