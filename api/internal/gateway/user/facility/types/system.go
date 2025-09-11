package types

type PaymentSystem struct {
	MethodType PaymentMethodType `json:"methodType"` // 決済システム種別
	Status     PaymentStatus     `json:"status"`     // 決済システム状態
}
