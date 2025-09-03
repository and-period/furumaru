package entity

import (
	"strings"
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

// Address - 住所情報
type Address struct {
	AddressRevision `gorm:"-"`
	ID              string         `gorm:"primaryKey;<-:create"` // 住所ID
	UserID          string         `gorm:""`                     // ユーザーID
	IsDefault       bool           `gorm:""`                     // デフォルト設定フラグ
	CreatedAt       time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt       time.Time      `gorm:""`                     // 更新日時
	DeletedAt       gorm.DeletedAt `gorm:"default:null"`         // 退会日時
}

type Addresses []*Address

type NewAddressParams struct {
	UserID         string
	IsDefault      bool
	Lastname       string
	Firstname      string
	LastnameKana   string
	FirstnameKana  string
	PostalCode     string
	PrefectureCode int32
	City           string
	AddressLine1   string
	AddressLine2   string
	PhoneNumber    string
}

func NewAddress(params *NewAddressParams) (*Address, error) {
	addressID := uuid.Base58Encode(uuid.New())
	rparams := &NewAddressRevisionParams{
		AddressID:      addressID,
		Lastname:       params.Lastname,
		Firstname:      params.Firstname,
		LastnameKana:   params.LastnameKana,
		FirstnameKana:  params.FirstnameKana,
		PostalCode:     params.PostalCode,
		PrefectureCode: params.PrefectureCode,
		City:           params.City,
		AddressLine1:   params.AddressLine1,
		AddressLine2:   params.AddressLine2,
		PhoneNumber:    params.PhoneNumber,
	}
	revision, err := NewAddressRevision(rparams)
	if err != nil {
		return nil, err
	}
	return &Address{
		ID:              addressID,
		UserID:          params.UserID,
		IsDefault:       params.IsDefault,
		AddressRevision: *revision,
	}, nil
}

func (a *Address) Name() string {
	if a == nil {
		return ""
	}
	str := strings.Join([]string{a.Lastname, a.Firstname}, " ")
	return strings.TrimSpace(str)
}

func (a *Address) NameKana() string {
	if a == nil {
		return ""
	}
	str := strings.Join([]string{a.LastnameKana, a.FirstnameKana}, " ")
	return strings.TrimSpace(str)
}

func (a *Address) FullPath() string {
	str := strings.Join([]string{a.Prefecture, a.City, a.AddressLine1, a.AddressLine2}, " ")
	return strings.TrimSpace(str)
}

func (a *Address) ShortPath() string {
	str := strings.Join([]string{a.Prefecture, a.City, a.AddressLine1}, " ")
	return strings.TrimSpace(str)
}

func (a *Address) Fill(revision *AddressRevision) {
	a.AddressRevision = *revision
}

func (as Addresses) IDs() []string {
	return set.UniqBy(as, func(a *Address) string {
		return a.ID
	})
}

func (as Addresses) Map() map[string]*Address {
	res := make(map[string]*Address, len(as))
	for _, a := range as {
		res[a.ID] = a
	}
	return res
}

func (as Addresses) MapByRevision() map[int64]*Address {
	res := make(map[int64]*Address, len(as))
	for _, a := range as {
		res[a.AddressRevision.ID] = a
	}
	return res
}

func (as Addresses) MapByUserID() map[string]*Address {
	res := make(map[string]*Address, len(as))
	for _, a := range as {
		res[a.UserID] = a
	}
	return res
}

func (as Addresses) Fill(revisions map[string]*AddressRevision) {
	for _, a := range as {
		revision, ok := revisions[a.ID]
		if !ok {
			revision = &AddressRevision{AddressID: a.ID}
		}
		a.Fill(revision)
	}
}
