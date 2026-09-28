// Package mijia 小米 / 米家订单协议实现。
//
// 米家通过 HTTP 主动推送订单，北向网关拿到原始报文后交给本客户端解析成领域命令，
// 这样「外部协议」的解析集中在 acl 层，接口层只做转发。
package mijia

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-spring/spring-core/gs"

	"think.com/go-ddd/domain/model/constant"
	"think.com/go-ddd/domain/model/valueobject"
	"think.com/go-ddd/domain/pl"
	"think.com/go-ddd/domain/pl/command"
	"think.com/go-ddd/domain/pl/query"
	"think.com/go-ddd/domain/pl/request"
	"think.com/go-ddd/domain/pl/response"
)

func init() {
	gs.Object(NewClient())
}

// Client 米家协议客户端
type Client struct{}

// NewClient 创建客户端
func NewClient() *Client {
	return &Client{}
}

// orderCreateReq 米家下单报文
type orderCreateReq struct {
	OrderNo        string                 `json:"orderNo"`
	OrderStatus    *int                   `json:"orderStatus"`
	OrderType      int                    `json:"orderType"`
	StoreCode      string                 `json:"storeCode"`
	OrderTitle     string                 `json:"orderTitle"`
	OrderPrice     int64                  `json:"orderPrice"`
	PayType        int                    `json:"payType"`
	DiscountPrice  int64                  `json:"discountPrice"`
	Currency       int                    `json:"currency"`
	UserId         int64                  `json:"userId"`
	UserName       string                 `json:"userName"`
	UserType       int                    `json:"userType"`
	Recipient      string                 `json:"recipient"`
	Mobile         string                 `json:"mobile"`
	Address        string                 `json:"address"`
	InvoiceName    string                 `json:"invoiceName"`
	InvoiceDetails string                 `json:"invoiceDetails"`
	FeeAmountMap   map[string]int64       `json:"feeAmountMap"`
	AttachInfos    map[string]interface{} `json:"attachInfos"`
	SkuInfos       []orderSkuInfo         `json:"skuInfos"`
}

type orderSkuInfo struct {
	ExternalSkuId   string `json:"externalSkuId"`
	ExternalSkuCode string `json:"externalSkuCode"`
	SkuName         string `json:"skuName"`
	SkuPayPrice     int64  `json:"skuPayPrice"`
	SkuBuyAmount    int    `json:"skuBuyAmount"`
}

// ParseCreateOrder 解析米家下单报文为创建订单命令
func (c *Client) ParseCreateOrder(body []byte) (*command.OrderCreateCommand, error) {
	var req orderCreateReq
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("米家下单报文解析失败: %w", err)
	}
	if len(req.SkuInfos) == 0 {
		return nil, fmt.Errorf("米家下单报文缺少 skuInfos")
	}

	feeAmountMap := make(map[constant.FeeType]int64, len(req.FeeAmountMap))
	for code, amount := range req.FeeAmountMap {
		feeType, err := strconv.Atoi(code)
		if err != nil {
			return nil, fmt.Errorf("feeAmountMap 的 key 必须是费用类型编码, key=%s", code)
		}
		feeAmountMap[constant.FeeType(feeType)] = amount
	}

	orderSkuInfos := make([]*pl.OrderSkuInfo, 0, len(req.SkuInfos))
	for _, sku := range req.SkuInfos {
		orderSkuInfos = append(orderSkuInfos, &pl.OrderSkuInfo{
			ExternalSkuId:   sku.ExternalSkuId,
			ExternalSkuCode: sku.ExternalSkuCode,
			SkuName:         sku.SkuName,
			SkuPayPrice:     sku.SkuPayPrice,
			SkuBuyAmount:    sku.SkuBuyAmount,
		})
	}

	// orderStatus 用指针接收：米家只推送已支付订单，报文中缺省该字段时按「已支付」处理；
	// 但必须区分「缺省」与「显式传 0」，否则 int 零值会把显式的未支付状态
	// 悄悄改写成已支付，让领域层的「未支付不允许接入」守卫失效。
	orderStatus := constant.OrderStatusPayed
	if req.OrderStatus != nil {
		orderStatus = constant.OrderStatus(*req.OrderStatus)
	}
	orderType := constant.OrderType(req.OrderType)
	if req.OrderType == 0 {
		orderType = constant.OrderTypeStandard
	}
	currency := constant.Currency(req.Currency)
	if req.Currency == 0 {
		currency = constant.CurrencyCny
	}

	return &command.OrderCreateCommand{
		ExternalOrderNo: req.OrderNo,
		OrderStatus:     orderStatus,
		OrderType:       orderType,
		StoreCode:       req.StoreCode,
		OrderSource:     constant.OrderSourceMiJia,
		OrderTitle:      req.OrderTitle,
		OrderPrice:      req.OrderPrice,
		PayType:         constant.PayType(req.PayType),
		DiscountPrice:   req.DiscountPrice,
		Currency:        currency,
		OrderSkuInfos:   orderSkuInfos,
		UserId:          req.UserId,
		UserName:        req.UserName,
		UserType:        valueobject.UserType(req.UserType),
		Address:         req.Address,
		Mobile:          req.Mobile,
		Recipient:       req.Recipient,
		InvoiceName:     req.InvoiceName,
		InvoiceDetails:  req.InvoiceDetails,
		FeeAmountMap:    feeAmountMap,
		AttachInfos:     req.AttachInfos,
	}, nil
}

// orderQueryReq 米家订单查询报文
type orderQueryReq struct {
	OrderNo         string `json:"orderNo"`
	ExternalOrderNo string `json:"externalOrderNo"`
}

// ParseQuery 解析米家订单查询报文
func (c *Client) ParseQuery(body []byte) (*query.OrderInfoQuery, error) {
	var req orderQueryReq
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("米家查询报文解析失败: %w", err)
	}
	return &query.OrderInfoQuery{
		OrderNo:         req.OrderNo,
		ExternalOrderNo: req.ExternalOrderNo,
		OrderSource:     constant.OrderSourceMiJia.Code(),
	}, nil
}

// ShippingCallback 发货回传
//
// 米家通过订单状态查询获取发货结果，这里按协议回传并返回结果。
func (c *Client) ShippingCallback(req *request.ShippingCallbackRequest) (*response.ShippingCallbackResponse, error) {
	return &response.ShippingCallbackResponse{
		CallbackResult: &pl.ShippingCallbackResult{
			CallStatus: 0,
			Result:     fmt.Sprintf("mijia callback success, externalOrderNo=%s", req.ExternalOrderNo),
		},
	}, nil
}
