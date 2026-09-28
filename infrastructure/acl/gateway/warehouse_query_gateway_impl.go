package gateway

import (
	"github.com/go-spring/spring-core/gs"

	domainpl "think.com/go-ddd/domain/pl"
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
	"think.com/go-ddd/domain/port/gateway"
)

func init() {
	gs.Object(new(WarehouseQueryGatewayImpl)).Export((*gateway.WarehouseQueryGateway)(nil))
}

// WarehouseQueryGatewayImpl 仓库域南向网关实现。
//
// 提供分仓所需的仓库库存、店铺/sku 指定仓映射；内置数据用于本地演示，
// 真实项目替换为调用仓储系统。
type WarehouseQueryGatewayImpl struct {
	// 是否使用内置仓库数据
	MockEnabled bool `value:"${oms.warehouse.mock:=true}"`
}

// Query 查询参与分仓的仓库与指定仓映射
func (g *WarehouseQueryGatewayImpl) Query(req *request.WarehouseQueryRequest) (*response.WarehouseQueryResponse, error) {
	return &response.WarehouseQueryResponse{
		WarehouseInfos: builtinWarehouses,
		// 演示：组合商品 SKU1004 指定从广州仓发货；其余 sku 走最优仓计算
		SkuMappingWarehouseMap: map[string]string{"SKU1004": "WH_GZ_01"},
		// 未配置店铺指定仓，让最优仓策略生效
		StoreMappingWarehouseMap: map[string]string{},
	}, nil
}

// builtinWarehouses 内置仓库数据（本地演示用）
//
// Score = shippingLevel - distance，得分越高越优先被选中。
var builtinWarehouses = []*domainpl.WarehouseInfo{
	{
		WarehouseCode: "WH_SH_01",
		AreaCode:      "310000",
		Distance:      12.5,
		ShippingLevel: 9,
		InventoryMap:  map[string]int64{"SKU1001": 100, "SKU1002": 50, "SKU1003": 0, "SKU1004": 0},
	},
	{
		WarehouseCode: "WH_GZ_01",
		AreaCode:      "440100",
		Distance:      850.0,
		ShippingLevel: 8,
		InventoryMap:  map[string]int64{"SKU1001": 30, "SKU1002": 100, "SKU1003": 0, "SKU1004": 20},
	},
	{
		WarehouseCode: "WH_BJ_01",
		AreaCode:      "110000",
		Distance:      1200.0,
		ShippingLevel: 7,
		InventoryMap:  map[string]int64{"SKU1001": 0, "SKU1002": 10, "SKU1003": 0, "SKU1004": 0},
	},
}
