package entity

import (
	"time"

	"gorm.io/datatypes"
)

// OrderEventType - 注文イベントログ種別
type OrderEventType int32

const (
	OrderEventTypeUnknown OrderEventType = 0
)

// OrderActivity - 注文イベントログ
type OrderActivity struct {
	ID           string            `gorm:"primaryKey;<-:create"`         // 注文商品ID
	OrderID      string            `gorm:""`                             // 注文履歴ID
	UserID       string            `gorm:""`                             // ユーザーID
	EventType    OrderEventType    `gorm:""`                             // イベントログ種別
	Detail       string            `gorm:""`                             // イベントログ詳細
	Metadata     map[string]string `gorm:"-"`                            // メタデータ
	MetadataJSON datatypes.JSON    `gorm:"default:null;column:metadata"` // メタデータ(JSON)
	CreatedAt    time.Time         `gorm:"<-:create"`                    // 登録日時
	UpdatedAt    time.Time         `gorm:""`                             // 更新日時
}
