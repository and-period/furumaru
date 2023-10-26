package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
)

// ライブ配信情報
type Live struct {
	LiveProducts `gorm:"-"`
	ID           string    `gorm:"primaryKey;<-:create"` // ライブ配信ID
	ScheduleID   string    `gorm:""`                     // 開催スケジュールID
	ProducerID   string    `gorm:""`                     // 生産者ID
	ProductIDs   []string  `gorm:"-"`                    // 商品ID一覧
	Comment      string    `gorm:""`                     // コメント
	StartAt      time.Time `gorm:""`                     // 配信開始日時
	EndAt        time.Time `gorm:""`                     // 配信終了日時
	CreatedAt    time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt    time.Time `gorm:""`                     // 更新日時
}

type Lives []*Live

type NewLiveParams struct {
	ScheduleID string
	ProducerID string
	ProductIDs []string
	Comment    string
	StartAt    time.Time
	EndAt      time.Time
}

func NewLive(params *NewLiveParams) *Live {
	liveID := uuid.Base58Encode(uuid.New())
	return &Live{
		ID:           liveID,
		ScheduleID:   params.ScheduleID,
		ProducerID:   params.ProducerID,
		ProductIDs:   params.ProductIDs,
		Comment:      params.Comment,
		StartAt:      params.StartAt,
		EndAt:        params.EndAt,
		LiveProducts: NewLiveProducts(liveID, params.ProductIDs),
	}
}

func (l *Live) Fill(products LiveProducts) {
	l.ProductIDs = products.ProductIDs()
	l.LiveProducts = products
}

func (ls Lives) IDs() []string {
	return set.UniqBy(ls, func(l *Live) string {
		return l.ID
	})
}

func (ls Lives) ProducerIDs() []string {
	return set.UniqBy(ls, func(l *Live) string {
		return l.ProducerID
	})
}

func (ls Lives) ProductIDs() []string {
	res := set.NewEmpty[string](len(ls))
	for i := range ls {
		res.Add(ls[i].ProductIDs...)
	}
	return res.Slice()
}

func (ls Lives) Fill(products map[string]LiveProducts) {
	for i := range ls {
		ls[i].Fill(products[ls[i].ID])
	}
}

func (ls Lives) GroupByScheduleID() map[string]Lives {
	res := make(map[string]Lives, len(ls))
	for _, live := range ls {
		if _, ok := res[live.ScheduleID]; !ok {
			res[live.ScheduleID] = make(Lives, 0, len(ls))
		}
		res[live.ScheduleID] = append(res[live.ScheduleID], live)
	}
	return res
}
