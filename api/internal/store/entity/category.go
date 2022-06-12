package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

// Category - 商品種別情報
type Category struct {
	ID        string    `gorm:"primaryKey;<-:create"` // カテゴリID
	Name      string    `gorm:""`                     // カテゴリ名
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type Categories []*Category

func NewCategory(name string) *Category {
	return &Category{
		ID:   uuid.Base58Encode(uuid.New()),
		Name: name,
	}
}
