package entity

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// WeightUnit - 重量単位
type WeightUnit int32

const (
	WeightUnitUnknown  WeightUnit = 0
	WeightUnitGram     WeightUnit = 1 // g(グラム)
	WeightUnitKilogram WeightUnit = 2 // kg(キログラム)
)

// DeliveryType - 配送方法
type DeliveryType int32

const (
	DeliveryTypeUnknown      DeliveryType = 0
	DeliveryTypeNormal       DeliveryType = 1 // 通常便
	DeliveryTypeRefrigerated DeliveryType = 2 // 冷蔵便
	DeliveryTypeFrozen       DeliveryType = 3 // 冷凍便
)

// Product - 商品情報
type Product struct {
	ID               string            `gorm:"primaryKey;<-:create"`                // 商品ID
	ProducerID       string            `gorm:""`                                    // 生産者ID
	CategoryID       string            `gorm:"default:null"`                        // 商品種別ID
	TypeID           string            `gorm:"default:null;column:product_type_id"` // 品目ID
	Name             string            `gorm:""`                                    // 商品名
	Description      string            `gorm:""`                                    // 商品説明
	Public           bool              `gorm:""`                                    // 公開フラグ
	Inventory        int64             `gorm:""`                                    // 在庫数
	Weight           int64             `gorm:""`                                    // 重量
	WeightUnit       WeightUnit        `gorm:""`                                    // 重量単位
	Item             int64             `gorm:""`                                    // 数量
	ItemUnit         string            `gorm:""`                                    // 数量単位
	ItemDescription  string            `gorm:""`                                    // 数量単位説明
	Media            MultiProductMedia `gorm:"-"`                                   // メディア一覧
	MediaJSON        datatypes.JSON    `gorm:"default:null;column:media"`           // メディア一覧(JSON)
	Price            int64             `gorm:""`                                    // 販売価格
	DeliveryType     DeliveryType      `gorm:""`                                    // 配送方法
	Box60Rate        int64             `gorm:""`                                    // 箱の占有率(サイズ:60)
	Box80Rate        int64             `gorm:""`                                    // 箱の占有率(サイズ:80)
	Box100Rate       int64             `gorm:""`                                    // 箱の占有率(サイズ:100)
	OriginPrefecture string            `gorm:""`                                    // 原産地(都道府県)
	OriginCity       string            `gorm:""`                                    // 原産地(市区町村)
	CreatedAt        time.Time         `gorm:"<-:create"`                           // 登録日時
	CreatedBy        string            `gorm:"<-:create"`                           // 登録者ID
	UpdatedAt        time.Time         `gorm:""`                                    // 更新日時
	UpdatedBy        string            `gorm:""`                                    // 更新者ID
	DeletedAt        gorm.DeletedAt    `gorm:"default:null"`                        // 削除日時
}

type Products []*Product

// ProductMedia - 商品メディア情報
type ProductMedia struct {
	URL         string `json:"url"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail"` // サムネイルとして使用
}

type MultiProductMedia []*ProductMedia

func (p *Product) Fill() error {
	var media MultiProductMedia
	if err := json.Unmarshal(p.MediaJSON, &media); err != nil {
		return err
	}
	p.Media = media
	return nil
}

func (p *Product) FillJSON() error {
	v, err := json.Marshal(p.Media)
	if err != nil {
		return err
	}
	p.MediaJSON = datatypes.JSON(v)
	return nil
}
