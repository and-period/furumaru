package request

type DraftOrderRequest struct {
	ShippingMessage string `json:"shippingMessage,omitempty"` // 発送連絡時のメッセージ
}

type CompleteOrderRequest struct {
	ShippingMessage string `json:"shippingMessage,omitempty"` // 発送連絡時のメッセージ
}

type RefundOrderRequest struct {
	Description string `json:"description,omitempty"` // 返金理由
}

type UpdateOrderFulfillmentRequest struct {
	ShippingCarrier int32  `json:"shippingCarrier,omitempty"` // 配送会社
	TrackingNumber  string `json:"trackingNumber,omitempty"`  // 伝票番号
}

type ExportOrdersRequest struct {
	ShippingCarrier       int32 `json:"shippingCarrier,omitempty"`       // 配送会社
	CharacterEncodingType int32 `json:"characterEncodingType,omitempty"` // 文字コード種別
}
