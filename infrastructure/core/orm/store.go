// Package orm 持久化组件（MyBatis / GORM 的等价物）。
//
// 这里给的是内存表实现，让示例可以零依赖跑通；接入 MySQL 时保持本文件的
// 方法签名不变、替换为真实 SQL 即可，上层 repository 实现无需改动。
package orm

import (
	"sync"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/infrastructure/core/orm/po"
)

func init() {
	gs.Object(NewStore())
}

// Store 内存表存储
type Store struct {
	mu sync.RWMutex

	seq map[string]int64

	orderBaseInfos          []*po.OrderBaseInfo
	orderSkuInfos           []*po.OrderSkuInfo
	orderSkuItemInfos       []*po.OrderSkuItemInfo
	orderSplitResults       []*po.OrderSplitResult
	shippingCallbackRecords []*po.ShippingCallbackRecord
}

// NewStore 创建存储
func NewStore() *Store {
	return &Store{seq: map[string]int64{}}
}

func (s *Store) nextId(table string) int64 {
	s.seq[table]++
	return s.seq[table]
}

// ---------------------------------------------------------------- 写入

// InsertOrderBaseInfo 插入订单主表
func (s *Store) InsertOrderBaseInfo(v *po.OrderBaseInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v.Id = s.nextId("order_base_info")
	s.orderBaseInfos = append(s.orderBaseInfos, v)
	return nil
}

// InsertOrderSkuInfos 批量插入订单 sku
func (s *Store) InsertOrderSkuInfos(list []*po.OrderSkuInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range list {
		v.Id = s.nextId("order_sku_info")
		s.orderSkuInfos = append(s.orderSkuInfos, v)
	}
	return nil
}

// InsertOrderSkuItemInfos 批量插入订单 sku item
func (s *Store) InsertOrderSkuItemInfos(list []*po.OrderSkuItemInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range list {
		v.Id = s.nextId("order_sku_item_info")
		s.orderSkuItemInfos = append(s.orderSkuItemInfos, v)
	}
	return nil
}

// InsertOrderSplitResults 批量插入拆单结果
func (s *Store) InsertOrderSplitResults(list []*po.OrderSplitResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range list {
		v.Id = s.nextId("order_split_result")
		s.orderSplitResults = append(s.orderSplitResults, v)
	}
	return nil
}

// InsertShippingCallbackRecord 插入发货回传记录
func (s *Store) InsertShippingCallbackRecord(v *po.ShippingCallbackRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v.Id = s.nextId("shipping_callback_record")
	s.shippingCallbackRecords = append(s.shippingCallbackRecords, v)
	return nil
}

// ---------------------------------------------------------------- 查询

// SelectOrderBaseInfos 按条件查询订单主表
func (s *Store) SelectOrderBaseInfos(where func(*po.OrderBaseInfo) bool) []*po.OrderBaseInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*po.OrderBaseInfo, 0)
	for _, v := range s.orderBaseInfos {
		if where == nil || where(v) {
			copied := *v
			result = append(result, &copied)
		}
	}
	return result
}

// SelectOrderSkuInfos 按条件查询订单 sku
func (s *Store) SelectOrderSkuInfos(where func(*po.OrderSkuInfo) bool) []*po.OrderSkuInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*po.OrderSkuInfo, 0)
	for _, v := range s.orderSkuInfos {
		if where == nil || where(v) {
			copied := *v
			result = append(result, &copied)
		}
	}
	return result
}

// SelectOrderSkuItemInfos 按条件查询订单 sku item
func (s *Store) SelectOrderSkuItemInfos(where func(*po.OrderSkuItemInfo) bool) []*po.OrderSkuItemInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*po.OrderSkuItemInfo, 0)
	for _, v := range s.orderSkuItemInfos {
		if where == nil || where(v) {
			copied := *v
			result = append(result, &copied)
		}
	}
	return result
}

// SelectOrderSplitResults 按条件查询拆单结果
func (s *Store) SelectOrderSplitResults(where func(*po.OrderSplitResult) bool) []*po.OrderSplitResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*po.OrderSplitResult, 0)
	for _, v := range s.orderSplitResults {
		if where == nil || where(v) {
			copied := *v
			result = append(result, &copied)
		}
	}
	return result
}

// SelectShippingCallbackRecords 按条件查询发货回传记录
func (s *Store) SelectShippingCallbackRecords(where func(*po.ShippingCallbackRecord) bool) []*po.ShippingCallbackRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*po.ShippingCallbackRecord, 0)
	for _, v := range s.shippingCallbackRecords {
		if where == nil || where(v) {
			copied := *v
			result = append(result, &copied)
		}
	}
	return result
}

// ---------------------------------------------------------------- 更新

// UpdateOrderBaseInfo 按条件更新订单主表，返回受影响行数
func (s *Store) UpdateOrderBaseInfo(where func(*po.OrderBaseInfo) bool, apply func(*po.OrderBaseInfo)) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	affected := 0
	for _, v := range s.orderBaseInfos {
		if where != nil && !where(v) {
			continue
		}
		apply(v)
		affected++
	}
	return affected
}

// UpdateOrderSkuItemInfos 按条件更新订单 sku item，返回受影响行数
func (s *Store) UpdateOrderSkuItemInfos(where func(*po.OrderSkuItemInfo) bool, apply func(*po.OrderSkuItemInfo)) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	affected := 0
	for _, v := range s.orderSkuItemInfos {
		if where != nil && !where(v) {
			continue
		}
		apply(v)
		affected++
	}
	return affected
}

// UpdateOrderSplitResults 按条件更新拆单结果，返回受影响行数
func (s *Store) UpdateOrderSplitResults(where func(*po.OrderSplitResult) bool, apply func(*po.OrderSplitResult)) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	affected := 0
	for _, v := range s.orderSplitResults {
		if where != nil && !where(v) {
			continue
		}
		apply(v)
		affected++
	}
	return affected
}
