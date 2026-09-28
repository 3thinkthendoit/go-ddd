package response

// OrderFulfillResponse 推送 WMS 履约结果
type OrderFulfillResponse struct {
	IsFulfill bool
	Msg       string
}
