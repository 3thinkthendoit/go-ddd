package pl

import "think.com/go-ddd/domain/model/constant"

// OrderInfo 订单查询结果（对外统一订单视图）
type OrderInfo struct {
	OrderNo         string
	ExternalOrderNo string
	OrderSource     constant.OrderSource
	OrderStatus     constant.OrderStatus
	SkuInfos        []*SkuItemInfo
}
