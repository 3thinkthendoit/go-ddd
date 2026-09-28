package constant

import "fmt"

// SkuCategory 商品种类
type SkuCategory int

const (
	SkuCategoryFreshFood SkuCategory = 1 // 生鲜
	SkuCategoryGoods     SkuCategory = 2 // 百货
)

var skuCategoryDesc = map[SkuCategory]string{
	SkuCategoryFreshFood: "生鲜",
	SkuCategoryGoods:     "百货",
}

// Code 返回商品种类编码
func (c SkuCategory) Code() int { return int(c) }

// Desc 返回商品种类描述
func (c SkuCategory) Desc() string { return skuCategoryDesc[c] }

func (c SkuCategory) String() string {
	if desc, ok := skuCategoryDesc[c]; ok {
		return desc
	}
	return fmt.Sprintf("SkuCategory(%d)", int(c))
}

// SkuCategoryOf 按编码查找商品种类
func SkuCategoryOf(code int) (SkuCategory, bool) {
	c := SkuCategory(code)
	_, ok := skuCategoryDesc[c]
	return c, ok
}
