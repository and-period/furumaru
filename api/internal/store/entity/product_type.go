package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

type ProductTypeOrderBy string

const (
	ProductTypeOrderByName ProductTypeOrderBy = "name"
)

// ProductType - 品目情報
type ProductType struct {
	ID         string    `gorm:"primaryKey;<-:create"` // 品目ID
	Name       string    `gorm:""`                     // 品目名
	IconURL    string    `gorm:""`                     // アイコンURL
	CategoryID string    `gorm:""`                     // カテゴリID
	CreatedAt  time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt  time.Time `gorm:""`                     // 更新日時
}

type ProductTypes []*ProductType

func NewProductType(name, iconURL, categoryID string) *ProductType {
	return &ProductType{
		ID:         uuid.Base58Encode(uuid.New()),
		Name:       name,
		IconURL:    iconURL,
		CategoryID: categoryID,
	}
}
