package entity

import "time"

// ProductType - 品目情報
type ProductType struct {
	ID         string    `gorm:"primaryKey;<-:craete"` // 品目ID
	Name       string    `gorm:""`                     // 品目名
	CategoryID string    `gorm:""`                     // カテゴリID
	CreatedAt  time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt  time.Time `gorm:""`                     // 更新日時
}
