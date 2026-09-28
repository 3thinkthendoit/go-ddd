package resp

import "think.com/go-ddd/domain/pl"

// SkuItemInfo sku item 级信息
type SkuItemInfo struct {
	SkuCode        string `json:"skuCode"`
	SkuName        string `json:"skuName"`
	SkuAmount      int    `json:"skuAmount"`
	Priority       int    `json:"priority"`
	ShippingAmount int    `json:"shippingAmount"`
}

// OrderInfoQueryResp 订单查询返回
type OrderInfoQueryResp struct {
	OrderNo         string        `json:"orderNo"`
	ExternalOrderNo string        `json:"externalOrderNo"`
	OrderSource     int           `json:"orderSource"`
	OrderSourceDesc string        `json:"orderSourceDesc"`
	OrderStatus     int           `json:"orderStatus"`
	OrderStatusDesc string        `json:"orderStatusDesc"`
	SkuItemInfos    []SkuItemInfo `json:"skuItemInfos"`
}

// NewOrderInfoQueryResp 领域订单视图 → 北向 DTO
func NewOrderInfoQueryResp(order *pl.OrderInfo) OrderInfoQueryResp {
	if order == nil {
		return OrderInfoQueryResp{}
	}
	skuItemInfos := make([]SkuItemInfo, 0, len(order.SkuInfos))
	for _, item := range order.SkuInfos {
		skuItem := SkuItemInfo{
			SkuAmount:      item.SkuAmount,
			Priority:       item.Priority,
			ShippingAmount: item.ShippingAmount,
		}
		if item.SkuFullInfo != nil {
			skuItem.SkuCode = item.SkuFullInfo.SkuCode
			skuItem.SkuName = item.SkuFullInfo.SkuName
		}
		skuItemInfos = append(skuItemInfos, skuItem)
	}
	return OrderInfoQueryResp{
		OrderNo:         order.OrderNo,
		ExternalOrderNo: order.ExternalOrderNo,
		OrderSource:     order.OrderSource.Code(),
		OrderSourceDesc: order.OrderSource.Desc(),
		OrderStatus:     order.OrderStatus.Code(),
		OrderStatusDesc: order.OrderStatus.Desc(),
		SkuItemInfos:    skuItemInfos,
	}
}

// NewOrderInfoQueryResps 批量转换
func NewOrderInfoQueryResps(orders []*pl.OrderInfo) []OrderInfoQueryResp {
	result := make([]OrderInfoQueryResp, 0, len(orders))
	for _, order := range orders {
		result = append(result, NewOrderInfoQueryResp(order))
	}
	return result
}
