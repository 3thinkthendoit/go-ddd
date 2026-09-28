package constant

import "fmt"

// OrderSource 订单来源（外部平台）
type OrderSource int

const (
	OrderSourceUnKnow   OrderSource = -100 // 未知
	OrderSourceTaoBao   OrderSource = 1    // 淘宝天猫
	OrderSourcePdd      OrderSource = 2    // 拼多多
	OrderSourceDouYin   OrderSource = 3    // 抖音
	OrderSourceJd       OrderSource = 4    // 京东
	OrderSourceKuaiShou OrderSource = 5    // 快手
	OrderSourceMiJia    OrderSource = 6    // 小米/米家
)

var orderSourceDesc = map[OrderSource]string{
	OrderSourceUnKnow:   "未知",
	OrderSourceTaoBao:   "淘宝天猫",
	OrderSourcePdd:      "拼多多",
	OrderSourceDouYin:   "抖音",
	OrderSourceJd:       "京东",
	OrderSourceKuaiShou: "快手",
	OrderSourceMiJia:    "小米",
}

// Code 返回来源编码
func (s OrderSource) Code() int { return int(s) }

// Desc 返回来源描述
func (s OrderSource) Desc() string { return orderSourceDesc[s] }

func (s OrderSource) String() string {
	if desc, ok := orderSourceDesc[s]; ok {
		return desc
	}
	return fmt.Sprintf("OrderSource(%d)", int(s))
}

// OrderSourceOfByCode 按编码查找订单来源，未知编码返回 OrderSourceUnKnow
func OrderSourceOfByCode(code int) OrderSource {
	s := OrderSource(code)
	if _, ok := orderSourceDesc[s]; ok {
		return s
	}
	return OrderSourceUnKnow
}
