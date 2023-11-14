package request

type UpdatePaymentSystemRequest struct {
	Status int32 `json:"status,omitempty"` // 決済システム状態
}
