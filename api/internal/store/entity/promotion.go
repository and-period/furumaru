package entity

import (
	"errors"
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

var errInvalidDiscount = errors.New("entity: invalid discount value")

// DiscountType - 割引計算方法
type DiscountType int32

const (
	DiscountTypeUnknown      DiscountType = 0
	DiscountTypeAmount       DiscountType = 1 // 固定額(円)
	DiscountTypeRate         DiscountType = 2 // 料率計算(%)
	DiscountTypeFreeShipping DiscountType = 3 // 送料無料
)

// PromotionCodeType - プロモーションコード種別
type PromotionCodeType int32

const (
	PromotionCodeTypeUnknown PromotionCodeType = 0
	PromotionCodeTypeOnce    PromotionCodeType = 1 // １回限り利用可能
	PromotionCodeTypeAlways  PromotionCodeType = 2 // 期間内回数無制限
)

type PromotionOrderBy string

const (
	PromotionOrderByTitle       PromotionOrderBy = "title"
	PromotionOrderByPublic      PromotionOrderBy = "public"
	PromotionOrderByPublishedAt PromotionOrderBy = "published_at"
	PromotionOrderByStartAt     PromotionOrderBy = "start_at"
	PromotionOrderByEndAt       PromotionOrderBy = "end_at"
	PromotionOrderByCreatedAt   PromotionOrderBy = "created_at"
	PromotionOrderByUpdatedAt   PromotionOrderBy = "updated_at"
)

// Promotion - プロモーション情報
type Promotion struct {
	ID           string            `gorm:"primaryKey;<-:create"` // プロモーションID
	Title        string            `gorm:""`                     // タイトル
	Description  string            `gorm:""`                     // 詳細説明
	Public       bool              `gorm:""`                     // 公開フラグ
	PublishedAt  time.Time         `gorm:""`                     // 公開日時
	DiscountType DiscountType      `gorm:""`                     // 割引計算方法
	DiscountRate int64             `gorm:""`                     // 割引額(%/円)
	Code         string            `gorm:"<-:create"`            // クーポンコード
	CodeType     PromotionCodeType `gorm:"<-:create"`            // クーポンコード種別
	StartAt      time.Time         `gorm:""`                     // クーポン使用可能日時(開始)
	EndAt        time.Time         `gorm:""`                     // クーポン使用可能日時(終了)
	CreatedAt    time.Time         `gorm:"<-:create"`            // 登録日時
	UpdatedAt    time.Time         `gorm:""`                     // 更新日時
}

type Promotions []*Promotion

type NewPromotionParams struct {
	Title        string
	Description  string
	Public       bool
	PublishedAt  time.Time
	DiscountType DiscountType
	DiscountRate int64
	Code         string
	CodeType     PromotionCodeType
	StartAt      time.Time
	EndAt        time.Time
}

func NewPromotion(params *NewPromotionParams) *Promotion {
	return &Promotion{
		ID:           uuid.Base58Encode(uuid.New()),
		Title:        params.Title,
		Description:  params.Description,
		Public:       params.Public,
		PublishedAt:  params.PublishedAt,
		DiscountType: params.DiscountType,
		DiscountRate: params.DiscountRate,
		Code:         params.Code,
		CodeType:     params.CodeType,
		StartAt:      params.StartAt,
		EndAt:        params.EndAt,
	}
}

func (p *Promotion) Validate() error {
	switch p.DiscountType {
	case DiscountTypeAmount:
		if p.DiscountRate <= 0 {
			return errInvalidDiscount
		}
	case DiscountTypeRate:
		if p.DiscountRate <= 0 || p.DiscountRate > 100 {
			return errInvalidDiscount
		}
	case DiscountTypeFreeShipping:
		p.DiscountRate = 0
	}
	return nil
}
