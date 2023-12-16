package entity

import (
	"encoding/json"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/pkg/set"
	"gorm.io/datatypes"
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
	Prefecture         string         `gorm:"-"`                                    // 都道府県
	PrefectureCode     int32          `gorm:"column:prefecture"`                    // 都道府県コード
	City               string         `gorm:""`                                     // 市区町村
	AddressLine1       string         `gorm:""`                                     // 町名・番地
	AddressLine2       string         `gorm:""`                                     // ビル名・号室など
	BusinessDays       []time.Weekday `gorm:"-"`                                    // 営業曜日(発送可能日)一覧
	BusinessDaysJSON   datatypes.JSON `gorm:"default:null;column:business_days"`    // 営業曜日(発送可能日)一覧(JSON)
	CreatedAt          time.Time      `gorm:"<-:create"`                            // 登録日時
	UpdatedAt          time.Time      `gorm:""`                                     // 更新日時
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
	PrefectureCode    int32
	City              string
	AddressLine1      string
	AddressLine2      string
	BusinessDays      []time.Weekday
}

func NewCoordinator(params *NewCoordinatorParams) (*Coordinator, error) {
	prefecture, err := codes.ToPrefectureJapanese(params.PrefectureCode)
	if err != nil {
		return nil, err
	}
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
		Prefecture:        prefecture,
		PrefectureCode:    params.PrefectureCode,
		City:              params.City,
		AddressLine1:      params.AddressLine1,
		AddressLine2:      params.AddressLine2,
		BusinessDays:      params.BusinessDays,
		Admin:             *params.Admin,
	}, nil
}

func (c *Coordinator) Fill(admin *Admin) (err error) {
	c.Thumbnails, err = common.NewImagesFromBytes(c.ThumbnailsJSON)
	if err != nil {
		return err
	}
	c.Headers, err = common.NewImagesFromBytes(c.HeadersJSON)
	if err != nil {
		return err
	}
	c.ProductTypeIDs, err = c.unmarshalProductTypeIDs()
	if err != nil {
		return err
	}
	c.BusinessDays, err = c.unmarshalBusinessDays()
	if err != nil {
		return err
	}
	c.Admin = *admin
	c.Prefecture, _ = codes.ToPrefectureJapanese(c.PrefectureCode)
	return nil
}

func (c *Coordinator) unmarshalProductTypeIDs() ([]string, error) {
	if c.ProductTypeIDsJSON == nil {
		return []string{}, nil
	}
	var ids []string
	return ids, json.Unmarshal(c.ProductTypeIDsJSON, &ids)
}

func (c *Coordinator) unmarshalBusinessDays() ([]time.Weekday, error) {
	if c.BusinessDaysJSON == nil {
		return []time.Weekday{}, nil
	}
	var days []time.Weekday
	return days, json.Unmarshal(c.BusinessDaysJSON, &days)
}

func (c *Coordinator) FillJSON() error {
	tagsIDs, err := CoordinatorMarshalProductTypeIDs(c.ProductTypeIDs)
	if err != nil {
		return err
	}
	days, err := CoordinatorMarshalBusinessDays(c.BusinessDays)
	if err != nil {
		return err
	}
	c.ProductTypeIDsJSON = datatypes.JSON(tagsIDs)
	c.BusinessDaysJSON = datatypes.JSON(days)
	return nil
}

func CoordinatorMarshalProductTypeIDs(types []string) ([]byte, error) {
	return json.Marshal(types)
}

func CoordinatorMarshalBusinessDays(days []time.Weekday) ([]byte, error) {
	return json.Marshal(days)
}

func (cs Coordinators) IDs() []string {
	res := make([]string, len(cs))
	for i := range cs {
		res[i] = cs[i].AdminID
	}
	return res
}

func (cs Coordinators) ProductTypeIDs() []string {
	res := set.NewEmpty[string](len(cs))
	for i := range cs {
		res.Add(cs[i].ProductTypeIDs...)
	}
	return res.Slice()
}

func (cs Coordinators) Fill(admins map[string]*Admin) error {
	for _, c := range cs {
		admin, ok := admins[c.AdminID]
		if !ok {
			admin = &Admin{ID: c.AdminID, Role: AdminRoleCoordinator}
		}
		if err := c.Fill(admin); err != nil {
			return err
		}
	}
	return nil
}
