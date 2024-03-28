package entity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var errOnlyOneThumbnail = errors.New("entity: only one thumbnail is available")

// ProductStatus - 販売状況
type ProductStatus int32

const (
	ProductStatusUnknown   ProductStatus = 0
	ProductStatusPrivate   ProductStatus = 1 // 非公開
	ProductStatusPresale   ProductStatus = 2 // 予約受付中
	ProductStatusForSale   ProductStatus = 3 // 販売中
	ProductStatusOutOfSale ProductStatus = 4 // 販売期間外
	ProductStatusArchived  ProductStatus = 5 // アーカイブ済み
)

// WeightUnit - 重量単位
type WeightUnit int32

const (
	WeightUnitUnknown  WeightUnit = 0
	WeightUnitGram     WeightUnit = 1 // g(グラム)
	WeightUnitKilogram WeightUnit = 2 // kg(キログラム)
)

// StorageMethodType - 保存方法
type StorageMethodType int32

const (
	StorageMethodTypeUnknown      StorageMethodType = 0
	StorageMethodTypeNormal       StorageMethodType = 1 // 常温保存
	StorageMethodTypeCoolDark     StorageMethodType = 2 // 冷暗所保存
	StorageMethodTypeRefrigerated StorageMethodType = 3 // 冷蔵保存
	StorageMethodTypeFrozen       StorageMethodType = 4 // 冷凍保存
)

// DeliveryType - 配送方法
type DeliveryType int32

const (
	DeliveryTypeUnknown      DeliveryType = 0
	DeliveryTypeNormal       DeliveryType = 1 // 常温便
	DeliveryTypeRefrigerated DeliveryType = 2 // 冷蔵便
	DeliveryTypeFrozen       DeliveryType = 3 // 冷凍便
)

type ProductOrderBy string

const (
	ProductOrderByName             ProductOrderBy = "name"
	ProductOrderBySoldOut          ProductOrderBy = "CASE WHEN (inventory = 0) THEN 1 ELSE 0 END"
	ProductOrderByPublic           ProductOrderBy = "public"
	ProductOrderByInventory        ProductOrderBy = "inventory"
	ProductOrderByOriginPrefecture ProductOrderBy = "origin_prefecture"
	ProductOrderByOriginCity       ProductOrderBy = "origin_city"
	ProductOrderByStartAt          ProductOrderBy = "start_at"
	ProductOrderByCreatedAt        ProductOrderBy = "created_at"
	ProductOrderByUpdatedAt        ProductOrderBy = "updated_at"
)

