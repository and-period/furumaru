package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/jinzhu/copier"
)

// ProductRevision - 商品変更履歴情報
type ProductRevision struct {
	ID        int64     `gorm:"primaryKey;<-:create"` // 変更履歴ID
	ProductID string    `gorm:""`                     // 商品ID
	Price     int64     `gorm:""`                     // 販売価格(税込)
	Cost      int64     `gorm:""`                     // 商品原価(税込)
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type ProductRevisions []*ProductRevision

type NewProductRevisionParams struct {
	ProductID string
	Price     int64
	Cost      int64
}

func NewProductRevision(params *NewProductRevisionParams) *ProductRevision {
	return &ProductRevision{
		ProductID: params.ProductID,
		Price:     params.Price,
		Cost:      params.Cost,
	}
}

func (rs ProductRevisions) ProductIDs() []string {
	return set.UniqBy(rs, func(r *ProductRevision) string {
		return r.ProductID
	})
}

func (rs ProductRevisions) MapByProductID() map[string]*ProductRevision {
	res := make(map[string]*ProductRevision, len(rs))
	for _, r := range rs {
		res[r.ProductID] = r
	}
	return res
}

func (rs ProductRevisions) Merge(products map[string]*Product) (Products, error) {
	res := make(Products, 0, len(rs))
	for _, r := range rs {
		product := &Product{}
		base, ok := products[r.ProductID]
		if !ok {
			base = &Product{ID: r.ProductID}
		}
		opt := copier.Option{IgnoreEmpty: true, DeepCopy: true}
		if err := copier.CopyWithOption(&product, &base, opt); err != nil {
			return nil, err
		}
		product.ProductRevision = *r
		res = append(res, product)
	}
	return res, nil
}
