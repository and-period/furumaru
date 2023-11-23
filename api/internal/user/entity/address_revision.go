package entity

import (
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/jinzhu/copier"
)

// AddressRevision - 住所変更履歴情報
type AddressRevision struct {
	ID             int64     `gorm:"primaryKey;<-:create"` // 変更履歴ID
	AddressID      string    `gorm:""`                     // 住所ID
	Lastname       string    `gorm:""`                     // 姓
	Firstname      string    `gorm:""`                     // 名
	LastnameKana   string    `gorm:""`                     // 姓（かな）
	FirstnameKana  string    `gorm:""`                     // 名（かな）
	PostalCode     string    `gorm:""`                     // 郵便番号
	Prefecture     string    `gorm:"-"`                    // 都道府県
	PrefectureCode int32     `gorm:"column:prefecture"`    // 都道府県コード
	City           string    `gorm:""`                     // 市区町村
	AddressLine1   string    `gorm:""`                     // 町名・番地
	AddressLine2   string    `gorm:""`                     // ビル名・号室など
	PhoneNumber    string    `gorm:"default:null"`         // 電話番号
	CreatedAt      time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time `gorm:""`                     // 更新日時
}

type AddressRevisions []*AddressRevision

type NewAddressRevisionParams struct {
	AddressID      string
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

func NewAddressRevision(params *NewAddressRevisionParams) (*AddressRevision, error) {
	prefecture, err := codes.ToPrefectureJapanese(params.PrefectureCode)
	if err != nil {
		return nil, err
	}
	return &AddressRevision{
		AddressID:      params.AddressID,
		Lastname:       params.Lastname,
		Firstname:      params.Firstname,
		LastnameKana:   params.LastnameKana,
		FirstnameKana:  params.FirstnameKana,
		PostalCode:     params.PostalCode,
		Prefecture:     prefecture,
		PrefectureCode: params.PrefectureCode,
		City:           params.City,
		AddressLine1:   params.AddressLine1,
		AddressLine2:   params.AddressLine2,
		PhoneNumber:    params.PhoneNumber,
	}, nil
}

func (r *AddressRevision) Fill() {
	r.Prefecture, _ = codes.ToPrefectureJapanese(r.PrefectureCode)
}

func (rs AddressRevisions) AddressIDs() []string {
	return set.UniqBy(rs, func(r *AddressRevision) string {
		return r.AddressID
	})
}

func (rs AddressRevisions) MapByAddressID() map[string]*AddressRevision {
	res := make(map[string]*AddressRevision, len(rs))
	for _, r := range rs {
		res[r.AddressID] = r
	}
	return res
}

func (rs AddressRevisions) Fill() {
	for i := range rs {
		rs[i].Fill()
	}
}

func (rs AddressRevisions) Merge(addresses map[string]*Address) (Addresses, error) {
	res := make(Addresses, 0, len(rs))
	for _, r := range rs {
		address := &Address{}
		base, ok := addresses[r.AddressID]
		if !ok {
			base = &Address{ID: r.AddressID}
		}
		opt := copier.Option{IgnoreEmpty: true, DeepCopy: true}
		if err := copier.CopyWithOption(&address, &base, opt); err != nil {
			return nil, err
		}
		address.AddressRevision = *r
		res = append(res, address)
	}
	return res, nil
}
