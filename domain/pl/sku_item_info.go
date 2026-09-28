package pl

import "think.com/go-ddd/domain/model/constant"

// SkuItemInfo sku item 级别的下单 / 发货 / 售后 / 签收信息。
//
// 拆单、履约、发货回传都围绕 item 维度展开，是 OMS 内部的核心流转对象。
type SkuItemInfo struct {
	SkuFullInfo    *SkuFullInfo
	SkuPayPrice    int64
	SkuAmount      int
	Priority       int
	ShippingAmount int
	ReturnAmount   int
	ConfirmAmount  int
	SkuType        constant.SkuType
	SubSkuList     []*SkuItemInfo
}
