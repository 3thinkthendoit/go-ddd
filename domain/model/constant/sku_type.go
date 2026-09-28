package constant

import "fmt"

// SkuType 商品类型
type SkuType int

const (
	SkuTypeSingle  SkuType = 1 // 单品
	SkuTypeCombine SkuType = 2 // 组合
	SkuTypeVirtual SkuType = 3 // 虚拟商品
)

var skuTypeDesc = map[SkuType]string{
	SkuTypeSingle:  "单品",
	SkuTypeCombine: "组合",
	SkuTypeVirtual: "虚拟商品",
}

// Code 返回商品类型编码
func (t SkuType) Code() int { return int(t) }

// Desc 返回商品类型描述
func (t SkuType) Desc() string { return skuTypeDesc[t] }

func (t SkuType) String() string {
	if desc, ok := skuTypeDesc[t]; ok {
		return desc
	}
	return fmt.Sprintf("SkuType(%d)", int(t))
}

// SkuTypeOf 按编码查找商品类型
func SkuTypeOf(code int) (SkuType, bool) {
	t := SkuType(code)
	_, ok := skuTypeDesc[t]
	return t, ok
}
