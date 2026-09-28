package gateway

import (
	"errors"
	"log"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/domain/model/constant"
	domainpl "think.com/go-ddd/domain/pl"
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
	"think.com/go-ddd/domain/port/gateway"
	infrahttp "think.com/go-ddd/infrastructure/core/http"
)

func init() {
	gs.Object(new(SkuInfoQueryGatewayImpl)).Export((*gateway.SkuInfoQueryGateway)(nil))
}

// SkuInfoQueryGatewayImpl 商品信息南向网关实现。
//
// 正常链路：调用商品中心 HTTP 接口；当商品中心不可用（或本地演示模式打开）时，
// 退化为内置商品目录，保证下单链路可以继续跑通。
type SkuInfoQueryGatewayImpl struct {
	httpClient *infrahttp.Client `autowire:""`

	// 商品中心地址
	Endpoint string `value:"${oms.sku-center.endpoint:=http://sku-center/api/sku/queryBySkuIds}"`
	// 是否允许退化为内置商品目录
	FallbackEnabled bool `value:"${oms.sku-center.fallback:=true}"`
}

// Query 批量查询商品完整信息
func (g *SkuInfoQueryGatewayImpl) Query(req *request.SkuInfoQueryRequest) (*response.SkuInfoQueryResponse, error) {
	if len(req.ExternalSkuIds) == 0 {
		return &response.SkuInfoQueryResponse{}, nil
	}

	params := map[string]interface{}{"externalSkuIds": req.ExternalSkuIds}
	data, err := g.httpClient.Post(g.Endpoint, params)
	if err != nil {
		log.Printf("[sku-center] 调用商品中心失败: %v", err)
	} else if skuInfos, parseErr := parseSkuCenterResponse(data); parseErr == nil {
		return &response.SkuInfoQueryResponse{SkuInfos: skuInfos}, nil
	} else {
		log.Printf("[sku-center] 解析商品中心返回失败: %v", parseErr)
	}

	if !g.FallbackEnabled {
		log.Printf("[sku-center] 商品中心不可用且未开启兜底, externalSkuIds=%v", req.ExternalSkuIds)
		return &response.SkuInfoQueryResponse{}, nil
	}
	log.Printf("[sku-center] 使用内置商品目录兜底, externalSkuIds=%v", req.ExternalSkuIds)

	skuInfos := make([]*domainpl.SkuFullInfo, 0, len(req.ExternalSkuIds))
	for _, externalSkuId := range req.ExternalSkuIds {
		if sku, ok := builtinSkuCatalog[externalSkuId]; ok {
			skuInfos = append(skuInfos, sku)
		}
	}
	return &response.SkuInfoQueryResponse{SkuInfos: skuInfos}, nil
}

// parseSkuCenterResponse 解析商品中心返回报文。
//
// TODO 接入真实商品中心时补充：报文结构确认后按字段映射成 []*pl.SkuFullInfo。
func parseSkuCenterResponse(_ string) ([]*domainpl.SkuFullInfo, error) {
	return nil, errors.New("商品中心返回报文解析未实现")
}

// builtinSkuCatalog 内置商品目录（本地演示用，接入商品中心后可删除）
var builtinSkuCatalog = map[string]*domainpl.SkuFullInfo{
	"EXT_SKU_1001": {
		ExternalSkuId:   "EXT_SKU_1001",
		ExternalSkuCode: "EXT-1001",
		SkuCode:         "SKU1001",
		SkuName:         "阳澄湖大闸蟹 4 对装",
		SkuPrice:        2999,
		SkuType:         constant.SkuTypeSingle,
		SkuCategory:     constant.SkuCategoryFreshFood,
		SkuAmount:       2,
	},
	"EXT_SKU_1002": {
		ExternalSkuId:   "EXT_SKU_1002",
		ExternalSkuCode: "EXT-1002",
		SkuCode:         "SKU1002",
		SkuName:         "316 不锈钢保温杯",
		SkuPrice:        1599,
		SkuType:         constant.SkuTypeSingle,
		SkuCategory:     constant.SkuCategoryGoods,
		SkuAmount:       1,
	},
	"EXT_SKU_1003": {
		ExternalSkuId:   "EXT_SKU_1003",
		ExternalSkuCode: "EXT-1003",
		SkuCode:         "SKU1003",
		SkuName:         "视频会员年卡（虚拟）",
		SkuPrice:        990,
		SkuType:         constant.SkuTypeVirtual,
		SkuCategory:     constant.SkuCategoryGoods,
		SkuAmount:       1,
	},
	"EXT_SKU_1004": {
		ExternalSkuId:   "EXT_SKU_1004",
		ExternalSkuCode: "EXT-1004",
		SkuCode:         "SKU1004",
		SkuName:         "厨房三件套（组合）",
		SkuPrice:        5999,
		SkuType:         constant.SkuTypeCombine,
		SkuCategory:     constant.SkuCategoryGoods,
		SkuAmount:       1,
	},
}
