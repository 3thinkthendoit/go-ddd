package create

import (
	"fmt"

	"think.com/go-ddd/domain/common/assert"
	"think.com/go-ddd/domain/model/constant"
	"think.com/go-ddd/domain/model/valueobject"
	"think.com/go-ddd/domain/pl"
)

// OrderSku 订单内的 sku 下单实体。
//
// 下单时只知道外部 sku 标识，商品域信息（内部 skuCode、价格、类型）由领域服务
// 通过 ModifySku 补全。
type OrderSku struct {
	skuInfo      valueobject.SkuInfo
	skuPayPrice  int64
	skuBuyAmount int
	skuType      constant.SkuType
	skuCategory  constant.SkuCategory
}

// NewOrderSku 由外部 sku 下单信息创建实体
func NewOrderSku(orderSkuInfo *pl.OrderSkuInfo) (*OrderSku, error) {
	if err := assert.NotNull(orderSkuInfo, "orderSkuInfo is null"); err != nil {
		return nil, err
	}
	skuInfo, err := valueobject.CreateSkuInfo(orderSkuInfo.ExternalSkuId, orderSkuInfo.ExternalSkuCode)
	if err != nil {
		return nil, err
	}
	if err := assert.True(orderSkuInfo.SkuBuyAmount > 0, "skuBuyAmount is invalid"); err != nil {
		return nil, err
	}
	return &OrderSku{
		skuInfo:      skuInfo,
		skuPayPrice:  orderSkuInfo.SkuPayPrice,
		skuBuyAmount: orderSkuInfo.SkuBuyAmount,
	}, nil
}

// ModifySku 用商品域的完整信息补全自身（领域方法）
func (o *OrderSku) ModifySku(skuFullInfo *pl.SkuFullInfo) error {
	if err := o.skuInfo.Init(skuFullInfo); err != nil {
		return err
	}
	o.skuType = o.skuInfo.SkuType()
	o.skuCategory = o.skuInfo.SkuCategory()
	return nil
}

// SkuInfo 商品基本信息
func (o *OrderSku) SkuInfo() valueobject.SkuInfo { return o.skuInfo }

// SkuPayPrice sku 实付价格（分）
func (o *OrderSku) SkuPayPrice() int64 { return o.skuPayPrice }

// SkuBuyAmount 下单数量
func (o *OrderSku) SkuBuyAmount() int { return o.skuBuyAmount }

// SkuType 商品类型
func (o *OrderSku) SkuType() constant.SkuType { return o.skuType }

// SkuCategory 商品种类
func (o *OrderSku) SkuCategory() constant.SkuCategory { return o.skuCategory }

// TotalPayPrice sku 行小计金额（分）
func (o *OrderSku) TotalPayPrice() int64 { return o.skuPayPrice * int64(o.skuBuyAmount) }

func (o *OrderSku) String() string {
	return fmt.Sprintf("OrderSku{skuCode=%s, amount=%d, payPrice=%d}",
		o.skuInfo.SkuCode(), o.skuBuyAmount, o.skuPayPrice)
}
