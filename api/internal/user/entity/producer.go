package entity

import (
	"encoding/json"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Producer - 生産者情報
type Producer struct {
	Admin          `gorm:"-"`
	AdminID        string         `gorm:"primaryKey;<-:create"`           // 管理者ID
	CoordinatorID  string         `gorm:"default:null"`                   // コーディネータID
	PhoneNumber    string         `gorm:""`                               // 電話番号
	StoreName      string         `gorm:""`                               // 店舗名
	ThumbnailURL   string         `gorm:""`                               // サムネイルURL
	Thumbnails     common.Images  `gorm:"-"`                              // サムネイル一覧(リサイズ済み)
	ThumbnailsJSON datatypes.JSON `gorm:"default:null;column:thumbnails"` // サムネイル一覧(JSON)
	HeaderURL      string         `gorm:""`                               // ヘッダー画像URL
	Headers        common.Images  `gorm:"-"`                              // ヘッダー画像一覧(リサイズ済み)
	HeadersJSON    datatypes.JSON `gorm:"default:null;column:headers"`    // ヘッダー画像一覧(JSON)
	PostalCode     string         `gorm:""`                               // 郵便番号
	Prefecture     string         `gorm:""`                               // 都道府県
	City           string         `gorm:""`                               // 市区町村
	AddressLine1   string         `gorm:""`                               // 町名・番地
	AddressLine2   string         `gorm:""`                               // ビル名・号室など
	CreatedAt      time.Time      `gorm:"<-:create"`                      // 登録日時
	UpdatedAt      time.Time      `gorm:""`                               // 更新日時
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`                   // 退会日時
}

type Producers []*Producer

type NewProducerParams struct {
	Admin        *Admin
	PhoneNumber  string
	StoreName    string
	ThumbnailURL string
	HeaderURL    string
	PostalCode   string
	Prefecture   string
	City         string
	AddressLine1 string
	AddressLine2 string
}

func NewProducer(params *NewProducerParams) *Producer {
	return &Producer{
		AdminID:      params.Admin.ID,
		PhoneNumber:  params.PhoneNumber,
		StoreName:    params.StoreName,
		ThumbnailURL: params.ThumbnailURL,
		HeaderURL:    params.HeaderURL,
		PostalCode:   params.PostalCode,
		Prefecture:   params.Prefecture,
		City:         params.City,
		AddressLine1: params.AddressLine1,
		AddressLine2: params.AddressLine2,
		Admin:        *params.Admin,
	}
}

func (p *Producer) Fill(admin *Admin) (err error) {
	var thumbnails, headers common.Images
	if thumbnails, err = common.NewImagesFromBytes(p.ThumbnailsJSON); err != nil {
		return err
	}
	if headers, err = common.NewImagesFromBytes(p.HeadersJSON); err != nil {
		return err
	}
	p.Admin = *admin
	p.Thumbnails = thumbnails
	p.Headers = headers
	return nil
}

func (p *Producer) FillJSON() error {
	thumbnails, err := p.marshalThumbnails()
	if err != nil {
		return err
	}
	headers, err := p.marshalHeaders()
	if err != nil {
		return err
	}
	p.ThumbnailsJSON = datatypes.JSON(thumbnails)
	p.HeadersJSON = datatypes.JSON(headers)
	return nil
}

func (p *Producer) marshalThumbnails() ([]byte, error) {
	if len(p.Thumbnails) == 0 {
		return []byte{}, nil
	}
	return json.Marshal(p.Thumbnails)
}

func (p *Producer) marshalHeaders() ([]byte, error) {
	if len(p.Headers) == 0 {
		return []byte{}, nil
	}
	return json.Marshal(p.Headers)
}

func (ps Producers) IDs() []string {
	res := make([]string, len(ps))
	for i := range ps {
		res[i] = ps[i].AdminID
	}
	return res
}

func (ps Producers) Related() Producers {
	res := make(Producers, 0, len(ps))
	for _, p := range ps {
		if p.CoordinatorID == "" {
			continue
		}
		res = append(res, p)
	}
	return res
}

func (ps Producers) Unrelated() Producers {
	res := make(Producers, 0, len(ps))
	for _, p := range ps {
		if p.CoordinatorID != "" {
			continue
		}
		res = append(res, p)
	}
	return res
}
