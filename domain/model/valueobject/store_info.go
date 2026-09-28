package valueobject

import (
	"fmt"

	"think.com/go-ddd/domain/common/assert"
)

// StoreInfo 店铺信息（值对象）。
type StoreInfo struct {
	storeId   string
	storeCode string
	storeName string
}

// NewStoreInfo 以下单时已知的店铺编码创建，其余属性待领域服务补全。
func NewStoreInfo(storeCode string) StoreInfo {
	assert.NotBlank(storeCode, "storeCode is blank")
	return StoreInfo{storeCode: storeCode}
}

// Init 补全店铺信息（调用其他域后回填）。
func (s *StoreInfo) Init(storeId, storeCode, storeName string) error {
	if err := assert.NotBlank(storeCode, "storeCode is blank"); err != nil {
		return err
	}
	s.storeId = storeId
	s.storeCode = storeCode
	s.storeName = storeName
	return nil
}

// StoreId 店铺 ID
func (s StoreInfo) StoreId() string { return s.storeId }

// StoreCode 店铺编码
func (s StoreInfo) StoreCode() string { return s.storeCode }

// StoreName 店铺名称
func (s StoreInfo) StoreName() string { return s.storeName }

func (s StoreInfo) String() string {
	return fmt.Sprintf("StoreInfo{storeId=%s, storeCode=%s, storeName=%s}", s.storeId, s.storeCode, s.storeName)
}
