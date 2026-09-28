package po

// OrderSkuItemInfo 订单 sku item 表（order_sku_item_info）
//
// item 级别记录下单、发货、退货、签收数量，是拆单与履约的最小单位。
type OrderSkuItemInfo struct {
	Id             int64
	OrderNo        string
	SkuCode        string
	SkuName        string
	SkuAmount      int
	PayPrice       int64
	Priority       int
	FeeAmountInfos string
	ShippingAmount int
	ReturnAmount   int
	ConfirmAmount  int
	SkuStatus      int
}
