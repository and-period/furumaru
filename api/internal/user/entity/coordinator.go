package entity

import (
	"encoding/json"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Coordinator - コーディネータ情報
type Coordinator struct {
	Admin              `gorm:"-"`
	AdminID            string         `gorm:"primaryKey;<-:create"`                 // 管理者ID
	PhoneNumber        string         `gorm:""`                                     // 電話番号
	MarcheName         string         `gorm:""`                                     // マルシェ名
	Username           string         `gorm:""`                                     // 表示名
	Profile            string         `gorm:""`                                     // 紹介文
	ProductTypeIDs     []string       `gorm:"-"`                                    // 取り扱い品目ID一覧
	ProductTypeIDsJSON datatypes.JSON `gorm:"default:null;column:product_type_ids"` // 取り扱い品目ID一覧(JSON)
	ThumbnailURL       string         `gorm:""`                                     // サムネイルURL
	Thumbnails         common.Images  `gorm:"-"`                                    // サムネイル一覧(リサイズ済み)
	ThumbnailsJSON     datatypes.JSON `gorm:"default:null;column:thumbnails"`       // サムネイル一覧(JSON)
	HeaderURL          string         `gorm:""`                                     // ヘッダー画像URL
	Headers            common.Images  `gorm:"-"`                                    // ヘッダー画像一覧(リサイズ済み)
	HeadersJSON        datatypes.JSON `gorm:"default:null;column:headers"`          // ヘッダー画像一覧(JSON)
	PromotionVideoURL  string         `gorm:""`                                     // 紹介動画URL
	BonusVideoURL      string         `gorm:""`                                     // 購入特典動画URL
	InstagramID        string         `gorm:""`                                     // SNS(Instagram)アカウント名
	FacebookID         string         `gorm:""`                                     // SNS(Facebook)アカウント名
	PostalCode         string         `gorm:""`                                     // 郵便番号
	Prefecture         string         `gorm:""`                                     // 都道府県
	City               string         `gorm:""`                                     // 市区町村
	AddressLine1       string         `gorm:""`                                     // 町名・番地
	AddressLine2       string         `gorm:""`                                     // ビル名・号室など
	CreatedAt          time.Time      `gorm:"<-:create"`                            // 登録日時
	UpdatedAt          time.Time      `gorm:""`                                     // 更新日時
	DeletedAt          gorm.DeletedAt `gorm:"default:null"`                         // 退会日時
}

type Coordinators []*Coordinator

type NewCoordinatorParams struct {
	Admin             *Admin
	PhoneNumber       string
	MarcheName        string
	Username          string
	Profile           string
	ProductTypeIDs    []string
	ThumbnailURL      string
	HeaderURL         string
	PromotionVideoURL string
	BonusVideoURL     string
	InstagramID       string
	FacebookID        string
	PostalCode        string
	Prefecture        string
	City              string
	AddressLine1      string
	AddressLine2      string
}

func NewCoordinator(params *NewCoordinatorParams) *Coordinator {
	return &Coordinator{
		AdminID:           params.Admin.ID,
		PhoneNumber:       params.PhoneNumber,
		MarcheName:        params.MarcheName,
		Username:          params.Username,
		Profile:           params.Profile,
		ProductTypeIDs:    params.ProductTypeIDs,
		ThumbnailURL:      params.ThumbnailURL,
		HeaderURL:         params.HeaderURL,
		PromotionVideoURL: params.PromotionVideoURL,
		BonusVideoURL:     params.BonusVideoURL,
		InstagramID:       params.InstagramID,
		FacebookID:        params.FacebookID,
		PostalCode:        params.PostalCode,
		Prefecture:        params.Prefecture,
		City:              params.City,
		AddressLine1:      params.AddressLine1,
		AddressLine2:      params.AddressLine2,
		Admin:             *params.Admin,
	}
}

func (c *Coordinator) Fill(admin *Admin) (err error) {
	var (
		thumbnails, headers common.Images
		productTypeIDs      []string
	)
	if thumbnails, err = common.NewImagesFromBytes(c.ThumbnailsJSON); err != nil {
		return err
	}
	if headers, err = common.NewImagesFromBytes(c.HeadersJSON); err != nil {
		return err
	}
	if productTypeIDs, err = c.unmarshalProductTypeIDs(c.ProductTypeIDsJSON); err != nil {
		return err
	}
	c.Admin = *admin
	c.Thumbnails = thumbnails
	c.Headers = headers
	c.ProductTypeIDs = productTypeIDs
	return nil
}

func (c *Coordinator) unmarshalProductTypeIDs(b []byte) ([]string, error) {
	if b == nil {
		return []string{}, nil
	}
	var ids []string
	return ids, json.Unmarshal(b, &ids)
}

func (c *Coordinator) FillJSON() error {
	v, err := CoordinatorMarshalProductTypeIDs(c.ProductTypeIDs)
	if err != nil {
		return err
	}
	c.ProductTypeIDsJSON = datatypes.JSON(v)
	return nil
}

func CoordinatorMarshalProductTypeIDs(types []string) ([]byte, error) {
	return json.Marshal(types)
}

func (cs Coordinators) IDs() []string {
	res := make([]string, len(cs))
	for i := range cs {
		res[i] = cs[i].AdminID
	}
	return res
}
