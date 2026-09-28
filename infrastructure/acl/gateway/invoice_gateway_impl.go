package gateway

import (
	"log"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/domain/port/gateway"
)

func init() {
	gs.Object(new(InvoiceGatewayImpl)).Export((*gateway.InvoiceGateway)(nil))
}

// InvoiceGatewayImpl 发票域南向网关实现。
type InvoiceGatewayImpl struct {
	// 是否使用内置模拟
	MockEnabled bool `value:"${oms.invoice.mock:=true}"`
	// 开票请求的 MQ topic
	Topic string `value:"${oms.invoice.topic:=oms-invoice-issue}"`
}

// Issue 开票
func (g *InvoiceGatewayImpl) Issue(orderNo string) error {
	// 真实实现：通过 MQ 通知发票域开票
	log.Printf("[invoice] 发送开票请求 orderNo=%s topic=%s", orderNo, g.Topic)
	return nil
}
