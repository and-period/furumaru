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
	errUnknownShippingSize           = errors.New("entity: unknown shipping size")
)

// ShippingSize - 配送種別
type ShippingType int32

const (
	ShippingTypeUnknown ShippingType = 0
	ShippingTypeNormal  ShippingType = 1 // 通常配送
	ShippingTypeFrozen  ShippingType = 2 // クール配送
	ShippingTypePickup  ShippingType = 3 // 店舗受け取り
)

// Shipping - 配送設定情報
type Shipping struct {
	ShippingRevision `gorm:"-"`
	ID               string         `gorm:"primaryKey;<-:create"` // 配送設定ID
	ShopID           string         `gorm:"default:null"`         // 店舗ID
	CoordinatorID    string         `gorm:""`                     // コーディネータID
	Name             string         `gorm:""`                     // 配送設定名
	InUse            bool           `gorm:""`                     // 使用中
	CreatedAt        time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt        time.Time      `gorm:""`                     // 更新日時
	DeletedAt        gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type Shippings []*Shipping

type NewShippingParams struct {
	ShopID            string
	Name              string
	CoordinatorID     string
	Box60Rates        ShippingRates
	Box60Frozen       int64
	Box80Rates        ShippingRates
	Box80Frozen       int64
	Box100Rates       ShippingRates
	Box100Frozen      int64
	HasFreeShipping   bool
	FreeShippingRates int64
	InUse             bool
}

func (t ShippingType) String() string {
	switch t {
	case ShippingTypeNormal:
		return "通常配送"
	case ShippingTypeFrozen:
		return "クール配送"
	default:
		return ""
	}
}

func NewShipping(params *NewShippingParams) *Shipping {
	rparams := &NewShippingRevisionParams{
		ShippingID:        params.CoordinatorID,
		Box60Rates:        params.Box60Rates,
		Box60Frozen:       params.Box60Frozen,
		Box80Rates:        params.Box80Rates,
		Box80Frozen:       params.Box80Frozen,
		Box100Rates:       params.Box100Rates,
		Box100Frozen:      params.Box100Frozen,
		HasFreeShipping:   params.HasFreeShipping,
		FreeShippingRates: params.FreeShippingRates,
	}
	revision := NewShippingRevision(rparams)
	return &Shipping{
		ID:               params.CoordinatorID, // PKはコーディネータと同一にする
		ShopID:           params.ShopID,
		CoordinatorID:    params.CoordinatorID,
		Name:             params.Name,
		InUse:            params.InUse,
		ShippingRevision: *revision,
	}
}

func (s *Shipping) IsDefault() bool {
	return s.ID == DefaultShippingID
}

func (s *Shipping) CalcShippingFee(
	shippingSize ShippingSize, shippingType ShippingType, total int64, prefectureCode int32,
) (int64, error) {
	if s.HasFreeShipping && total >= s.FreeShippingRates {
		return 0, nil // 送料設定による無料配送
	}
	var (
		rate       *ShippingRate
		additional int64
		err        error
	)
	switch shippingSize {
	case ShippingSize60:
		if shippingType == ShippingTypeFrozen {
			additional = s.Box60Frozen
		}
		rate, err = s.Box60Rates.Find(prefectureCode)
	case ShippingSize80:
		if shippingType == ShippingTypeFrozen {
			additional = s.Box80Frozen
		}
		rate, err = s.Box80Rates.Find(prefectureCode)
	case ShippingSize100:
		if shippingType == ShippingTypeFrozen {
			additional = s.Box100Frozen
		}
		rate, err = s.Box100Rates.Find(prefectureCode)
	default:
		return 0, errUnknownShippingSize
	}
	if err != nil {
		return 0, err
	}
	return rate.Price + additional, nil
}

func (s *Shipping) Fill(revision *ShippingRevision) {
	s.ShippingRevision.Fill()
	s.ShippingRevision = *revision
}

func (ss Shippings) Fill(revisions map[string]*ShippingRevision) {
	for _, s := range ss {
		revision, ok := revisions[s.ID]
		if !ok {
			revision = &ShippingRevision{ShippingID: s.ID}
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
