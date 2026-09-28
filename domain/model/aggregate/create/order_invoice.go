package create

import "fmt"

// OrderInvoice 发票信息（值对象）
type OrderInvoice struct {
	invoiceName    string
	invoiceDetails string
}

// NewOrderInvoice 创建发票信息
func NewOrderInvoice(invoiceName, invoiceDetails string) *OrderInvoice {
	return &OrderInvoice{invoiceName: invoiceName, invoiceDetails: invoiceDetails}
}

// InvoiceName 发票抬头
func (i *OrderInvoice) InvoiceName() string { return i.invoiceName }

// InvoiceDetails 发票明细
func (i *OrderInvoice) InvoiceDetails() string { return i.invoiceDetails }

func (i *OrderInvoice) String() string {
	return fmt.Sprintf("OrderInvoice{invoiceName=%s}", i.invoiceName)
}
