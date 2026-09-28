package po

import "time"

// OrderSplitResult 拆单结果表（order_split_result）
//
// 一条记录 = 一个发货单（parent_order_no）下的某个 sku 在某仓库的发货明细。
type OrderSplitResult struct {
	Id            int64
	ParentOrderNo string
	OrderNo       string
	SkuCode       string
	SkuName       string
	SkuAmount     int
	WarehouseCode string
	// 是否已推送 WMS：0 未推送 1 已推送
	FulfillStatus int
	CreateTime    time.Time
	UpdateTime    time.Time
}
