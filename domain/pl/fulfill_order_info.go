package pl

// FulfillSku 发货单内的 sku 明细
type FulfillSku struct {
	SkuCode   string
	SkuAmount int
}

// FulfillOrderInfo 发货单信息（拆单后产生，用于推送 WMS）
type FulfillOrderInfo struct {
	OmsOrderNo    string
	OrderNo       string
	WarehouseCode string
	Skus          []FulfillSku
}
