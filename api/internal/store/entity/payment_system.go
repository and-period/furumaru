package entity

import "time"

// PaymentSystemStatus - 決済システム状態
type PaymentSystemStatus int32

const (
	PaymentSystemStatusUnknown PaymentSystemStatus = 0
	PaymentSystemStatusInUse   PaymentSystemStatus = 1 // 利用可能
	PaymentSystemStatusOutage  PaymentSystemStatus = 2 // 停止中（障害発生・メンテナンス）
)

type PaymentSystem struct {
	MethodType   PaymentMethodType   `gorm:"primaryKey;<-:create"` // 決済種別
	ProviderType PaymentProviderType `gorm:"default:1"`            // 決済プロバイダー種別
	Status       PaymentSystemStatus `gorm:""`                     // 決済システム状態
	CreatedAt    time.Time           `gorm:"<-:create"`            // 登録日時
	UpdatedAt    time.Time           `gorm:""`                     // 更新日時
}

type PaymentSystems []*PaymentSystem
