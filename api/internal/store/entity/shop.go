package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

// Shop - 店舗情報
type Shop struct {
	ID             string         `gorm:"primaryKey;<-:create"` // 店舗ID
	CoordinatorID  string         `gorm:""`                     // コーディネータID
	ProducerIDs    []string       `gorm:"-"`                    // 生産者ID一覧
	ProductTypeIDs []string       `gorm:"-"`                    // 取り扱い商品種別ID一覧
	BusinessDays   []time.Weekday `gorm:"-"`                    // 営業曜日(発送可能日)一覧
	Name           string         `gorm:""`                     // 店舗名
	Activated      bool           `gorm:""`                     // 有効フラグ
	CreatedAt      time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time      `gorm:""`                     // 更新日時
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type Shops []*Shop

type ShopParams struct {
	CoordinatorID  string
	ProductTypeIDs []string
	BusinessDays   []time.Weekday
	Name           string
}

func NewShop(params *ShopParams) *Shop {
	return &Shop{
		ID:             uuid.Base58Encode(uuid.New()),
		CoordinatorID:  params.CoordinatorID,
		ProductTypeIDs: params.ProductTypeIDs,
		BusinessDays:   params.BusinessDays,
		Name:           params.Name,
		Activated:      true,
	}
}

func (s *Shop) Fill(producers ShopProducers) {
	s.ProducerIDs = producers.ProducerIDs()
}

func (ss Shops) IDs() []string {
	return set.UniqBy(ss, func(s *Shop) string {
		return s.ID
	})
}

func (ss Shops) Fill(producers map[string]ShopProducers) {
	for _, s := range ss {
		s.Fill(producers[s.ID])
	}
}
