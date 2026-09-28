// Package demo 本地演示数据。
//
// 仅用于让示例项目在没有外部依赖时也能跑通「接单 → 分仓拆单 → 推 WMS → 发货回传」
// 全链路；接入真实环境时删除本包即可，业务代码不依赖它。
package demo

import (
	"fmt"
	"time"

	"think.com/go-ddd/domain/model/constant"
	"think.com/go-ddd/domain/model/valueobject"
	"think.com/go-ddd/domain/pl"
	"think.com/go-ddd/domain/pl/command"
)

// BuildOrderCommand 构造一个示例创建订单命令
func BuildOrderCommand(prefix string, source constant.OrderSource,
	payType constant.PayType, storeCode, title string) *command.OrderCreateCommand {
	now := time.Now()
	return &command.OrderCreateCommand{
		ExternalOrderNo: fmt.Sprintf("%s%s", prefix, now.Format("20060102150405")),
		OrderStatus:     constant.OrderStatusPayed,
		OrderType:       constant.OrderTypeStandard,
		StoreCode:       storeCode,
		OrderSource:     source,
		OrderTitle:      title,
		OrderPrice:      10196,
		PayType:         payType,
		DiscountPrice:   500,
		Currency:        constant.CurrencyCny,
		UserId:          now.UnixMilli() % 100000,
		UserName:        "demo-buyer",
		UserType:        valueobject.UserTypeToC,
		Recipient:       "张三",
		Mobile:          "13800000000",
		Address:         "浙江省杭州市余杭区文一西路 969 号",
		InvoiceName:     "张三",
		InvoiceDetails:  "电子普通发票",
		FeeAmountMap: map[constant.FeeType]int64{
			constant.FeeTypeTranFee:     1200,
			constant.FeeTypeDiscountFee: 500,
		},
		AttachInfos: map[string]interface{}{"channel": prefix},
		OrderSkuInfos: []*pl.OrderSkuInfo{
			// 生鲜：优先发货
			{ExternalSkuId: "EXT_SKU_1001", ExternalSkuCode: "EXT-1001", SkuName: "阳澄湖大闸蟹", SkuPayPrice: 2999, SkuBuyAmount: 2},
			// 百货单品
			{ExternalSkuId: "EXT_SKU_1002", ExternalSkuCode: "EXT-1002", SkuName: "保温杯", SkuPayPrice: 1599, SkuBuyAmount: 1},
			// 虚拟商品：无需发货
			{ExternalSkuId: "EXT_SKU_1003", ExternalSkuCode: "EXT-1003", SkuName: "视频会员年卡", SkuPayPrice: 990, SkuBuyAmount: 1},
			// 组合商品：被指定从广州仓发货
			{ExternalSkuId: "EXT_SKU_1004", ExternalSkuCode: "EXT-1004", SkuName: "厨房三件套", SkuPayPrice: 5999, SkuBuyAmount: 1},
		},
	}
}
