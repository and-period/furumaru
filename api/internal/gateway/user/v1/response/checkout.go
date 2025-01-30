package response

type CheckoutResponse struct {
	URL string `json:"url"` // 遷移先URL
}

type CheckoutStateResponse struct {
	OrderID string `json:"orderId"` // 注文履歴ID
	Status  int32  `json:"status"`  // 注文ステータス
}

type PreCheckoutExperienceResponse struct {
	RequestID  string      `json:"requestId"`  // 支払い時にAPIへ送信するキー(重複判定用)
	Experience *Experience `json:"experience"` // 体験情報
	Promotion  *Promotion  `json:"promotion"`  // プロモーション情報
	SubTotal   int64       `json:"subtotal"`   // 購入金額(税込)
	Discount   int64       `json:"discount"`   // 割引金額(税込)
	Total      int64       `json:"total"`      // 合計金額(税込)
}
