package entity

import (
	"errors"
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"github.com/shopspring/decimal"
)

var errInvalidDiscount = errors.New("entity: invalid discount value")

// PromotionStatus - プロモーションの状態
type PromotionStatus int32

const (
	PromotionStatusUnknown  PromotionStatus = 0
	PromotionStatusPrivate  PromotionStatus = 1 // 非公開
	PromotionStatusWaiting  PromotionStatus = 2 // 利用開始前
	PromotionStatusEnabled  PromotionStatus = 3 // 利用可能
	PromotionStatusFinished PromotionStatus = 4 // 利用終了
)

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

// PromotionTargetType - プロモーション対象種別
type PromotionTargetType int32

const (
	PromotionTargetTypeAllShop      PromotionTargetType = 0 // 全ての店舗
	PromotionTargetTypeSpecificShop PromotionTargetType = 1 // 特定の店舗のみ
)

// Promotion - プロモーション情報
type Promotion struct {
	ID           string              `gorm:"primaryKey;<-:create"` // プロモーションID
	ShopID       string              `gorm:"default:null"`         // ショップID
	Status       PromotionStatus     `gorm:"-"`                    // 状態
	Title        string              `gorm:""`                     // タイトル
	Description  string              `gorm:""`                     // 詳細説明
	Public       bool                `gorm:""`                     // Deprecated: 公開フラグ
	TargetType   PromotionTargetType `gorm:""`                     // 対象種別
	DiscountType DiscountType        `gorm:""`                     // 割引計算方法
	DiscountRate int64               `gorm:""`                     // 割引額(%/円)
	Code         string              `gorm:"<-:create"`            // クーポンコード
	CodeType     PromotionCodeType   `gorm:"<-:create"`            // クーポンコード種別
	StartAt      time.Time           `gorm:""`                     // クーポン使用可能日時(開始)
	EndAt        time.Time           `gorm:""`                     // クーポン使用可能日時(終了)
	CreatedAt    time.Time           `gorm:"<-:create"`            // 登録日時
	UpdatedAt    time.Time           `gorm:""`                     // 更新日時
}

type Promotions []*Promotion

type NewPromotionParams struct {
	ShopID       string
	Title        string
	Description  string
	Public       bool
	DiscountType DiscountType
	DiscountRate int64
	Code         string
	CodeType     PromotionCodeType
	StartAt      time.Time
	EndAt        time.Time
}

func NewPromotion(params *NewPromotionParams) *Promotion {
	targetType := PromotionTargetTypeAllShop
	if params.ShopID != "" {
		targetType = PromotionTargetTypeSpecificShop
	}
	return &Promotion{
		ID:           uuid.Base58Encode(uuid.New()),
		ShopID:       params.ShopID,
		Title:        params.Title,
		Description:  params.Description,
		Public:       params.Public,
		TargetType:   targetType,
		DiscountType: params.DiscountType,
		DiscountRate: params.DiscountRate,
		Code:         params.Code,
		CodeType:     params.CodeType,
		StartAt:      params.StartAt,
		EndAt:        params.EndAt,
	}
}

func (p *Promotion) CalcDiscount(total int64, shippingFee int64) int64 {
	if p == nil {
		return 0
	}
	switch p.DiscountType {
	case DiscountTypeAmount:
		if total < p.DiscountRate {
			return total
		}
		return p.DiscountRate
	case DiscountTypeRate:
		if p.DiscountRate == 0 {
			return 0
		}
		dtotal := decimal.NewFromInt(total)
		rate := decimal.NewFromInt(p.DiscountRate).Div(decimal.NewFromInt(100))
		return dtotal.Mul(rate).IntPart()
	case DiscountTypeFreeShipping:
		return shippingFee
	default:
		return 0
	}
}

func (p *Promotion) IsEnabled(shopID string) bool {
	if p == nil {
		return false
	}
	if p.Status != PromotionStatusEnabled {
		return false
	}
	switch p.TargetType {
	case PromotionTargetTypeAllShop:
		return true
	case PromotionTargetTypeSpecificShop:
		return p.ShopID == shopID
	default:
		return false
	}
}

func (p *Promotion) Fill(now time.Time) {
	p.SetStatus(now)
}

func (p *Promotion) SetStatus(now time.Time) {
	switch {
	case !p.Public:
		p.Status = PromotionStatusPrivate
	case now.Before(p.StartAt):
		p.Status = PromotionStatusWaiting
	case now.Before(p.EndAt):
		p.Status = PromotionStatusEnabled
	default:
		p.Status = PromotionStatusFinished
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

func (ps Promotions) IDs() []string {
	return set.UniqBy(ps, func(p *Promotion) string {
		return p.ID
	})
}

func (ps Promotions) ShopIDs() []string {
	return set.UniqBy(ps, func(p *Promotion) string {
		return p.ShopID
	})
}

func (ps Promotions) Fill(now time.Time) {
	for i := range ps {
		ps[i].Fill(now)
	}
}
