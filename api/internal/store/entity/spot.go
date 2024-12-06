package entity

import (
	"errors"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
)

// SpotUserType - 投稿者の種別
type SpotUserType int32

const (
	SpotUserTypeUnknown     SpotUserType = 0
	SpotUserTypeUser        SpotUserType = 1 // ユーザー
	SpotUserTypeCoordinator SpotUserType = 2 // コーディネータ
	SpotUserTypeProducer    SpotUserType = 3 // 生産者
)

// Spot - スポット情報
type Spot struct {
	ID              string       `gorm:"primaryKey;<-:create"`             // スポットID
	TypeID          string       `gorm:"column:spot_type_id;default:null"` // 種別ID
	UserType        SpotUserType `gorm:""`                                 // 投稿者の種別
	UserID          string       `gorm:""`                                 // ユーザーID
	Name            string       `gorm:""`                                 // スポット名
	Description     string       `gorm:""`                                 // 説明
	ThumbnailURL    string       `gorm:""`                                 // サムネイル画像URL
	Longitude       float64      `gorm:""`                                 // 座標情報:経度
	Latitude        float64      `gorm:""`                                 // 座標情報:緯度
	PostalCode      string       `gorm:""`                                 // 郵便番号
	Prefecture      string       `gorm:"-"`                                // 都道府県
	PrefectureCode  int32        `gorm:"column:prefecture"`                // 都道府県コード
	City            string       `gorm:""`                                 // 市区町村
	AddressLine1    string       `gorm:""`                                 // 町名・番地
	AddressLine2    string       `gorm:""`                                 // ビル名・号室など
	Approved        bool         `gorm:""`                                 // 承認フラグ
	ApprovedAdminID string       `gorm:""`                                 // 承認した管理者ID
	CreatedAt       time.Time    `gorm:"<-:create"`                        // 登録日時
	UpdatedAt       time.Time    `gorm:""`                                 // 更新日時
}

type Spots []*Spot

type SpotParams struct {
	SpotTypeID   string
	UserType     SpotUserType
	UserID       string
	Name         string
	Description  string
	ThumbnailURL string
	Longitude    float64
	Latitude     float64
	PostalCode   string
	Prefecture   string
	City         string
	AddressLine1 string
	AddressLine2 string
}

func NewSpotByUser(params *SpotParams) (*Spot, error) {
	prefectureCode, err := codes.ToPrefectureValue(params.Prefecture)
	if err != nil {
		return nil, fmt.Errorf("entity: invalid prefecture name: %w", err)
	}
	res := &Spot{
		ID:              uuid.Base58Encode(uuid.New()),
		TypeID:          params.SpotTypeID,
		UserType:        SpotUserTypeUser,
		UserID:          params.UserID,
		Name:            params.Name,
		Description:     params.Description,
		ThumbnailURL:    params.ThumbnailURL,
		Approved:        true, // デフォルトで承認済みにする
		ApprovedAdminID: "",
		Longitude:       params.Longitude,
		Latitude:        params.Latitude,
		PostalCode:      params.PostalCode,
		Prefecture:      params.Prefecture,
		PrefectureCode:  prefectureCode,
		City:            params.City,
		AddressLine1:    params.AddressLine1,
		AddressLine2:    params.AddressLine2,
	}
	if err := res.Validate(); err != nil {
		return nil, err
	}
	return res, nil
}

func NewSpotByAdmin(params *SpotParams) (*Spot, error) {
	prefectureCode, err := codes.ToPrefectureValue(params.Prefecture)
	if err != nil {
		return nil, fmt.Errorf("entity: invalid prefecture name: %w", err)
	}
	res := &Spot{
		ID:              uuid.Base58Encode(uuid.New()),
		TypeID:          params.SpotTypeID,
		UserType:        params.UserType,
		UserID:          params.UserID,
		Name:            params.Name,
		Description:     params.Description,
		ThumbnailURL:    params.ThumbnailURL,
		Approved:        true,
		ApprovedAdminID: params.UserID,
		Longitude:       params.Longitude,
		Latitude:        params.Latitude,
		PostalCode:      params.PostalCode,
		Prefecture:      params.Prefecture,
		PrefectureCode:  prefectureCode,
		City:            params.City,
		AddressLine1:    params.AddressLine1,
		AddressLine2:    params.AddressLine2,
	}
	if err := res.Validate(); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Spot) Validate() error {
	if c.Longitude < -180 || 180 < c.Longitude {
		return errors.New("entity: longitude is invalid")
	}
	if c.Latitude < -90 || 90 < c.Latitude {
		return errors.New("entity: latitude is invalid")
	}
	return nil
}

func (c *Spot) Fill() (err error) {
	c.Prefecture, err = codes.ToPrefectureJapanese(c.PrefectureCode)
	return
}

func (cs Spots) TypeIDs() []string {
	return set.UniqBy(cs, func(c *Spot) string {
		return c.TypeID
	})
}

func (cs Spots) UserIDs() []string {
	return set.UniqBy(cs, func(c *Spot) string {
		return c.UserID
	})
}

func (cs Spots) Fill() error {
	for _, c := range cs {
		if err := c.Fill(); err != nil {
			return err
		}
	}
	return nil
}

func (cs Spots) GroupByUserType() map[SpotUserType]Spots {
	res := make(map[SpotUserType]Spots, 3)
	for _, c := range cs {
		if _, ok := res[c.UserType]; !ok {
			res[c.UserType] = make(Spots, 0, len(cs))
		}
		res[c.UserType] = append(res[c.UserType], c)
	}
	return res
}
