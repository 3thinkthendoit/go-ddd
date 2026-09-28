package response

import "think.com/go-ddd/domain/pl"

// WarehouseQueryResponse 仓库信息查询结果
type WarehouseQueryResponse struct {
	// 参与分仓计算的仓库（含库存）
	WarehouseInfos []*pl.WarehouseInfo
	// sku 指定发货仓：key 为 skuCode
	SkuMappingWarehouseMap map[string]string
	// 店铺指定发货仓：key 为 storeCode
	StoreMappingWarehouseMap map[string]string
}
