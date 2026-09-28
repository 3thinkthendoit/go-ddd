package dp

import (
	"fmt"
	"sync/atomic"
	"time"

	"think.com/go-ddd/domain/model/constant"
)

// orderNoSeq 同一毫秒内的订单号自增序列，避免并发下单时订单号冲突。
var orderNoSeq uint32

// nextOrderNo 生成订单号：毫秒时间戳 + 4 位自增序列。
func nextOrderNo() string {
	seq := atomic.AddUint32(&orderNoSeq, 1) % 10000
	return fmt.Sprintf("%d%04d", time.Now().UnixMilli(), seq)
}

// OrderId 订单标识（DP：Domain Primitive）。
//
// 把「订单号 + 外部订单号 + 订单来源」这三个必须一起出现、且各自都不足以定位订单
// 的属性收敛成一个整体，避免在领域内部到处传递裸 string。
type OrderId struct {
	orderNo         string
	externalOrderNo string
	orderSource     constant.OrderSource
}

// NewOrderId 由外部订单号 + 订单来源创建，内部生成 OMS 订单号。
func NewOrderId(externalOrderNo string, orderSource constant.OrderSource) *OrderId {
	return &OrderId{
		orderNo:         nextOrderNo(),
		externalOrderNo: externalOrderNo,
		orderSource:     orderSource,
	}
}

// NewOrderIdByOrderNo 由已有的 OMS 订单号还原 OrderId。
func NewOrderIdByOrderNo(orderNo string) *OrderId {
	return &OrderId{orderNo: orderNo}
}

// OrderNo OMS 订单号
func (o *OrderId) OrderNo() string { return o.orderNo }

// ExternalOrderNo 外部平台订单号
func (o *OrderId) ExternalOrderNo() string { return o.externalOrderNo }

// OrderSource 订单来源
func (o *OrderId) OrderSource() constant.OrderSource { return o.orderSource }

// Init 补全订单来源与外部订单号。
//
// 发货回传等场景只拿到 OMS 订单号，需要先查询订单再补全其余标识。
func (o *OrderId) Init(externalOrderNo string, orderSource constant.OrderSource) {
	o.externalOrderNo = externalOrderNo
	o.orderSource = orderSource
}

func (o *OrderId) String() string {
	return fmt.Sprintf("OrderId{orderNo=%s, externalOrderNo=%s, orderSource=%s}",
		o.orderNo, o.externalOrderNo, o.orderSource)
}
