package pl

// WarehouseInfo 仓库信息（库存 / 处理能力 / 距离）
type WarehouseInfo struct {
	WarehouseCode string
	AreaCode      string
	Distance      float64
	ShippingLevel int
	InventoryMap  map[string]int64
}
