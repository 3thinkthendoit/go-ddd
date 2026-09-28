package gateway

// InvoiceGateway 发票域南向网关
type InvoiceGateway interface {
	// Issue 开票
	Issue(orderNo string) error
}
