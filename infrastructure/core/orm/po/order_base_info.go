// Package po 持久化对象（Persistent Object）。
//
// 与数据库表一一对应，只承载存储结构，不参与领域逻辑；
// 领域对象 ↔ PO 的转换统一放在 infrastructure/acl/pl 里完成。
package po

import "time"

// OrderBaseInfo 订单主表（order_base_info）
type OrderBaseInfo struct {
	Id              int64
	OrderNo         string
	ExternalOrderNo string
	StoreCode       string
	StoreName       string
	OrderTitle      string
	OrderPrice      int64
	OrderSource     int
	OrderStatus     int
	OrderType       int
	Remark          string
	CreateTime      time.Time
	UpdateTime      time.Time
}
