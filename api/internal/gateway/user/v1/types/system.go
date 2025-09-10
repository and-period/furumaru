package types

// PaymentSystemStatus - 決済システム状態
type PaymentSystemStatus int32

const (
	PaymentSystemStatusUnknown PaymentSystemStatus = 0
	PaymentSystemStatusInUse   PaymentSystemStatus = 1 // 利用可能
	PaymentSystemStatusOutage  PaymentSystemStatus = 2 // 停止中
)

type PaymentSystem struct {
	MethodType PaymentMethodType   `json:"methodType"` // 決済システム種別
	Status     PaymentSystemStatus `json:"status"`     // 決済システム状態
}

type PaymentSystemsResponse struct {
	Systems []*PaymentSystem `json:"systems"` // 決済システム一覧
}
