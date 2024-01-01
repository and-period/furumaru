package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

type CategoryOrderBy string

const (
	CategoryOrderByName CategoryOrderBy = "name"
)

// Category - 商品種別情報
type Category struct {
	ID        string    `gorm:"primaryKey;<-:create"` // カテゴリID
	Name      string    `gorm:""`                     // カテゴリ名
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type Categories []*Category

type NewCategoryParams struct {
	Name string
}

func NewCategory(params *NewCategoryParams) *Category {
	return &Category{
		ID:   uuid.Base58Encode(uuid.New()),
		Name: params.Name,
	}
}

func (cs Categories) MapByName() map[string]*Category {
	res := make(map[string]*Category, len(cs))
	for _, c := range cs {
		res[c.Name] = c
	}
	return res
}
