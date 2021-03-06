package entity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/and-period/furumaru/api/internal/store/codes"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/datatypes"
)

var (
	errInvalidShippingRateFormat     = errors.New("entity: invalid shipping rate format")
	errInvalidShippingRatePrefLength = errors.New("entity: unmatch shipping rate prefecture length")
	errNotUniqueShippingRateNumber   = errors.New("entity: shipping rate number must be unique")
)

// Shipping - 配送設定情報
type Shipping struct {
	ID                 string         `gorm:"primaryKey;<-:create"`             // 配送設定ID
	Name               string         `gorm:""`                                 // 配送設定名
	Box60Rates         ShippingRates  `gorm:"-"`                                // 箱サイズ60の通常便配送料一覧
	Box60RatesJSON     datatypes.JSON `gorm:"default:null;column:box60_rates"`  // 箱サイズ60の通常便配送料一覧(JSON)
	Box60Refrigerated  int64          `gorm:""`                                 // 箱サイズ60の冷蔵便追加配送料
	Box60Frozen        int64          `gorm:""`                                 // 箱サイズ60の冷凍便追加配送料
	Box80Rates         ShippingRates  `gorm:"-"`                                // 箱サイズ80の通常便配送料一覧
	Box80RatesJSON     datatypes.JSON `gorm:"default:null;column:box80_rates"`  // 箱サイズ80の通常便配送料一覧(JSON)
	Box80Refrigerated  int64          `gorm:""`                                 // 箱サイズ80の冷蔵便追加配送料
	Box80Frozen        int64          `gorm:""`                                 // 箱サイズ80の冷凍便追加配送料
	Box100Rates        ShippingRates  `gorm:"-"`                                // 箱サイズ100の通常便配送料一覧
	Box100RatesJSON    datatypes.JSON `gorm:"default:null;column:box100_rates"` // 箱サイズ100の通常便配送料一覧(JSON)
	Box100Refrigerated int64          `gorm:""`                                 // 箱サイズ100の冷蔵便追加配送料
	Box100Frozen       int64          `gorm:""`                                 // 箱サイズ100の冷凍便追加配送料
	HasFreeShipping    bool           `gorm:""`                                 // 送料無料オプションの有無
	FreeShippingRates  int64          `gorm:""`                                 // 送料無料になる金額
	CreatedAt          time.Time      `gorm:"<-:create"`                        // 登録日時
	UpdatedAt          time.Time      `gorm:""`                                 // 更新日時
}

type Shippings []*Shipping

// ShippingRate - 配送料金情報
type ShippingRate struct {
	Number      int64   `json:"number"`      // No.
	Name        string  `json:"name"`        // 配送料金設定名
	Price       int64   `json:"price"`       // 配送料金
	Prefectures []int64 `json:"prefectures"` // 対象都道府県一覧
}

type ShippingRates []*ShippingRate

type NewShippingParams struct {
	Name               string
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
	return &Shipping{
		ID:                 uuid.Base58Encode(uuid.New()),
		Name:               params.Name,
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
}

func (s *Shipping) Fill() error {
	var box60Rates, box80Rates, box100Rates ShippingRates
	if err := json.Unmarshal(s.Box60RatesJSON, &box60Rates); err != nil {
		return err
	}
	if err := json.Unmarshal(s.Box80RatesJSON, &box80Rates); err != nil {
		return err
	}
	if err := json.Unmarshal(s.Box100RatesJSON, &box100Rates); err != nil {
		return err
	}
	s.Box60Rates = box60Rates
	s.Box80Rates = box80Rates
	s.Box100Rates = box100Rates
	return nil
}

func (s *Shipping) FillJSON() error {
	box60Rates, err := s.Box60Rates.Marshal()
	if err != nil {
		return err
	}
	box80Rates, err := s.Box80Rates.Marshal()
	if err != nil {
		return err
	}
	box100Rates, err := s.Box100Rates.Marshal()
	if err != nil {
		return err
	}
	s.Box60RatesJSON = datatypes.JSON(box60Rates)
	s.Box80RatesJSON = datatypes.JSON(box80Rates)
	s.Box100RatesJSON = datatypes.JSON(box100Rates)
	return nil
}

func (ss Shippings) Fill() error {
	for i := range ss {
		if err := ss[i].Fill(); err != nil {
			return err
		}
	}
	return nil
}

func NewShippingRate(num int64, name string, price int64, prefs []int64) *ShippingRate {
	return &ShippingRate{
		Number:      num,
		Name:        name,
		Price:       price,
		Prefectures: prefs,
	}
}

func (rs ShippingRates) Validate() error {
	var total int
	set := set.New(len(rs))
	for i := range rs {
		if rs[i].Number < 1 { // No.の形式チェック
			return errInvalidShippingRateFormat
		}
		if rs[i].Price < 0 { // 配送料金の形式チェック
			return errInvalidShippingRateFormat
		}
		if set.FindOrAdd(rs[i].Number) { // No.の重複チェック
			return errNotUniqueShippingRateNumber
		}
		if err := codes.ValidatePrefectureValues(rs[i].Prefectures...); err != nil { // 都道府県の存在性チェック
			return err
		}
		total += len(rs[i].Prefectures)
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
