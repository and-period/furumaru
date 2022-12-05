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

// Activity - 注文イベントログ
type Activity struct {
	ID           string            `gorm:"primaryKey;<-:create"`         // 注文イベントログID
	OrderID      string            `gorm:""`                             // 注文履歴ID
	UserID       string            `gorm:""`                             // ユーザーID
	EventType    OrderEventType    `gorm:""`                             // イベントログ種別
	Detail       string            `gorm:""`                             // イベントログ詳細
	Metadata     map[string]string `gorm:"-"`                            // メタデータ
	MetadataJSON datatypes.JSON    `gorm:"default:null;column:metadata"` // メタデータ(JSON)
	CreatedAt    time.Time         `gorm:"<-:create"`                    // 登録日時
	UpdatedAt    time.Time         `gorm:""`                             // 更新日時
}

type Activities []*Activity

func (as Activities) GroupByOrderID() map[string]Activities {
	res := make(map[string]Activities, len(as))
	for _, a := range as {
		if _, ok := res[a.OrderID]; !ok {
			res[a.OrderID] = make(Activities, 0, len(as))
		}
		res[a.OrderID] = append(res[a.OrderID], a)
	}
	return res
}
