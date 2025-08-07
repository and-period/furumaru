package response

type CheckoutResponse struct {
	URL string `json:"url"` // 遷移先URL
}

type CheckoutStateResponse struct {
	OrderID string `json:"orderId"` // 注文履歴ID
	Status  int32  `json:"status"`  // 注文ステータス
}
