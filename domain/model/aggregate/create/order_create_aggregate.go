package create

import (
	"errors"
	"fmt"
	"sort"

	"think.com/go-ddd/domain/common/assert"
	"think.com/go-ddd/domain/model/constant"
	"think.com/go-ddd/domain/model/dp"
	"think.com/go-ddd/domain/model/valueobject"
	"think.com/go-ddd/domain/pl"
	"think.com/go-ddd/domain/pl/command"
)

// OrderCreateAggregate 订单创建聚合根。
//
// 聚合内所有状态变更都通过领域方法完成，外部只能读不能改，
// 保证「下单 → 完善信息 → 校验 → 金额分摊 → 优先级 → 持久化」这条不变式链。
type OrderCreateAggregate struct {
	orderId         *dp.OrderId
	orderStatus     constant.OrderStatus
	desc            string
	orderType       constant.OrderType
	storeInfo       valueobject.StoreInfo
	orderTitle      string
	skuInfos        []*OrderSku
	buyer           valueobject.UserInfo
	orderPay        *OrderPay
	shippingAddress *ShippingAddress
	invoiceInfo     *OrderInvoice
	skuItems        []*OrderSkuItem
	attachInfos     map[string]interface{}
}

// CreateOrderCreateAggregate 由命令创建聚合根（唯一入口）
func CreateOrderCreateAggregate(cmd *command.OrderCreateCommand) (*OrderCreateAggregate, error) {
	if err := assert.NotNull(cmd, "command is null"); err != nil {
		return nil, err
	}
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	aggregate := &OrderCreateAggregate{}
	if err := aggregate.baseCreate(cmd); err != nil {
		return nil, err
	}
	return aggregate, nil
}

// baseCreate 组装聚合根的基础状态
func (a *OrderCreateAggregate) baseCreate(cmd *command.OrderCreateCommand) error {
	a.orderId = dp.NewOrderId(cmd.ExternalOrderNo, cmd.OrderSource)
	a.orderTitle = cmd.OrderTitle
	a.shippingAddress = NewShippingAddress(cmd.Recipient, cmd.Mobile, cmd.Address)
	a.attachInfos = cmd.AttachInfos
	a.buyer = valueobject.NewUserInfo(cmd.UserId, cmd.UserName, cmd.UserType)
	a.invoiceInfo = NewOrderInvoice(cmd.InvoiceName, cmd.InvoiceDetails)
	a.orderType = cmd.OrderType
	a.storeInfo = valueobject.NewStoreInfo(cmd.StoreCode)
	a.orderPay = NewOrderPay(cmd.Currency, cmd.PayType, cmd.OrderPrice, cmd.DiscountPrice, cmd.FeeAmountMap)
	a.orderStatus = cmd.OrderStatus

	skuInfos := make([]*OrderSku, 0, len(cmd.OrderSkuInfos))
	for _, orderSkuInfo := range cmd.OrderSkuInfos {
		orderSku, err := NewOrderSku(orderSkuInfo)
		if err != nil {
			return err
		}
		skuInfos = append(skuInfos, orderSku)
	}
	a.skuInfos = skuInfos
	return nil
}

// Check 订单接入校验（领域方法，高内聚）
func (a *OrderCreateAggregate) Check() error {
	// 只允许已支付订单接入（未支付 / 已取消 / 未定义状态一律拒绝）
	if err := assert.Truef(a.orderStatus == constant.OrderStatusPayed,
		"orderNo=%s, orderStatus=%s, 只允许已支付订单接入", a.orderId.OrderNo(), a.orderStatus); err != nil {
		return err
	}
	if err := assert.NotNull(a.invoiceInfo, "发票信息不正确"); err != nil {
		return err
	}
	if err := assert.True(a.buyer.UserId() != 0, "用户信息不正确"); err != nil {
		return err
	}
	if err := assert.NotBlank(a.shippingAddress.Address(), "下单地址信息不正确"); err != nil {
		return err
	}
	return assert.NotEmpty(len(a.skuInfos), "sku 下单信息为空")
}

// PriceCalculate sku 金额拆分计算（领域行为）。
//
// 优惠金额、运费等附加费用按 sku 行小计的权重比例拆分到每个 sku item，
// 余数由最后一个 item 承担，保证分摊后总额与订单总额严格一致。
func (a *OrderCreateAggregate) PriceCalculate() error {
	if len(a.skuInfos) == 0 {
		return errors.New("skuInfos is empty, can not price calculate")
	}

	var totalWeight int64
	for _, sku := range a.skuInfos {
		totalWeight += sku.TotalPayPrice()
	}

	feeTypes := sortedFeeTypes(a.orderPay.FeeAmountMap())
	remaining := make(map[constant.FeeType]int64, len(feeTypes))
	for _, feeType := range feeTypes {
		remaining[feeType] = a.orderPay.FeeAmountMap()[feeType]
	}

	items := make([]*OrderSkuItem, 0, len(a.skuInfos))
	for idx, sku := range a.skuInfos {
		feeAmountInfos := make(map[constant.FeeType]int64, len(feeTypes))
		for _, feeType := range feeTypes {
			var share int64
			switch {
			case idx == len(a.skuInfos)-1:
				// 最后一个 item 承担全部舍入余数
				share = remaining[feeType]
			case totalWeight > 0:
				share = a.orderPay.FeeAmountMap()[feeType] * sku.TotalPayPrice() / totalWeight
			default:
				share = a.orderPay.FeeAmountMap()[feeType] / int64(len(a.skuInfos))
			}
			remaining[feeType] -= share
			if share != 0 {
				feeAmountInfos[feeType] = share
			}
		}
		items = append(items, NewOrderSkuItem(sku.SkuInfo().SkuCode(), sku.SkuBuyAmount(), sku.SkuPayPrice(), feeAmountInfos))
	}
	a.skuItems = items
	return nil
}

