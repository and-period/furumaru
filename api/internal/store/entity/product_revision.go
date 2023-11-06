package entity

import "time"

// ProductRevision - 商品変更履歴情報
type ProductRevision struct {
	ID        int64     `gorm:"primaryKey;<-:create"` // 変更履歴ID
	ProductID string    `gorm:""`                     // 商品ID
	Price     int64     `gorm:""`                     // 販売価格
	Cost      int64     `gorm:""`                     // 商品原価
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

func (ps ProductRevisions) Map() map[string]*ProductRevision {
	res := make(map[string]*ProductRevision, len(ps))
	for _, p := range ps {
		res[p.ProductID] = p
	}
	return res
}
