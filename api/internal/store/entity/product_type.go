package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

// ProductType - 品目情報
type ProductType struct {
	ID         string    `gorm:"primaryKey;<-:create"` // 品目ID
	Name       string    `gorm:""`                     // 品目名
	CategoryID string    `gorm:""`                     // カテゴリID
	CreatedAt  time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt  time.Time `gorm:""`                     // 更新日時
}

type ProductTypes []*ProductType

func NewProductType(name, categoryID string) *ProductType {
	return &ProductType{
		ID:         uuid.Base58Encode(uuid.New()),
		Name:       name,
		CategoryID: categoryID,
	}
}
