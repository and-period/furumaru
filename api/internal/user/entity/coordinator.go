package entity

import (
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
)

// Coordinator - コーディネータ情報
type Coordinator struct {
	Admin             `gorm:"-"`
	AdminID           string    `gorm:"primaryKey;<-:create"` // 管理者ID
	PhoneNumber       string    `gorm:""`                     // 電話番号
	Username          string    `gorm:""`                     // 表示名
	Profile           string    `gorm:""`                     // 紹介文
	ThumbnailURL      string    `gorm:""`                     // サムネイルURL
	HeaderURL         string    `gorm:""`                     // ヘッダー画像URL
	PromotionVideoURL string    `gorm:""`                     // 紹介動画URL
	BonusVideoURL     string    `gorm:""`                     // 購入特典動画URL
	InstagramID       string    `gorm:""`                     // SNS(Instagram)アカウント名
	FacebookID        string    `gorm:""`                     // SNS(Facebook)アカウント名
	PostalCode        string    `gorm:""`                     // 郵便番号
	Prefecture        string    `gorm:"-"`                    // 都道府県
	PrefectureCode    int32     `gorm:"column:prefecture"`    // 都道府県コード
	City              string    `gorm:""`                     // 市区町村
	AddressLine1      string    `gorm:""`                     // 町名・番地
	AddressLine2      string    `gorm:""`                     // ビル名・号室など
	CreatedAt         time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt         time.Time `gorm:""`                     // 更新日時
}

type Coordinators []*Coordinator

type NewCoordinatorParams struct {
	Admin             *Admin
	PhoneNumber       string
	Username          string
	Profile           string
	ThumbnailURL      string
	HeaderURL         string
	PromotionVideoURL string
	BonusVideoURL     string
	InstagramID       string
	FacebookID        string
	PostalCode        string
	PrefectureCode    int32
	City              string
	AddressLine1      string
	AddressLine2      string
}

func NewCoordinator(params *NewCoordinatorParams) (*Coordinator, error) {
	prefecture, err := codes.ToPrefectureJapanese(params.PrefectureCode)
	if err != nil {
		return nil, err
	}
	return &Coordinator{
		AdminID:           params.Admin.ID,
		PhoneNumber:       params.PhoneNumber,
		Username:          params.Username,
		Profile:           params.Profile,
		ThumbnailURL:      params.ThumbnailURL,
		HeaderURL:         params.HeaderURL,
		PromotionVideoURL: params.PromotionVideoURL,
		BonusVideoURL:     params.BonusVideoURL,
		InstagramID:       params.InstagramID,
		FacebookID:        params.FacebookID,
		PostalCode:        params.PostalCode,
		Prefecture:        prefecture,
		PrefectureCode:    params.PrefectureCode,
		City:              params.City,
		AddressLine1:      params.AddressLine1,
		AddressLine2:      params.AddressLine2,
		Admin:             *params.Admin,
	}, nil
}

func (c *Coordinator) Fill(admin *Admin, groups AdminGroupUsers) {
	admin.Fill(groups)
	c.Admin = *admin
	c.Prefecture, _ = codes.ToPrefectureJapanese(c.PrefectureCode)
}

func (cs Coordinators) IDs() []string {
	res := make([]string, len(cs))
	for i := range cs {
		res[i] = cs[i].AdminID
	}
	return res
}

func (cs Coordinators) Fill(admins map[string]*Admin, groups map[string]AdminGroupUsers) {
	for _, c := range cs {
		admin, ok := admins[c.AdminID]
		if !ok {
			admin = &Admin{ID: c.AdminID, Type: AdminTypeCoordinator}
		}
		c.Fill(admin, groups[c.AdminID])
	}
}