// PriorityProcessing 发货优先级处理（领域行为）。
//
// 生鲜优先发货；虚拟商品无需发货。
func (a *OrderCreateAggregate) PriorityProcessing() error {
	if len(a.skuItems) == 0 {
		return nil
	}
	skuMap := make(map[string]*OrderSku, len(a.skuInfos))
	for _, sku := range a.skuInfos {
		skuMap[sku.SkuInfo().SkuCode()] = sku
	}

	for _, item := range a.skuItems {
		sku, ok := skuMap[item.SkuCode()]
		if !ok {
			return fmt.Errorf("skuCode=%s 查询不到 sku 信息", item.SkuCode())
		}
		priority := defaultShippingPriority
		if sku.SkuCategory() == constant.SkuCategoryFreshFood {
			priority = freshFoodShippingPriority
		}
		if sku.SkuType() == constant.SkuTypeVirtual {
			priority = virtualSkuShippingPriority
		}
		item.PriorityProcessing(priority)
	}
	return nil
}

// ModifyOrderSku 完善 sku 下单信息（外部 sku 转内部 sku、组合商品识别等）
func (a *OrderCreateAggregate) ModifyOrderSku(skuInfoMap map[string]*pl.SkuFullInfo) error {
	for _, orderSku := range a.skuInfos {
		externalSkuId := orderSku.SkuInfo().ExternalSkuId()
		skuFullInfo, ok := skuInfoMap[externalSkuId]
		if !ok || skuFullInfo == nil {
			return fmt.Errorf("根据外部 skuId=[%s] 查询不到商品信息", externalSkuId)
		}
		if err := orderSku.ModifySku(skuFullInfo); err != nil {
			return err
		}
	}
	return nil
}

// Hangup 挂起订单（风控不通过时走业务补偿，而不是直接丢弃）
func (a *OrderCreateAggregate) Hangup(desc string) {
	a.orderStatus = constant.OrderStatusHangUp
	a.desc = desc
}

// OrderId 订单标识
func (a *OrderCreateAggregate) OrderId() *dp.OrderId { return a.orderId }

// OrderStatus 订单状态
func (a *OrderCreateAggregate) OrderStatus() constant.OrderStatus { return a.orderStatus }

// Desc 备注 / 挂起原因
func (a *OrderCreateAggregate) Desc() string { return a.desc }

// OrderType 订单类型
func (a *OrderCreateAggregate) OrderType() constant.OrderType { return a.orderType }

// StoreInfo 店铺信息
func (a *OrderCreateAggregate) StoreInfo() valueobject.StoreInfo { return a.storeInfo }

// OrderTitle 订单标题
func (a *OrderCreateAggregate) OrderTitle() string { return a.orderTitle }

// SkuInfos sku 下单信息
func (a *OrderCreateAggregate) SkuInfos() []*OrderSku { return a.skuInfos }

// Buyer 下单用户
func (a *OrderCreateAggregate) Buyer() valueobject.UserInfo { return a.buyer }

// OrderPay 支付信息
func (a *OrderCreateAggregate) OrderPay() *OrderPay { return a.orderPay }

// ShippingAddress 收货地址
func (a *OrderCreateAggregate) ShippingAddress() *ShippingAddress { return a.shippingAddress }

// InvoiceInfo 发票信息
func (a *OrderCreateAggregate) InvoiceInfo() *OrderInvoice { return a.invoiceInfo }

// SkuItems sku item 级下单信息
func (a *OrderCreateAggregate) SkuItems() []*OrderSkuItem { return a.skuItems }

// AttachInfos 附加信息
func (a *OrderCreateAggregate) AttachInfos() map[string]interface{} { return a.attachInfos }

// 发货优先级常量：0 无需发货，数字越大优先级越高
const (
	virtualSkuShippingPriority = 0
	defaultShippingPriority    = 1
	freshFoodShippingPriority  = 3
)

// sortedFeeTypes 对费用类型排序，保证分摊结果稳定可复现
func sortedFeeTypes(feeAmountMap map[constant.FeeType]int64) []constant.FeeType {
	feeTypes := make([]constant.FeeType, 0, len(feeAmountMap))
	for feeType := range feeAmountMap {
		feeTypes = append(feeTypes, feeType)
	}
	sort.Slice(feeTypes, func(i, j int) bool { return feeTypes[i] < feeTypes[j] })
	return feeTypes
}
