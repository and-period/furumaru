package entity

import (
	"errors"
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"gorm.io/gorm"
)

const DefaultShippingID = "default"

var (
	errInvalidShippingRateFormat     = errors.New("entity: invalid shipping rate format")
	errInvalidShippingRatePrefLength = errors.New("entity: unmatch shipping rate prefecture length")
	errNotUniqueShippingRateNumber   = errors.New("entity: shipping rate number must be unique")
)

// Shipping - 配送設定情報
type Shipping struct {
	ShippingRevision `gorm:"-"`
	ID               string         `gorm:"primaryKey;<-:create"` // 配送設定ID
	CoordinatorID    string         `gorm:""`                     // コーディネータID
	CreatedAt        time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt        time.Time      `gorm:""`                     // 更新日時
	DeletedAt        gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type Shippings []*Shipping

type NewShippingParams struct {
	CoordinatorID      string
	Box60Rates         ShippingRates
	Box60Refrigerated  int64
	Box60Frozen        int64
	Box80Rates         ShippingRates
	Box80Refrigerated  int64
	Box80Frozen        int64
	Box100Rates        ShippingRates
	Box100Refrigerated int64
	Box100Frozen       int64
	HasFreeShipping    bool
	FreeShippingRates  int64
}

func NewShipping(params *NewShippingParams) *Shipping {
	rparams := &NewShippingRevisionParams{
		ShippingID:         params.CoordinatorID,
		Box60Rates:         params.Box60Rates,
		Box60Refrigerated:  params.Box60Refrigerated,
		Box60Frozen:        params.Box60Frozen,
		Box80Rates:         params.Box80Rates,
		Box80Refrigerated:  params.Box80Refrigerated,
		Box80Frozen:        params.Box80Frozen,
		Box100Rates:        params.Box100Rates,
		Box100Refrigerated: params.Box100Refrigerated,
		Box100Frozen:       params.Box100Frozen,
		HasFreeShipping:    params.HasFreeShipping,
		FreeShippingRates:  params.FreeShippingRates,
	}
	revision := NewShippingRevision(rparams)
	return &Shipping{
		ID:               params.CoordinatorID, // PKはコーディネータと同一にする
		CoordinatorID:    params.CoordinatorID,
		ShippingRevision: *revision,
	}
}

func (s *Shipping) IsDefault() bool {
	return s.ID == DefaultShippingID
}

func (s *Shipping) Fill(revision *ShippingRevision) {
	s.ShippingRevision = *revision
}

func (ss Shippings) Fill(revisions map[string]*ShippingRevision) {
	for _, s := range ss {
		revision, ok := revisions[s.ID]
		if !ok {
			continue
		}
		s.Fill(revision)
	}
}

func (ss Shippings) IDs() []string {
	return set.UniqBy(ss, func(s *Shipping) string {
		return s.ID
	})
}

func (ss Shippings) CoordinatorIDs() []string {
	return set.UniqBy(ss, func(s *Shipping) string {
		return s.CoordinatorID
	})
}

func (ss Shippings) Map() map[string]*Shipping {
	res := make(map[string]*Shipping, len(ss))
	for _, s := range ss {
		res[s.ID] = s
	}
	return res
}
