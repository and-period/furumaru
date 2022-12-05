package entity

import "time"

// Cart - カート情報
type Cart struct {
	ID         string       `gorm:"primaryKey;<-:create"` // カートID
	UserID     string       `gorm:""`                     // ユーザーID
	ScheduleID string       `gorm:"default:null"`         // マルシェ開催スケジュールID
	BoxNumber  int64        `gorm:""`                     // 箱の通番
	BoxType    DeliveryType `gorm:""`                     // 箱の種別
	BoxSize    ShippingSize `gorm:""`                     // 箱のサイズ
	CreatedAt  time.Time    `gorm:"<-:create"`            // 登録日時
	UpdatedAt  time.Time    `gorm:""`                     // 更新日時
}

type Carts []*Cart
