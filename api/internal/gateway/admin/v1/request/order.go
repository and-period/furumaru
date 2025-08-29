package request

type DraftOrderRequest struct {
	ShippingMessage string `json:"shippingMessage"` // 発送連絡時のメッセージ
}

type CompleteOrderRequest struct {
	ShippingMessage string `json:"shippingMessage"` // 発送連絡時のメッセージ
}

type RefundOrderRequest struct {
	Description string `json:"description"` // 返金理由
}

type UpdateOrderFulfillmentRequest struct {
	ShippingCarrier int32  `json:"shippingCarrier"` // 配送会社
	TrackingNumber  string `json:"trackingNumber"`  // 伝票番号
}

type ExportOrdersRequest struct {
	ShippingCarrier       int32 `json:"shippingCarrier"`       // 配送会社
	CharacterEncodingType int32 `json:"characterEncodingType"` // 文字コード種別
}
