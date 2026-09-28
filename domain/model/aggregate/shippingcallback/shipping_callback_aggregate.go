package shippingcallback

import (
	"encoding/json"
	"fmt"
	"sort"

	"think.com/go-ddd/domain/common/assert"
	"think.com/go-ddd/domain/model/constant"
	"think.com/go-ddd/domain/model/dp"
	"think.com/go-ddd/domain/pl"
	"think.com/go-ddd/domain/pl/command"
)

// ShippingCallbackAggregate 发货回传聚合根。
//
// 只拿到 WMS 发货单号时先创建聚合，再由领域服务查询订单信息补全 sku 明细，
// 最后判断「是否需要回传上游平台」并落回传记录。
type ShippingCallbackAggregate struct {
	orderId *dp.OrderId

	// 是否需要回传上游平台
	callback bool

	// 订单 sku item：key 为内部 skuCode
	orderSkuItems map[string]*ShippingSkuItem

	// WMS 回传的发货明细，用于补全 orderSkuItems
	shippingInfos []command.ShippingInfo

	// WMS 发货单号（拆单后的父单号）
	wmsOrderNo string

	shippingCallbackRecord *ShippingCallbackRecord
}

// CreateShippingCallbackAggregate 由发货回传命令创建聚合根
func CreateShippingCallbackAggregate(cmd *command.SkuShippingCommand) (*ShippingCallbackAggregate, error) {
	if err := assert.NotNull(cmd, "command is null"); err != nil {
		return nil, err
	}
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	return &ShippingCallbackAggregate{
		orderId:       dp.NewOrderIdByOrderNo(cmd.OrderNo),
		orderSkuItems: map[string]*ShippingSkuItem{},
		shippingInfos: cmd.ShippingInfos,
		wmsOrderNo:    cmd.WmsOrderNo,
	}, nil
}

// InitBaseInfo 用订单信息补全 sku 明细，并累加 WMS 回传的发货数量。
//
// 先按订单上的「已发货数量」建 item，再把本次 WMS 回传的数量累加上去，
// 这样同一订单多次（部分）发货时累计值才正确。
func (a *ShippingCallbackAggregate) InitBaseInfo(skuItemInfos []*pl.SkuItemInfo) error {
	for _, skuItemInfo := range skuItemInfos {
		if skuItemInfo == nil || skuItemInfo.SkuFullInfo == nil {
			continue
		}
		skuCode := skuItemInfo.SkuFullInfo.SkuCode
		a.orderSkuItems[skuCode] = NewShippingSkuItem(
			skuCode, skuItemInfo.SkuFullInfo.ExternalSkuId,
			skuItemInfo.SkuAmount, skuItemInfo.ShippingAmount)
	}
	// 合并 WMS 回传的发货信息（累加本次发货数量）
	for _, shippingInfo := range a.shippingInfos {
		item, ok := a.orderSkuItems[shippingInfo.SkuCode]
		if !ok {
			return fmt.Errorf("skuCode=%s 无法匹配订单 skuCode", shippingInfo.SkuCode)
		}
		item.AddShippingInfo(shippingInfo.SkuAmount, shippingInfo.ExpressCode, shippingInfo.ExpressNo)
	}
	return nil
}

// CompleteOrderId 补全订单标识（订单来源 / 外部订单号），回传上游平台时需要
func (a *ShippingCallbackAggregate) CompleteOrderId(externalOrderNo string, orderSource constant.OrderSource) {
	a.orderId.Init(externalOrderNo, orderSource)
}

// Check 回传校验：判断是否真的需要回传上游平台
func (a *ShippingCallbackAggregate) Check() error {
	a.callback = false
	if len(a.orderSkuItems) == 0 {
		return fmt.Errorf("根据 orderNo=[%s] 查询不到订单信息", a.orderId.OrderNo())
	}
	// 抖音平台无需发货回传
	if a.orderId.OrderSource() == constant.OrderSourceDouYin {
		return nil
	}
	for _, item := range a.orderSkuItems {
		if item.ShippingAmount() > 0 {
			a.callback = true
			break
		}
	}
	return nil
}

// shippingSkuItemSnapshot 回传记录里的 sku 明细快照。
//
// ShippingSkuItem 的字段都是未导出的，直接 json.Marshal 只会得到空对象，
// 所以留痕时必须先转成带 json tag 的快照结构。
type shippingSkuItemSnapshot struct {
	SkuCode        string `json:"skuCode"`
	ExternalSkuId  string `json:"externalSkuId"`
	SkuAmount      int    `json:"skuAmount"`
	ShippingAmount int    `json:"shippingAmount"`
	ExpressCode    string `json:"expressCode"`
	ExpressNo      string `json:"expressNo"`
}

// HandleCallbackResult 处理回传结果并留痕
func (a *ShippingCallbackAggregate) HandleCallbackResult(callbackResult *pl.ShippingCallbackResult) error {
	if err := assert.NotNull(callbackResult, "callbackResult is null"); err != nil {
		return err
	}
	snapshots := make([]shippingSkuItemSnapshot, 0, len(a.orderSkuItems))
	for _, item := range a.orderSkuItems {
		snapshots = append(snapshots, shippingSkuItemSnapshot{
			SkuCode:        item.SkuCode(),
			ExternalSkuId:  item.ExternalSkuId(),
			SkuAmount:      item.SkuAmount(),
			ShippingAmount: item.ShippingAmount(),
			ExpressCode:    item.ExpressCode(),
			ExpressNo:      item.ExpressNo(),
		})
	}
	// 按 skuCode 排序，保证留痕 JSON 稳定可比对
	sort.Slice(snapshots, func(i, j int) bool { return snapshots[i].SkuCode < snapshots[j].SkuCode })

	skuItems, err := json.Marshal(snapshots)
	if err != nil {
		return err
	}

	a.shippingCallbackRecord = NewShippingCallbackRecord(
		a.orderId.OrderNo(),
		a.orderId.ExternalOrderNo(),
		a.orderId.OrderSource(),
		string(skuItems),
		callbackResult.CallStatus,
		callbackResult.Result,
	)
	return nil
}

// OrderId 订单标识
func (a *ShippingCallbackAggregate) OrderId() *dp.OrderId { return a.orderId }

// WmsOrderNo WMS 发货单号（拆单后的父单号）
func (a *ShippingCallbackAggregate) WmsOrderNo() string { return a.wmsOrderNo }

// Callback 是否需要回传上游平台
func (a *ShippingCallbackAggregate) Callback() bool { return a.callback }

// OrderSkuItems 订单 sku item：key 为内部 skuCode
func (a *ShippingCallbackAggregate) OrderSkuItems() map[string]*ShippingSkuItem {
	return a.orderSkuItems
}

// ShippingCallbackRecord 回传记录
func (a *ShippingCallbackAggregate) ShippingCallbackRecord() *ShippingCallbackRecord {
	return a.shippingCallbackRecord
}
