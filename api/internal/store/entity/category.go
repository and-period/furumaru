package entity

import "time"

// Category - 商品種別情報
type Category struct {
	ID        string    `gorm:"primaryKey;<-:craete"` // カテゴリID
	Name      string    `gorm:""`                     // カテゴリ名
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}
