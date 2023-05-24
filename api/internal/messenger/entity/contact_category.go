package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

// ContactCategory - 問い合わせ種別情報
type ContactCategory struct {
	ID        string    `gorm:"primaryKey;<-:create"` // カテゴリID
	Title     string    `gorm:""`                     // カテゴリ名
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type ContactCategories []*ContactCategory

func NewContactCategory(title string) *ContactCategory {
	return &ContactCategory{
		ID:    uuid.Base58Encode(uuid.New()),
		Title: title,
	}
}
