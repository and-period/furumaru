package response

type GuestCheckoutResponse struct {
	URL string `json:"url"` // 遷移先URL
}

type GuestCheckoutStateResponse struct {
	OrderID string `json:"orderId"` // 注文履歴ID
	Status  int32  `json:"status"`  // 注文ステータス
}
