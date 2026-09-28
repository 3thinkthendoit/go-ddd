package shippingcallback

import (
	"encoding/json"
	"strings"
	"testing"

	"think.com/go-ddd/domain/model/constant"
	"think.com/go-ddd/domain/pl"
	"think.com/go-ddd/domain/pl/command"
)

// newTestAggregate 构造一个只含一条发货信息的回传聚合
func newTestAggregate(t *testing.T, shippingAmount int) *ShippingCallbackAggregate {
	t.Helper()
	aggregate, err := CreateShippingCallbackAggregate(&command.SkuShippingCommand{
		OrderNo:    "O10001",
		WmsOrderNo: "W10001",
		ShippingInfos: []command.ShippingInfo{
			{SkuCode: "SKU1001", SkuAmount: shippingAmount, ExpressCode: "SF", ExpressNo: "SF001"},
		},
	})
	if err != nil {
		t.Fatalf("创建聚合失败: %v", err)
	}
	return aggregate
}

// TestInitBaseInfoAccumulatesShippingAmount 回归测试：
// WMS 回传的是「本次发货数量」（增量），多次回传必须累加，
// 直接赋值会导致部分发货永远凑不满、主订单无法流转为已发货。
func TestInitBaseInfoAccumulatesShippingAmount(t *testing.T) {
	// 订单上 SKU1001 已发 1 件，总量 2 件
	aggregate := newTestAggregate(t, 1)
	if err := aggregate.InitBaseInfo([]*pl.SkuItemInfo{
		{
			SkuFullInfo:    &pl.SkuFullInfo{SkuCode: "SKU1001", ExternalSkuId: "EXT_SKU_1001"},
			SkuAmount:      2,
			Priority:       1,
			ShippingAmount: 1,
		},
	}); err != nil {
		t.Fatalf("InitBaseInfo 失败: %v", err)
	}

	item, ok := aggregate.OrderSkuItems()["SKU1001"]
	if !ok {
		t.Fatal("未找到 SKU1001")
	}
	if got, want := item.ShippingAmount(), 2; got != want {
		t.Fatalf("已发货数量未累加: got=%d want=%d", got, want)
	}
	if got, want := item.ExpressNo(), "SF001"; got != want {
		t.Fatalf("快递单号未更新: got=%s want=%s", got, want)
	}
}

// TestInitBaseInfoKeepsOtherSkuShippingAmount 回归测试：
// 本次回传只涉及部分 sku 时，其他 sku 的已发货数量不能被冲成 0。
func TestInitBaseInfoKeepsOtherSkuShippingAmount(t *testing.T) {
	aggregate := newTestAggregate(t, 1)
	if err := aggregate.InitBaseInfo([]*pl.SkuItemInfo{
		{
			SkuFullInfo:    &pl.SkuFullInfo{SkuCode: "SKU1001", ExternalSkuId: "EXT_SKU_1001"},
			SkuAmount:      2,
			Priority:       1,
			ShippingAmount: 1,
		},
		{
			SkuFullInfo:    &pl.SkuFullInfo{SkuCode: "SKU1002", ExternalSkuId: "EXT_SKU_1002"},
			SkuAmount:      1,
			Priority:       1,
			ShippingAmount: 1,
		},
	}); err != nil {
		t.Fatalf("InitBaseInfo 失败: %v", err)
	}

	if got := aggregate.OrderSkuItems()["SKU1002"].ShippingAmount(); got != 1 {
		t.Fatalf("无关 sku 的已发货数量被覆盖: got=%d want=1", got)
	}
}

// TestCheckOnlyCallbacksWhenShipped 回传判断：有发货数量才回传上游，抖音不回传。
func TestCheckOnlyCallbacksWhenShipped(t *testing.T) {
	aggregate := newTestAggregate(t, 1)
	if err := aggregate.InitBaseInfo([]*pl.SkuItemInfo{
		{SkuFullInfo: &pl.SkuFullInfo{SkuCode: "SKU1001"}, SkuAmount: 2, Priority: 1},
	}); err != nil {
		t.Fatalf("InitBaseInfo 失败: %v", err)
	}
	if err := aggregate.Check(); err != nil {
		t.Fatalf("Check 失败: %v", err)
	}
	if !aggregate.Callback() {
		t.Fatal("有发货数量时应当回传上游")
	}

	aggregate.CompleteOrderId("EXT001", constant.OrderSourceDouYin)
	if err := aggregate.Check(); err != nil {
		t.Fatalf("Check 失败: %v", err)
	}
	if aggregate.Callback() {
		t.Fatal("抖音平台无需发货回传")
	}
}

// TestHandleCallbackResultSerializesSkuItems 留痕记录的 sku 明细必须是真实内容，
// 而不是未导出字段导致的空对象 {"SKU1001":{}}。
func TestHandleCallbackResultSerializesSkuItems(t *testing.T) {
	aggregate := newTestAggregate(t, 1)
	if err := aggregate.InitBaseInfo([]*pl.SkuItemInfo{
		{
			SkuFullInfo:    &pl.SkuFullInfo{SkuCode: "SKU1001", ExternalSkuId: "EXT_SKU_1001"},
			SkuAmount:      2,
			Priority:       1,
			ShippingAmount: 1,
		},
	}); err != nil {
		t.Fatalf("InitBaseInfo 失败: %v", err)
	}
	if err := aggregate.HandleCallbackResult(&pl.ShippingCallbackResult{CallStatus: 0, Result: "ok"}); err != nil {
		t.Fatalf("HandleCallbackResult 失败: %v", err)
	}

	record := aggregate.ShippingCallbackRecord()
	if record == nil {
		t.Fatal("未生成回传记录")
	}
	if strings.Contains(record.SkuItems(), `"SKU1001":{}`) {
		t.Fatalf("sku 明细被序列化成空对象: %s", record.SkuItems())
	}

	var snapshots []shippingSkuItemSnapshot
	if err := json.Unmarshal([]byte(record.SkuItems()), &snapshots); err != nil {
		t.Fatalf("sku 明细不是合法 JSON 数组: %v, raw=%s", err, record.SkuItems())
	}
	if len(snapshots) != 1 {
		t.Fatalf("sku 明细条数不对: got=%d want=1", len(snapshots))
	}
	if snapshots[0].SkuCode != "SKU1001" || snapshots[0].ShippingAmount != 2 {
		t.Fatalf("sku 明细内容不对: %+v", snapshots[0])
	}
}
