package types

type PaymentSystem struct {
	MethodType int32 `json:"methodType"` // 決済システム種別
	Status     int32 `json:"status"`     // 決済システム状態
	CreatedAt  int64 `json:"createdAt"`  // 登録日時
	UpdatedAt  int64 `json:"updatedAt"`  // 更新日時
}

type UpdatePaymentSystemRequest struct {
	Status int32 `json:"status" validate:"required"` // 決済システム状態
}

type PaymentSystemsResponse struct {
	Systems []*PaymentSystem `json:"systems"` // 決済システム一覧
}
