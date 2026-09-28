// Package pl 领域层的转换层（PL：Payload Layer）。
//
// 存放跨域 / 跨进程传输的数据结构：command、event、query、request、response
// 以及南向网关之间共享的 DTO。它们只承载数据、不含业务行为，但仍属于领域层，
// 这样领域服务与外域协作时不必依赖基础设施层或接口层的对象。
//
// 约定：pl 下的对象是纯数据载体，统一使用导出字段（Go 惯例）；
// 只有 domain/model 下的聚合、实体、值对象、DP 才做封装与不变式保护。
package pl

import (
	"think.com/go-ddd/domain/common/assert"
	"think.com/go-ddd/domain/model/constant"
)

// SkuFullInfo 商品域返回的完整商品信息
type SkuFullInfo struct {
	ExternalSkuId   string
	SkuCode         string
	ExternalSkuCode string
	SkuName         string
	SkuPrice        int64
	SkuType         constant.SkuType
	SkuCategory     constant.SkuCategory
	SkuAmount       int
	SubSkuInfos     []*SkuFullInfo
}

// Validate 校验商品信息完整性
func (s *SkuFullInfo) Validate() error {
	if err := assert.NotBlank(s.ExternalSkuId, "externalSkuId is blank"); err != nil {
		return err
	}
	if err := assert.NotBlank(s.SkuCode, "skuCode is blank"); err != nil {
		return err
	}
	if err := assert.NotBlank(s.ExternalSkuCode, "externalSkuCode is blank"); err != nil {
		return err
	}
	if err := assert.NotBlank(s.SkuName, "skuName is blank"); err != nil {
		return err
	}
	if err := assert.True(s.SkuPrice >= 0, "skuPrice is invalid"); err != nil {
		return err
	}
	return assert.True(s.SkuAmount > 0, "skuAmount is invalid")
}
