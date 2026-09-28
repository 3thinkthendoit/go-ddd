// Package pl 南向网关的 PL 转换。
//
// 领域对象 ↔ PO ↔ 对外 DTO 的转换全部收敛在这里，避免转换逻辑散落到各个
// gateway / repository 实现里。
package pl

import (
	"encoding/json"
	"time"

	"think.com/go-ddd/domain/model/aggregate/create"
	"think.com/go-ddd/domain/model/aggregate/orderfulfill"
	"think.com/go-ddd/domain/model/aggregate/shippingcallback"
	"think.com/go-ddd/domain/model/constant"
	"think.com/go-ddd/domain/pl"
	"think.com/go-ddd/infrastructure/core/orm/po"
)

// ToOrderBaseInfoPO 订单创建聚合 → 订单主表 PO
func ToOrderBaseInfoPO(aggregate *create.OrderCreateAggregate) *po.OrderBaseInfo {
	now := time.Now()
	return &po.OrderBaseInfo{
		OrderNo:         aggregate.OrderId().OrderNo(),
		ExternalOrderNo: aggregate.OrderId().ExternalOrderNo(),
		StoreCode:       aggregate.StoreInfo().StoreCode(),
		StoreName:       aggregate.StoreInfo().StoreName(),
		OrderTitle:      aggregate.OrderTitle(),
		OrderPrice:      aggregate.OrderPay().PayAmount(),
		OrderSource:     aggregate.OrderId().OrderSource().Code(),
		OrderStatus:     aggregate.OrderStatus().Code(),
		OrderType:       aggregate.OrderType().Code(),
		Remark:          aggregate.Desc(),
		CreateTime:      now,
		UpdateTime:      now,
	}
}

// ToOrderSkuPOs 订单创建聚合 → 订单 sku 表 PO
func ToOrderSkuPOs(aggregate *create.OrderCreateAggregate) []*po.OrderSkuInfo {
	list := make([]*po.OrderSkuInfo, 0, len(aggregate.SkuInfos()))
	for _, orderSku := range aggregate.SkuInfos() {
		skuInfo := orderSku.SkuInfo()
		list = append(list, &po.OrderSkuInfo{
			OrderNo:         aggregate.OrderId().OrderNo(),
			ExternalSkuId:   skuInfo.ExternalSkuId(),
			ExternalSkuCode: skuInfo.ExternalSkuCode(),
			SkuCode:         skuInfo.SkuCode(),
			SkuName:         skuInfo.SkuName(),
			SkuPayPrice:     orderSku.SkuPayPrice(),
			SkuBuyAmount:    orderSku.SkuBuyAmount(),
			SkuType:         orderSku.SkuType().Code(),
			SkuCategory:     orderSku.SkuCategory().Code(),
		})
	}
	return list
}

// ToOrderSkuItemPOs 订单创建聚合 → 订单 sku item 表 PO
func ToOrderSkuItemPOs(aggregate *create.OrderCreateAggregate) []*po.OrderSkuItemInfo {
	skuNameMap := make(map[string]string, len(aggregate.SkuInfos()))
	for _, orderSku := range aggregate.SkuInfos() {
		skuNameMap[orderSku.SkuInfo().SkuCode()] = orderSku.SkuInfo().SkuName()
	}

	list := make([]*po.OrderSkuItemInfo, 0, len(aggregate.SkuItems()))
	for _, item := range aggregate.SkuItems() {
		feeAmountInfos := "{}"
		if data, err := json.Marshal(item.FeeAmountInfos()); err == nil {
			feeAmountInfos = string(data)
		}
		list = append(list, &po.OrderSkuItemInfo{
			OrderNo:        aggregate.OrderId().OrderNo(),
			SkuCode:        item.SkuCode(),
			SkuName:        skuNameMap[item.SkuCode()],
			SkuAmount:      item.SkuAmount(),
			PayPrice:       item.PayPrice(),
			Priority:       item.Priority(),
			FeeAmountInfos: feeAmountInfos,
			SkuStatus:      constant.OrderSkuStatusNotShipped.Code(),
		})
	}
	return list
}

