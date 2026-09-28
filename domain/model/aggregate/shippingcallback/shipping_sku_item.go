package shippingcallback

import "fmt"

// ShippingSkuItem 发货回传聚合内的 sku item 实体。
type ShippingSkuItem struct {
	skuCode        string
	externalSkuId  string
	skuAmount      int
	shippingAmount int
	expressCode    string
	expressNo      string
}

// NewShippingSkuItem 创建 sku item。
//
// shippingAmount 需要带入订单上已有的发货数量，避免本次回传只覆盖部分 sku 时，
// 把其他 sku 的已发货数量冲成 0。
func NewShippingSkuItem(skuCode, externalSkuId string, skuAmount, shippingAmount int) *ShippingSkuItem {
	return &ShippingSkuItem{
		skuCode:        skuCode,
		externalSkuId:  externalSkuId,
		skuAmount:      skuAmount,
		shippingAmount: shippingAmount,
	}
}

// AddShippingInfo 累加一次 WMS 回传的发货信息（领域方法）。
//
// WMS 回传的 shippingAmount 是**本次发货数量**（增量），不是累计已发数量：
// 同一个 sku 可能被拆到多张发货单、或一张发货单分多次发出，因此这里必须累加。
// 直接赋值会把之前的发货量冲掉，导致部分发货永远凑不满、主订单无法流转为「已发货」。
func (s *ShippingSkuItem) AddShippingInfo(shippingAmount int, expressCode, expressNo string) {
	s.shippingAmount += shippingAmount
	s.expressCode = expressCode
	s.expressNo = expressNo
}

// SkuCode 内部 skuCode
func (s *ShippingSkuItem) SkuCode() string { return s.skuCode }

// ExternalSkuId 外部 sku 标识
func (s *ShippingSkuItem) ExternalSkuId() string { return s.externalSkuId }

// SkuAmount 下单数量
func (s *ShippingSkuItem) SkuAmount() int { return s.skuAmount }

// ShippingAmount 已发货数量
func (s *ShippingSkuItem) ShippingAmount() int { return s.shippingAmount }

// ExpressCode 快递公司编码
func (s *ShippingSkuItem) ExpressCode() string { return s.expressCode }

// ExpressNo 快递单号
func (s *ShippingSkuItem) ExpressNo() string { return s.expressNo }

func (s *ShippingSkuItem) String() string {
	return fmt.Sprintf("ShippingSkuItem{skuCode=%s, skuAmount=%d, shippingAmount=%d, expressNo=%s}",
		s.skuCode, s.skuAmount, s.shippingAmount, s.expressNo)
}
