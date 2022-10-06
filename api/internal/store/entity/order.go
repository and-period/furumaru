package entity

import "time"

// 支払いステータス
type PaymentStatus int32

const (
	PaymentStatusFailed     PaymentStatus = 0 // 失敗
	PaymentStatusProcessing PaymentStatus = 1 // 申請中
	PaymentStatusConfirmed  PaymentStatus = 2 // 仮売上・オーソリ
	PaymentStatusCaptured   PaymentStatus = 3 // 実売上・キャプチャ
	PaymentStatusCanceled   PaymentStatus = 4 // キャンセル
)

// 配送ステータス
type FulfillmentStatus int32

const (
	FulfillmentStatusUnfulfilled FulfillmentStatus = 0 // 未発送
	FulfillmentStatusFulfilled   FulfillmentStatus = 1 // 発送済み
)

// 注文キャンセル種別
type CancelType int32

const (
	CancelTypeUnknown CancelType = 0
)

// Order - 注文履歴情報
type Order struct {
	ID                string            `gorm:"primaryKey;<-:create"` // 注文履歴ID
	UserID            string            `gorm:""`                     // ユーザーID
	PaymentStatus     PaymentStatus     `gorm:""`                     // 支払いステータス
	FulfillmentStatus FulfillmentStatus `gorm:""`                     // 配送ステータス
	CancelType        CancelType        `gorm:""`                     // 注文キャンセル種別
	CancelReason      string            `gorm:""`                     // 注文キャンセル理由
	CanceledAt        time.Time         `gorm:"default:null"`         // 注文キャンセル日時
	OrderedAt         time.Time         `gorm:"default:null"`         // 決済要求日時
	ConfirmedAt       time.Time         `gorm:"default:null"`         // 決済実行日時
	CapturedAt        time.Time         `gorm:"default:null"`         // 決済確定日時
	DeliveredAt       time.Time         `gorm:"default:null"`         // 配送日時
	CreatedAt         time.Time         `gorm:"<-:create"`            // 登録日時
	UpdatedAt         time.Time         `gorm:""`                     // 更新日時
}