// Product - 商品情報
type Product struct {
	ProductRevision       `gorm:"-"`
	ID                    string            `gorm:"primaryKey;<-:create"`                   // 商品ID
	CoordinatorID         string            `gorm:""`                                       // コーディネータID
	ProducerID            string            `gorm:""`                                       // 生産者ID
	TypeID                string            `gorm:"column:product_type_id"`                 // 品目ID
	TagIDs                []string          `gorm:"-"`                                      // 商品タグID一覧
	TagIDsJSON            datatypes.JSON    `gorm:"default:null;column:product_tag_ids"`    // 商品タグID一覧(JSON)
	Name                  string            `gorm:""`                                       // 商品名
	Description           string            `gorm:""`                                       // 商品説明
	Public                bool              `gorm:""`                                       // 公開フラグ
	Status                ProductStatus     `gorm:"-"`                                      // 販売状況
	Inventory             int64             `gorm:""`                                       // 在庫数
	Weight                int64             `gorm:""`                                       // 重量
	WeightUnit            WeightUnit        `gorm:""`                                       // 重量単位
	Item                  int64             `gorm:""`                                       // 数量
	ItemUnit              string            `gorm:""`                                       // 数量単位
	ItemDescription       string            `gorm:""`                                       // 数量単位説明
	ThumbnailURL          string            `gorm:"-"`                                      // サムネイルURL
	Thumbnails            common.Images     `gorm:"-"`                                      // リサイズ済サムネイル一覧
	Media                 MultiProductMedia `gorm:"-"`                                      // メディア一覧
	MediaJSON             datatypes.JSON    `gorm:"default:null;column:media"`              // メディア一覧(JSON)
	ExpirationDate        int64             `gorm:""`                                       // 賞味期限(単位:日)
	RecommendedPoints     []string          `gorm:"-"`                                      // おすすめポイント一覧
	RecommendedPointsJSON datatypes.JSON    `gorm:"default:null;column:recommended_points"` // おすすめポイント一覧(JSON)
	StorageMethodType     StorageMethodType `gorm:""`                                       // 保存方法
	DeliveryType          DeliveryType      `gorm:""`                                       // 配送方法
	Box60Rate             int64             `gorm:""`                                       // 箱の占有率(サイズ:60)
	Box80Rate             int64             `gorm:""`                                       // 箱の占有率(サイズ:80)
	Box100Rate            int64             `gorm:""`                                       // 箱の占有率(サイズ:100)
	OriginPrefecture      string            `gorm:"-"`                                      // 原産地(都道府県)
	OriginPrefectureCode  int32             `gorm:"column:origin_prefecture"`               // 原産地(都道府県コード)
	OriginCity            string            `gorm:""`                                       // 原産地(市区町村)
	StartAt               time.Time         `gorm:""`                                       // 販売開始日時
	EndAt                 time.Time         `gorm:""`                                       // 販売終了日時
	CreatedAt             time.Time         `gorm:"<-:create"`                              // 登録日時
	UpdatedAt             time.Time         `gorm:""`                                       // 更新日時
	DeletedAt             gorm.DeletedAt    `gorm:"default:null"`                           // 削除日時
}

type Products []*Product

// ProductMedia - 商品メディア情報
type ProductMedia struct {
	URL         string        `json:"url"`         // メディアURL
	IsThumbnail bool          `json:"isThumbnail"` // サムネイルとして使用
	Images      common.Images `json:"images"`      // リサイズ済み画像一覧
}

type MultiProductMedia []*ProductMedia

type NewProductParams struct {
	CoordinatorID        string
	ProducerID           string
	TypeID               string
	TagIDs               []string
	Name                 string
	Description          string
	Public               bool
	Inventory            int64
	Weight               int64
	WeightUnit           WeightUnit
	Item                 int64
	ItemUnit             string
	ItemDescription      string
	Media                MultiProductMedia
	Price                int64
	Cost                 int64
	ExpirationDate       int64
	RecommendedPoints    []string
	StorageMethodType    StorageMethodType
	DeliveryType         DeliveryType
	Box60Rate            int64
	Box80Rate            int64
	Box100Rate           int64
	OriginPrefectureCode int32
	OriginCity           string
	StartAt              time.Time
	EndAt                time.Time
}

func NewProduct(params *NewProductParams) (*Product, error) {
	productID := uuid.Base58Encode(uuid.New())
	prefecture, err := codes.ToPrefectureJapanese(params.OriginPrefectureCode)
	if err != nil {
		return nil, err
	}
	rparams := &NewProductRevisionParams{
		ProductID: productID,
		Price:     params.Price,
		Cost:      params.Cost,
	}
	revision := NewProductRevision(rparams)
	return &Product{
		ID:                   productID,
		CoordinatorID:        params.CoordinatorID,
		ProducerID:           params.ProducerID,
		TypeID:               params.TypeID,
		TagIDs:               params.TagIDs,
		Name:                 params.Name,
		Description:          params.Description,
		Public:               params.Public,
		Inventory:            params.Inventory,
		Weight:               params.Weight,
		WeightUnit:           params.WeightUnit,
		Item:                 params.Item,
		ItemUnit:             params.ItemUnit,
		ItemDescription:      params.ItemDescription,
		Media:                params.Media,
		ExpirationDate:       params.ExpirationDate,
		RecommendedPoints:    params.RecommendedPoints,
		StorageMethodType:    params.StorageMethodType,
		DeliveryType:         params.DeliveryType,
		Box60Rate:            params.Box60Rate,
		Box80Rate:            params.Box80Rate,
		Box100Rate:           params.Box100Rate,
		OriginPrefecture:     prefecture,
		OriginPrefectureCode: params.OriginPrefectureCode,
		OriginCity:           params.OriginCity,
		StartAt:              params.StartAt,
		EndAt:                params.EndAt,
		ProductRevision:      *revision,
	}, nil
}

