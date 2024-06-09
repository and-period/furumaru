package entity

import (
	"sort"
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
)

// ライブ配信関連商品情報
type LiveProduct struct {
	LiveID    string    `gorm:"primaryKey;<-:create"` // ライブ配信ID
	ProductID string    `gorm:"primaryKey;<-:create"` // 商品ID
	Priority  int64     `gorm:"default:0"`            // 表示優先度
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type LiveProducts []*LiveProduct

type NewLiveProductParams struct {
	LiveID    string
	ProductID string
	Priority  int64
}

type NewLiveProductsParams struct {
	LiveID     string
	ProductIDs []string
}

func NewLiveProduct(params *NewLiveProductParams) *LiveProduct {
	return &LiveProduct{
		LiveID:    params.LiveID,
		ProductID: params.ProductID,
		Priority:  params.Priority,
	}
}

func NewLiveProducts(liveID string, productIDs []string) LiveProducts {
	res := make(LiveProducts, len(productIDs))
	for i := range productIDs {
		params := &NewLiveProductParams{
			LiveID:    liveID,
			ProductID: productIDs[i],
			Priority:  int64(i + 1),
		}
		res[i] = NewLiveProduct(params)
	}
	return res
}

func (ps LiveProducts) ProductIDs() []string {
	res := set.NewEmpty[string](len(ps))
	orderedPs := ps.SortByPrimary()
	for i := range orderedPs {
		res.Add(orderedPs[i].ProductID)
	}
	return res.Slice()
}

func (ps LiveProducts) GroupByLiveID() map[string]LiveProducts {
	res := make(map[string]LiveProducts, len(ps))
	for _, p := range ps {
		if _, ok := res[p.LiveID]; !ok {
			res[p.LiveID] = make(LiveProducts, 0, len(ps))
		}
		res[p.LiveID] = append(res[p.LiveID], p)
	}
	return res
}

func (ps LiveProducts) SortByPrimary() LiveProducts {
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].Priority < ps[j].Priority
	})
	return ps
}
