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
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type LiveProducts []*LiveProduct

func NewLiveProduct(liveID, productID string) *LiveProduct {
	return &LiveProduct{
		LiveID:    liveID,
		ProductID: productID,
	}
}

func NewLiveProducts(liveID string, productIDs []string) LiveProducts {
	res := make(LiveProducts, len(productIDs))
	for i := range productIDs {
		res[i] = NewLiveProduct(liveID, productIDs[i])
	}
	return res
}

func (ps LiveProducts) ProductIDs() []string {
	res := set.NewEmpty[string](len(ps))
	for i := range ps {
		res.Add(ps[i].ProductID)
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

func (ps LiveProducts) SortByCreatedAt() LiveProducts {
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].CreatedAt.Before(ps[j].CreatedAt)
	})
	return ps
}