func (p *Product) ShippingType() ShippingType {
	switch p.DeliveryType {
	case DeliveryTypeNormal, DeliveryTypeRefrigerated:
		return ShippingTypeNormal
	case DeliveryTypeFrozen:
		return ShippingTypeFrozen
	default:
		return ShippingTypeUnknown
	}
}

func (p *Product) Validate() error {
	if len(p.RecommendedPoints) > 3 {
		return errors.New("entity: limit exceeded recommended points")
	}
	return p.Media.Validate()
}

func (p *Product) Fill(revision *ProductRevision, now time.Time) (err error) {
	p.TagIDs, err = p.unmarshalTagIDs()
	if err != nil {
		return
	}
	p.Media, err = p.unmarshalMedia()
	if err != nil {
		return
	}
	p.RecommendedPoints, err = p.unmarshalRecommendedPoints()
	if err != nil {
		return
	}
	p.SetStatus(now)
	p.SetThumbnail()
	p.ProductRevision = *revision
	p.OriginPrefecture, _ = codes.ToPrefectureJapanese(p.OriginPrefectureCode)
	return
}

func (p *Product) SetStatus(now time.Time) {
	switch {
	case !p.DeletedAt.Time.IsZero():
		p.Status = ProductStatusArchived
	case !p.Public:
		p.Status = ProductStatusPrivate
	case now.Before(p.StartAt):
		p.Status = ProductStatusPresale
	case now.Before(p.EndAt):
		p.Status = ProductStatusForSale
	default:
		p.Status = ProductStatusOutOfSale
	}
}

func (p *Product) SetThumbnail() {
	for _, media := range p.Media {
		if !media.IsThumbnail {
			continue
		}
		p.ThumbnailURL = media.URL
		p.Thumbnails = media.Images
	}
}

func (p *Product) WeightGram() int64 {
	if p.WeightUnit == WeightUnitGram {
		return p.Weight
	}
	return p.Weight * 1e3
}

func (p *Product) unmarshalTagIDs() ([]string, error) {
	if p.TagIDsJSON == nil {
		return []string{}, nil
	}
	var ids []string
	return ids, json.Unmarshal(p.TagIDsJSON, &ids)
}

func (p *Product) unmarshalMedia() (MultiProductMedia, error) {
	if p.MediaJSON == nil {
		return MultiProductMedia{}, nil
	}
	var media MultiProductMedia
	return media, json.Unmarshal(p.MediaJSON, &media)
}

func (p *Product) unmarshalRecommendedPoints() ([]string, error) {
	if p.RecommendedPointsJSON == nil {
		return []string{}, nil
	}
	var points []string
	return points, json.Unmarshal(p.RecommendedPointsJSON, &points)
}

func (p *Product) FillJSON() error {
	media, err := p.Media.Marshal()
	if err != nil {
		return err
	}
	tagIDs, err := ProductMarshalTagIDs(p.TagIDs)
	if err != nil {
		return err
	}
	points, err := ProductMarshalRecommendedPoints(p.RecommendedPoints)
	if err != nil {
		return err
	}
	p.MediaJSON = datatypes.JSON(media)
	p.TagIDsJSON = datatypes.JSON(tagIDs)
	p.RecommendedPointsJSON = datatypes.JSON(points)
	return nil
}

func ProductMarshalTagIDs(tagIDs []string) ([]byte, error) {
	return json.Marshal(tagIDs)
}

