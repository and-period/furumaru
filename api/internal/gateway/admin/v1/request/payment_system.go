package request

type UpdatePaymentSystemRequest struct {
	Status int32 `json:"status" binding:"required"` // 決済システム状態
}
