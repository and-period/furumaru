package request

type DraftOrderRequest struct {
	ShippingMessage string `json:"shippingMessage,omitempty"` // 発送連絡時のメッセージ
}

type CompleteOrderRequest struct {
	ShippingMessage string `json:"shippingMessage,omitempty"` // 発送連絡時のメッセージ
}

type UpdateOrderFulfillmentRequest struct {
	ShippingCarrier int32  `json:"shippingCarrier,omitempty"` // 配送会社
	TrackingNumber  string `json:"trackingNumber,omitempty"`  // 伝票番号
}
