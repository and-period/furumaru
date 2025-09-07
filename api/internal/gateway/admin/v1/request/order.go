package request

type DraftOrderRequest struct {
	ShippingMessage string `json:"shippingMessage" validate:"omitempty,max=2000"` // 発送連絡時のメッセージ
}

type CompleteOrderRequest struct {
	ShippingMessage string `json:"shippingMessage" validate:"omitempty,max=2000"` // 発送連絡時のメッセージ
}

type RefundOrderRequest struct {
	Description string `json:"description" validate:"required,max=2000"` // 返金理由
}

type UpdateOrderFulfillmentRequest struct {
	ShippingCarrier int32  `json:"shippingCarrier" validate:"required"` // 配送会社
	TrackingNumber  string `json:"trackingNumber" validate:"required"`  // 伝票番号
}

type ExportOrdersRequest struct {
	ShippingCarrier       int32 `json:"shippingCarrier" validate:"required"` // 配送会社
	CharacterEncodingType int32 `json:"characterEncodingType" validate:""`   // 文字コード種別
}
