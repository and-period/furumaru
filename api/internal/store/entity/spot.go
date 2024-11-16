package entity

import (
	"errors"
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

// SpotUserType - 投稿者の種別
type SpotUserType int32

const (
	SpotUserTypeUnknown SpotUserType = 0
	SpotUserTypeUser    SpotUserType = 1 // ユーザー
	SpotUserTypeAdmin   SpotUserType = 2 // 管理者
)

// Spot - スポット情報
type Spot struct {
	ID              string       `gorm:"primaryKey;<-:create"` // スポットID
	UserType        SpotUserType `gorm:""`                     // 投稿者の種別
	UserID          string       `gorm:""`                     // ユーザーID
	Name            string       `gorm:""`                     // スポット名
	Description     string       `gorm:""`                     // 説明
	ThumbnailURL    string       `gorm:""`                     // サムネイル画像URL
	Longitude       float64      `gorm:""`                     // 座標情報:経度
	Latitude        float64      `gorm:""`                     // 座標情報:緯度
	Approved        bool         `gorm:""`                     // 承認フラグ
	ApprovedAdminID string       `gorm:""`                     // 承認した管理者ID
	CreatedAt       time.Time    `gorm:"<-:create"`            // 登録日時
	UpdatedAt       time.Time    `gorm:""`                     // 更新日時
}

type Spots []*Spot

type SpotParams struct {
	UserID       string
	Name         string
	Description  string
	ThumbnailURL string
	Longitude    float64
	Latitude     float64
}

func NewSpotByUser(params *SpotParams) (*Spot, error) {
	res := &Spot{
		ID:              uuid.Base58Encode(uuid.New()),
		UserType:        SpotUserTypeUser,
		UserID:          params.UserID,
		Name:            params.Name,
		Description:     params.Description,
		ThumbnailURL:    params.ThumbnailURL,
		Approved:        true, // デフォルトで承認済みにする
		ApprovedAdminID: "",
		Longitude:       params.Longitude,
		Latitude:        params.Latitude,
	}
	if err := res.Validate(); err != nil {
		return nil, err
	}
	return res, nil
}

func NewSpotByAdmin(params *SpotParams) (*Spot, error) {
	res := &Spot{
		ID:              uuid.Base58Encode(uuid.New()),
		UserType:        SpotUserTypeAdmin,
		UserID:          params.UserID,
		Name:            params.Name,
		Description:     params.Description,
		ThumbnailURL:    params.ThumbnailURL,
		Approved:        true,
		ApprovedAdminID: params.UserID,
		Longitude:       params.Longitude,
		Latitude:        params.Latitude,
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
