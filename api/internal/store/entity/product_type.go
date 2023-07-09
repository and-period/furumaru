package entity

import (
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/datatypes"
)

type ProductTypeOrderBy string

const (
	ProductTypeOrderByName ProductTypeOrderBy = "name"
)

// ProductType - 品目情報
type ProductType struct {
	ID         string         `gorm:"primaryKey;<-:create"`      // 品目ID
	Name       string         `gorm:""`                          // 品目名
	IconURL    string         `gorm:""`                          // アイコンURL
	Icons      common.Images  `gorm:"-"`                         // アイコン一覧(リサイズ済み)
	IconsJSON  datatypes.JSON `gorm:"default:null;column:icons"` // アイコン一覧(JSON)
	CategoryID string         `gorm:""`                          // カテゴリID
	CreatedAt  time.Time      `gorm:"<-:create"`                 // 登録日時
	UpdatedAt  time.Time      `gorm:""`                          // 更新日時
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

func (t *ProductType) Fill() error {
	icons, err := common.NewImagesFromBytes(t.IconsJSON)
	if err != nil {
		return err
	}
	t.Icons = icons
	return nil
}

func (ts ProductTypes) Fill() error {
	for i := range ts {
		if err := ts[i].Fill(); err != nil {
			return err
		}
	}
	return nil
}

func (ts ProductTypes) CategoryIDs() []string {
	return set.UniqBy(ts, func(t *ProductType) string {
		return t.CategoryID
	})
}
