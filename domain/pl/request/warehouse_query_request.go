package request

// WarehouseQueryRequest 仓库信息查询请求
type WarehouseQueryRequest struct {
	// 需要查库存的 skuCode
	SkuCodes []string
	// 店铺编码，用于查询店铺指定仓
	StoreCode string
	// 收货区域编码，用于就近分仓
	AreaCode string
}
