package request

type UpdatePaymentSystemRequest struct {
	Status int32 `json:"status"` // 決済システム状態
}
