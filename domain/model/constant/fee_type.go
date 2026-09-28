package constant

import "fmt"

// FeeType 费用类型
type FeeType int

const (
	FeeTypeInstallFee  FeeType = 0 // 安装费
	FeeTypeTranFee     FeeType = 1 // 运费
	FeeTypeDiscountFee FeeType = 2 // 优惠金额
	FeeTypeCommission  FeeType = 3 // 提成
)

var feeTypeDesc = map[FeeType]string{
	FeeTypeInstallFee:  "安装费",
	FeeTypeTranFee:     "运费",
	FeeTypeDiscountFee: "优惠金额",
	FeeTypeCommission:  "提成",
}

// Code 返回费用类型编码
func (f FeeType) Code() int { return int(f) }

// Desc 返回费用类型描述
func (f FeeType) Desc() string { return feeTypeDesc[f] }

func (f FeeType) String() string {
	if desc, ok := feeTypeDesc[f]; ok {
		return desc
	}
	return fmt.Sprintf("FeeType(%d)", int(f))
}

// FeeTypeOf 按编码查找费用类型
func FeeTypeOf(code int) (FeeType, bool) {
	f := FeeType(code)
	_, ok := feeTypeDesc[f]
	return f, ok
}
