package request

type DraftOrderRequest struct {
	ShippingMessage string `json:"shippingMessage" binding:"omitempty,max=2000"` // 発送連絡時のメッセージ
}

type CompleteOrderRequest struct {
	ShippingMessage string `json:"shippingMessage" binding:"omitempty,max=2000"` // 発送連絡時のメッセージ
}

type RefundOrderRequest struct {
	Description string `json:"description" binding:"required,max=2000"` // 返金理由
}

type UpdateOrderFulfillmentRequest struct {
	ShippingCarrier int32  `json:"shippingCarrier" binding:"required"` // 配送会社
	TrackingNumber  string `json:"trackingNumber" binding:"required"`  // 伝票番号
}

type ExportOrdersRequest struct {
	ShippingCarrier       int32 `json:"shippingCarrier" binding:"required"` // 配送会社
	CharacterEncodingType int32 `json:"characterEncodingType" binding:""`   // 文字コード種別
}
