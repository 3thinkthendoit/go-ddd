package create

import "fmt"

// ShippingAddress 收货地址（值对象）
type ShippingAddress struct {
	recipient   string
	contactInfo string
	addressCode int
	address     string
}

// NewShippingAddress 下单时创建，此时还没有标准地址编码
func NewShippingAddress(recipient, contactInfo, address string) *ShippingAddress {
	return &ShippingAddress{recipient: recipient, contactInfo: contactInfo, address: address}
}

// NewShippingAddressWithCode 由标准地址编码还原
func NewShippingAddressWithCode(addressCode int, address string) *ShippingAddress {
	return &ShippingAddress{addressCode: addressCode, address: address}
}

// Recipient 收件人
func (a *ShippingAddress) Recipient() string { return a.recipient }

// ContactInfo 联系方式
func (a *ShippingAddress) ContactInfo() string { return a.contactInfo }

// AddressCode 标准地址编码
func (a *ShippingAddress) AddressCode() int { return a.addressCode }

// Address 详细地址
func (a *ShippingAddress) Address() string { return a.address }

// InitAddressCode 基础信息域回填标准地址编码
func (a *ShippingAddress) InitAddressCode(addressCode int) { a.addressCode = addressCode }

func (a *ShippingAddress) String() string {
	return fmt.Sprintf("ShippingAddress{recipient=%s, contactInfo=%s, address=%s}",
		a.recipient, a.contactInfo, a.address)
}
