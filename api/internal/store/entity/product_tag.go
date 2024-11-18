package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

// ProductTag - 商品タグ情報
type ProductTag struct {
	ID        string    `gorm:"primaryKey;<-:create"` // 商品タグID
	Name      string    `gorm:""`                     // 商品タグ名
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type ProductTags []*ProductTag

func NewProductTag(name string) *ProductTag {
	return &ProductTag{
		ID:   uuid.Base58Encode(uuid.New()),
		Name: name,
	}
}

func (ts ProductTags) IDs() []string {
	res := make([]string, len(ts))
	for i := range ts {
		res[i] = ts[i].ID
	}
	return res
}