func ProductMarshalRecommendedPoints(points []string) ([]byte, error) {
	return json.Marshal(points)
}

func ProductMarshalBusinessDays(days []time.Weekday) ([]byte, error) {
	return json.Marshal(days)
}

func (ps Products) Fill(revisions map[string]*ProductRevision, now time.Time) error {
	for _, p := range ps {
		revision, ok := revisions[p.ID]
		if !ok {
			revision = &ProductRevision{ProductID: p.ID}
		}
		if err := p.Fill(revision, now); err != nil {
			return err
		}
	}
	return nil
}

func (ps Products) Box60Rate() int64 {
	var rate int64
	for i := range ps {
		rate += ps[i].Box60Rate
	}
	return rate
}

func (ps Products) Box80Rate() int64 {
	var rate int64
	for i := range ps {
		rate += ps[i].Box80Rate
	}
	return rate
}

func (ps Products) Box100Rate() int64 {
	var rate int64
	for i := range ps {
		rate += ps[i].Box100Rate
	}
	return rate
}

func (ps Products) WeightGram() int64 {
	var weight int64
	for i := range ps {
		weight += ps[i].WeightGram()
	}
	return weight
}

func (ps Products) IDs() []string {
	return set.UniqBy(ps, func(p *Product) string {
		return p.ID
	})
}

func (ps Products) CoordinatorIDs() []string {
	return set.UniqBy(ps, func(p *Product) string {
		return p.CoordinatorID
	})
}

func (ps Products) ProducerIDs() []string {
	return set.UniqBy(ps, func(p *Product) string {
		return p.ProducerID
	})
}

func (ps Products) ProductTypeIDs() []string {
	return set.UniqBy(ps, func(p *Product) string {
		return p.TypeID
	})
}

func (ps Products) ProductTagIDs() []string {
	res := set.NewEmpty[string](len(ps))
	for i := range ps {
		res.Add(ps[i].TagIDs...)
	}
	return res.Slice()
}

func (ps Products) Map() map[string]*Product {
	res := make(map[string]*Product, len(ps))
	for _, p := range ps {
		res[p.ID] = p
	}
	return res
}

func (ps Products) MapByRevision() map[int64]*Product {
	res := make(map[int64]*Product, len(ps))
	for _, p := range ps {
		res[p.ProductRevision.ID] = p
	}
	return res
}

func (ps Products) Filter(productIDs ...string) Products {
	set := set.New(productIDs...)
	res := make(Products, 0, len(ps))
	for i := range ps {
		if !set.Contains(ps[i].ID) {
			continue
		}
		res = append(res, ps[i])
	}
	return res
}

func (ps Products) FilterBySales() Products {
	res := make(Products, 0, len(ps))
	for _, p := range ps {
		if p.Status != ProductStatusForSale {
			continue
		}
		res = append(res, p)
	}
	return res
}

func (ps Products) FilterByPublished() Products {
	res := make(Products, 0, len(ps))
	for _, p := range ps {
		if p.Status == ProductStatusPrivate || p.Status == ProductStatusArchived {
			continue
		}
		res = append(res, p)
	}
	return res
}

func NewProductMedia(url string, isThumbnail bool) *ProductMedia {
	return &ProductMedia{
		URL:         url,
		IsThumbnail: isThumbnail,
	}
}

func (m *ProductMedia) SetImages(images common.Images) {
	m.Images = images
}

func (m MultiProductMedia) MapByURL() map[string]*ProductMedia {
	res := make(map[string]*ProductMedia, len(m))
	for _, media := range m {
		res[media.URL] = media
	}
	return res
}

func (m MultiProductMedia) Validate() error {
	var exists bool
	for _, media := range m {
		if !media.IsThumbnail {
			continue
		}
		if exists {
			return errOnlyOneThumbnail
		}
		exists = true
	}
	return nil
}

func (m MultiProductMedia) Marshal() ([]byte, error) {
	if len(m) == 0 {
		return []byte{}, nil
	}
	return json.Marshal(m)
}
