package command

import (
	"think.com/go-ddd/domain/common/assert"
	"think.com/go-ddd/domain/model/constant"
	"think.com/go-ddd/domain/model/valueobject"
	"think.com/go-ddd/domain/pl"
)

// OrderCreateCommand 创建订单命令。
//
// 由北向网关（HTTP / RPC / MQ / 定时任务）解析外部协议后组装，
// 是外部世界进入订单领域的唯一入口。
type OrderCreateCommand struct {
	// 外部订单号
	ExternalOrderNo string
	// 订单状态
	OrderStatus constant.OrderStatus
	// 订单类型
	OrderType constant.OrderType
	// 店铺编码
	StoreCode string
	// 订单来源
	OrderSource constant.OrderSource
	// 订单标题
	OrderTitle string
	// 订单金额（单位：分）
	OrderPrice int64
	// 支付方式
	PayType constant.PayType
	// 优惠金额（单位：分）
	DiscountPrice int64
	// 币种
	Currency constant.Currency
	// Sku 下单信息
	OrderSkuInfos []*pl.OrderSkuInfo
	// 下单用户信息
	UserId   int64
	UserName string
	UserType valueobject.UserType
	// 收货地址
	Address   string
	Mobile    string
	Recipient string
	// 发票信息
	InvoiceName    string
	InvoiceDetails string
	// 附加费用信息
	FeeAmountMap map[constant.FeeType]int64
	// 附加信息
	AttachInfos map[string]interface{}
}

// Validate 校验命令必填项
func (c *OrderCreateCommand) Validate() error {
	if err := assert.NotBlank(c.ExternalOrderNo, "externalOrderNo is blank"); err != nil {
		return err
	}
	if err := assert.True(c.OrderSource != constant.OrderSourceUnKnow, "orderSource is unknown"); err != nil {
		return err
	}
	if err := assert.True(c.OrderPrice > 0, "orderPrice is invalid"); err != nil {
		return err
	}
	if err := assert.NotBlank(c.StoreCode, "storeCode is blank"); err != nil {
		return err
	}
	if err := assert.NotEmpty(len(c.OrderSkuInfos), "orderSkuInfos is empty"); err != nil {
		return err
	}
	return assert.NotBlank(c.Recipient, "recipient is blank")
}
