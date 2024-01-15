package entity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/jinzhu/copier"
	"gorm.io/datatypes"
)

var errNotFoundShippingRate = errors.New("entity: not found shipping rate")

// ShippingRevision - 配送設定変更履歴情報
type ShippingRevision struct {
	ID                int64          `gorm:"primarykey;<-:create"`             // 変更履歴ID
	ShippingID        string         `gorm:""`                                 // 配送設定ID
	Box60Rates        ShippingRates  `gorm:"-"`                                // 箱サイズ60の通常便配送料一覧
	Box60RatesJSON    datatypes.JSON `gorm:"default:null;column:box60_rates"`  // 箱サイズ60の通常便配送料一覧(JSON)
	Box60Frozen       int64          `gorm:""`                                 // 箱サイズ60の冷凍便追加配送料(税込)
	Box80Rates        ShippingRates  `gorm:"-"`                                // 箱サイズ80の通常便配送料一覧
	Box80RatesJSON    datatypes.JSON `gorm:"default:null;column:box80_rates"`  // 箱サイズ80の通常便配送料一覧(JSON)
	Box80Frozen       int64          `gorm:""`                                 // 箱サイズ80の冷凍便追加配送料(税込)
	Box100Rates       ShippingRates  `gorm:"-"`                                // 箱サイズ100の通常便配送料一覧
	Box100RatesJSON   datatypes.JSON `gorm:"default:null;column:box100_rates"` // 箱サイズ100の通常便配送料一覧(JSON)
	Box100Frozen      int64          `gorm:""`                                 // 箱サイズ100の冷凍便追加配送料(税込)
	HasFreeShipping   bool           `gorm:""`                                 // 送料無料オプションの有無
	FreeShippingRates int64          `gorm:""`                                 // 送料無料になる金額(税込)
	CreatedAt         time.Time      `gorm:"<-:create"`                        // 登録日時
	UpdatedAt         time.Time      `gorm:""`                                 // 更新日時
}

type ShippingRevisions []*ShippingRevision

// ShippingRate - 配送料金情報
type ShippingRate struct {
	Number          int64   `json:"number"`      // No.
	Name            string  `json:"name"`        // 配送料金設定名
	Price           int64   `json:"price"`       // 配送料金(税込)
	PrefectureCodes []int32 `json:"prefectures"` // 対象都道府県一覧
}

type ShippingRates []*ShippingRate

type NewShippingRevisionParams struct {
	ShippingID        string
	Box60Rates        ShippingRates
	Box60Frozen       int64
	Box80Rates        ShippingRates
	Box80Frozen       int64
	Box100Rates       ShippingRates
	Box100Frozen      int64
	HasFreeShipping   bool
	FreeShippingRates int64
}

func NewShippingRevision(params *NewShippingRevisionParams) *ShippingRevision {
	return &ShippingRevision{
		ShippingID:        params.ShippingID,
		Box60Rates:        params.Box60Rates,
		Box60Frozen:       params.Box60Frozen,
		Box80Rates:        params.Box80Rates,
		Box80Frozen:       params.Box80Frozen,
		Box100Rates:       params.Box100Rates,
		Box100Frozen:      params.Box100Frozen,
		HasFreeShipping:   params.HasFreeShipping,
		FreeShippingRates: params.FreeShippingRates,
	}
}

func (r *ShippingRevision) Fill() (err error) {
	if r.Box60Rates, err = r.unmarshalRates(r.Box60RatesJSON); err != nil {
		return err
	}
	if r.Box80Rates, err = r.unmarshalRates(r.Box80RatesJSON); err != nil {
		return err
	}
	if r.Box100Rates, err = r.unmarshalRates(r.Box100RatesJSON); err != nil {
		return err
	}
	return nil
}

func (r *ShippingRevision) unmarshalRates(b []byte) (ShippingRates, error) {
	if b == nil {
		return ShippingRates{}, nil
	}
	var rates ShippingRates
	return rates, json.Unmarshal(b, &rates)
}

func (r *ShippingRevision) FillJSON() error {
	box60Rates, err := r.Box60Rates.Marshal()
	if err != nil {
		return err
	}
	box80Rates, err := r.Box80Rates.Marshal()
	if err != nil {
		return err
	}
	box100Rates, err := r.Box100Rates.Marshal()
	if err != nil {
		return err
	}
	r.Box60RatesJSON = datatypes.JSON(box60Rates)
	r.Box80RatesJSON = datatypes.JSON(box80Rates)
	r.Box100RatesJSON = datatypes.JSON(box100Rates)
	return nil
}

func (rs ShippingRevisions) ShippingIDs() []string {
	return set.UniqBy(rs, func(r *ShippingRevision) string {
		return r.ShippingID
	})
}

func (rs ShippingRevisions) MapByShippingID() map[string]*ShippingRevision {
	res := make(map[string]*ShippingRevision, len(rs))
	for _, r := range rs {
		res[r.ShippingID] = r
	}
	return res
}

func (rs ShippingRevisions) Fill() error {
	for i := range rs {
		if err := rs[i].Fill(); err != nil {
			return err
		}
	}
	return nil
}

func (rs ShippingRevisions) Merge(shippings map[string]*Shipping) (Shippings, error) {
	res := make(Shippings, 0, len(rs))
	for _, r := range rs {
		shipping := &Shipping{}
		base, ok := shippings[r.ShippingID]
		if !ok {
			base = &Shipping{ID: r.ShippingID}
		}
		opt := copier.Option{IgnoreEmpty: true, DeepCopy: true}
		if err := copier.CopyWithOption(&shipping, &base, opt); err != nil {
			return nil, err
		}
		shipping.ShippingRevision = *r
		res = append(res, shipping)
	}
	return res, nil
}

func NewShippingRate(num int64, name string, price int64, prefs []int32) *ShippingRate {
	return &ShippingRate{
		Number:          num,
		Name:            name,
		Price:           price,
		PrefectureCodes: prefs,
	}
}

func (rs ShippingRates) Find(prefectureCode int32) (*ShippingRate, error) {
	for _, rate := range rs {
		set := set.New(rate.PrefectureCodes...)
		if set.Contains(prefectureCode) {
			return rate, nil
		}
	}
	return nil, errNotFoundShippingRate
}

func (rs ShippingRates) Validate() error {
	var total int
	set := set.NewEmpty[int64](len(rs))
	for i := range rs {
		if rs[i].Number < 1 { // No.の形式チェック
			return errInvalidShippingRateFormat
		}
		if rs[i].Price < 0 { // 配送料金の形式チェック
			return errInvalidShippingRateFormat
		}
		if _, exists := set.FindOrAdd(rs[i].Number); exists { // No.の重複チェック
			return errNotUniqueShippingRateNumber
		}
		if err := codes.ValidatePrefectureValues(rs[i].PrefectureCodes...); err != nil { // 都道府県の存在性チェック
			return err
		}
		total += len(rs[i].PrefectureCodes)
	}
	if total != len(codes.PrefectureValues) { // 都道府県が全て指定されているかのチェック(重複チェック含め)
		return errInvalidShippingRatePrefLength
	}
	return nil
}

func (rs ShippingRates) Marshal() ([]byte, error) {
	if len(rs) == 0 {
		return []byte{}, nil
	}
	return json.Marshal(rs)
}
