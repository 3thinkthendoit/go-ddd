package request

// OrderFulfillSku 发货单内的 sku 明细
type OrderFulfillSku struct {
	SkuCode   string
	SkuAmount int
}

// OrderFulfillRequest 推送 WMS 履约请求
type OrderFulfillRequest struct {
	// 发货单号（拆单后的父单号）
	OmsOrderNo string
	// 原订单号
	OrderNo string
	// 仓库编码
	WarehouseCode string
	// 发货明细
	Skus []OrderFulfillSku
}
