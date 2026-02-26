package types

// PaymentSystemStatus - 決済システム状態
type PaymentSystemStatus int32

const (
	PaymentSystemStatusUnknown PaymentSystemStatus = 0
	PaymentSystemStatusInUse   PaymentSystemStatus = 1 // 利用可能
	PaymentSystemStatusOutage  PaymentSystemStatus = 2 // 停止中
)

// PaymentProviderType - 決済プロバイダー種別
type PaymentProviderType int32

const (
	PaymentProviderTypeUnknown PaymentProviderType = 0
	PaymentProviderTypeKomoju  PaymentProviderType = 1 // KOMOJU
	PaymentProviderTypeStripe  PaymentProviderType = 2 // Stripe
)

type PaymentSystem struct {
	MethodType   PaymentMethodType   `json:"methodType"`   // 決済システム種別
	ProviderType PaymentProviderType `json:"providerType"` // 決済プロバイダー種別
	Status       PaymentSystemStatus `json:"status"`       // 決済システム状態
	CreatedAt    int64               `json:"createdAt"`    // 登録日時
	UpdatedAt    int64               `json:"updatedAt"`    // 更新日時
}

type UpdatePaymentSystemRequest struct {
	Status       PaymentSystemStatus `json:"status" validate:"required"`       // 決済システム状態
	ProviderType PaymentProviderType `json:"providerType" validate:"required"` // 決済プロバイダー種別
}

type PaymentSystemsResponse struct {
	Systems []*PaymentSystem `json:"systems"` // 決済システム一覧
}
