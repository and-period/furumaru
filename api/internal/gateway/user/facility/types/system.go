package types

type PaymentSystem struct {
	MethodType int32 `json:"methodType"` // 決済システム種別
	Status     int32 `json:"status"`     // 決済システム状態
}
