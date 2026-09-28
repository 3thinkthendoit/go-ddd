package constant

import "fmt"

// Currency 币种
type Currency int

const (
	CurrencyCny Currency = 1 // 人民币
	CurrencyUsd Currency = 2 // 美元
)

var currencyDesc = map[Currency]string{
	CurrencyCny: "CNY",
	CurrencyUsd: "USD",
}

// Code 返回币种编码
func (c Currency) Code() int { return int(c) }

// Desc 返回币种符号
func (c Currency) Desc() string { return currencyDesc[c] }

func (c Currency) String() string {
	if desc, ok := currencyDesc[c]; ok {
		return desc
	}
	return fmt.Sprintf("Currency(%d)", int(c))
}

// CurrencyOf 按编码查找币种
func CurrencyOf(code int) (Currency, bool) {
	c := Currency(code)
	_, ok := currencyDesc[c]
	return c, ok
}
