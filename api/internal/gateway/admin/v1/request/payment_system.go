package request

type UpdatePaymentSystemRequest struct {
	Status int32 `json:"status" validate:"required"` // 決済システム状態
}
