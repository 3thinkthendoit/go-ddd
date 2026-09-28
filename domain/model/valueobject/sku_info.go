package valueobject

import (
	"fmt"

	"think.com/go-ddd/domain/common/assert"
	"think.com/go-ddd/domain/model/constant"
	"think.com/go-ddd/domain/pl"
)

// SkuInfo 商品基本信息（值对象）。
//
// 下单时只知道外部 sku 标识，真实的 skuCode / 价格 / 类型需要调用商品域补全，
// 补全后本对象整体替换，不做局部修改。
type SkuInfo struct {
	skuId           string
	externalSkuId   string
	skuCode         string
	externalSkuCode string
	skuName         string
	skuPrice        int64
	skuType         constant.SkuType
	skuCategory     constant.SkuCategory
}

// NewSkuInfo 直接由内部标识构造
func NewSkuInfo(skuId, skuCode string) SkuInfo {
	return SkuInfo{skuId: skuId, skuCode: skuCode}
}

// CreateSkuInfo 由外部 sku 标识构造（下单入口使用）
func CreateSkuInfo(externalSkuId, externalSkuCode string) (SkuInfo, error) {
	if err := assert.NotBlank(externalSkuId, "externalSkuId is blank"); err != nil {
		return SkuInfo{}, err
	}
	if err := assert.NotBlank(externalSkuCode, "externalSkuCode is blank"); err != nil {
		return SkuInfo{}, err
	}
	return SkuInfo{externalSkuId: externalSkuId, externalSkuCode: externalSkuCode}, nil
}

// Init 用商品域返回的完整信息补全自身
func (s *SkuInfo) Init(full *pl.SkuFullInfo) error {
	if err := full.Validate(); err != nil {
		return err
	}
	s.skuId = full.ExternalSkuId
	s.skuType = full.SkuType
	s.externalSkuId = full.ExternalSkuId
	s.externalSkuCode = full.ExternalSkuCode
	s.skuCode = full.SkuCode
	s.skuName = full.SkuName
	s.skuPrice = full.SkuPrice
	s.skuCategory = full.SkuCategory
	return nil
}

// SkuId 内部 sku 标识
func (s SkuInfo) SkuId() string { return s.skuId }

// ExternalSkuId 外部 sku 标识
func (s SkuInfo) ExternalSkuId() string { return s.externalSkuId }

// SkuCode 内部 skuCode
func (s SkuInfo) SkuCode() string { return s.skuCode }

// ExternalSkuCode 外部 skuCode
func (s SkuInfo) ExternalSkuCode() string { return s.externalSkuCode }

// SkuName 商品名称
func (s SkuInfo) SkuName() string { return s.skuName }

// SkuPrice 商品价格（单位：分）
func (s SkuInfo) SkuPrice() int64 { return s.skuPrice }

// SkuType 商品类型
func (s SkuInfo) SkuType() constant.SkuType { return s.skuType }

// SkuCategory 商品种类
func (s SkuInfo) SkuCategory() constant.SkuCategory { return s.skuCategory }

// IsInit 判断商品信息是否已经被补全
func (s SkuInfo) IsInit() bool { return s.skuCode != "" }

func (s SkuInfo) String() string {
	return fmt.Sprintf("SkuInfo{skuCode=%s, externalSkuCode=%s, skuName=%s}", s.skuCode, s.externalSkuCode, s.skuName)
}