// ToSkuItemInfos PO 的 sku item 列表 → 领域 PL 的 sku item 列表
func ToSkuItemInfos(list []*po.OrderSkuItemInfo) []*pl.SkuItemInfo {
	result := make([]*pl.SkuItemInfo, 0, len(list))
	for _, item := range list {
		result = append(result, &pl.SkuItemInfo{
			SkuFullInfo: &pl.SkuFullInfo{
				SkuCode:   item.SkuCode,
				SkuName:   item.SkuName,
				SkuPrice:  item.PayPrice,
				SkuAmount: item.SkuAmount,
			},
			SkuPayPrice:    item.PayPrice,
			SkuAmount:      item.SkuAmount,
			Priority:       item.Priority,
			ShippingAmount: item.ShippingAmount,
			ReturnAmount:   item.ReturnAmount,
			ConfirmAmount:  item.ConfirmAmount,
		})
	}
	return result
}

// ToOrderInfo 订单主表 PO + sku item PO → 订单查询结果
func ToOrderInfo(order *po.OrderBaseInfo, skuItemInfos []*po.OrderSkuItemInfo) *pl.OrderInfo {
	return &pl.OrderInfo{
		OrderNo:         order.OrderNo,
		ExternalOrderNo: order.ExternalOrderNo,
		OrderSource:     constant.OrderSourceOfByCode(order.OrderSource),
		OrderStatus:     constant.OrderStatus(order.OrderStatus),
		SkuInfos:        ToSkuItemInfos(skuItemInfos),
	}
}

// ToOrderFulfillAggregate 订单主表 PO + sku item PO → 履约聚合根
func ToOrderFulfillAggregate(order *po.OrderBaseInfo, skuItemInfos []*po.OrderSkuItemInfo) (*orderfulfill.OrderFulfillAggregate, error) {
	return orderfulfill.NewOrderFulfillAggregate(order.OrderNo, order.StoreCode, ToSkuItemInfos(skuItemInfos))
}

// ToOrderSplitResultPOs 履约聚合的拆单结果 → 拆单结果表 PO
func ToOrderSplitResultPOs(aggregate *orderfulfill.OrderFulfillAggregate) []*po.OrderSplitResult {
	now := time.Now()
	list := make([]*po.OrderSplitResult, 0)
	for _, splitResults := range aggregate.SplitOrders() {
		for _, result := range splitResults {
			list = append(list, &po.OrderSplitResult{
				ParentOrderNo: result.ParentOrderNo(),
				OrderNo:       result.OrderNo(),
				SkuCode:       result.SkuCode(),
				SkuAmount:     result.SkuAmount(),
				WarehouseCode: result.WarehouseCode(),
				FulfillStatus: 0,
				CreateTime:    now,
				UpdateTime:    now,
			})
		}
	}
	return list
}

// ToFulfillOrderInfos 拆单结果 PO 列表 → 发货单信息（按父单号 + 仓库聚合）
func ToFulfillOrderInfos(list []*po.OrderSplitResult) []*pl.FulfillOrderInfo {
	index := make(map[string]*pl.FulfillOrderInfo)
	order := make([]*pl.FulfillOrderInfo, 0)
	for _, result := range list {
		item, ok := index[result.ParentOrderNo]
		if !ok {
			item = &pl.FulfillOrderInfo{
				OmsOrderNo:    result.ParentOrderNo,
				OrderNo:       result.OrderNo,
				WarehouseCode: result.WarehouseCode,
				Skus:          make([]pl.FulfillSku, 0, 1),
			}
			index[result.ParentOrderNo] = item
			order = append(order, item)
		}
		item.Skus = append(item.Skus, pl.FulfillSku{SkuCode: result.SkuCode, SkuAmount: result.SkuAmount})
	}
	return order
}

// ToShippingCallbackRecordPO 发货回传聚合 → 发货回传记录表 PO
func ToShippingCallbackRecordPO(aggregate *shippingcallback.ShippingCallbackAggregate) *po.ShippingCallbackRecord {
	record := aggregate.ShippingCallbackRecord()
	if record == nil {
		return nil
	}
	return &po.ShippingCallbackRecord{
		OrderNo:         record.OrderNo(),
		ExternalOrderNo: record.ExternalOrderNo(),
		OrderSource:     record.OrderSource().Code(),
		SkuItems:        record.SkuItems(),
		CallStatus:      record.CallStatus(),
		Result:          record.Result(),
		CreateTime:      time.Now(),
	}
}
