package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

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

type Promotion struct {
	types.Promotion
}

type Promotions []*Promotion

func NewPromotionStatus(typ entity.PromotionStatus) PromotionStatus {
	switch typ {
	case entity.PromotionStatusPrivate:
		return PromotionStatusPrivate
	case entity.PromotionStatusWaiting:
		return PromotionStatusWaiting
	case entity.PromotionStatusEnabled:
		return PromotionStatusEnabled
	case entity.PromotionStatusFinished:
		return PromotionStatusFinished
	default:
		return PromotionStatusUnknown
	}
}

func (s PromotionStatus) Response() int32 {
	return int32(s)
}

func NewDiscountType(typ entity.DiscountType) DiscountType {
	switch typ {
	case entity.DiscountTypeAmount:
		return DiscountTypeAmount
	case entity.DiscountTypeRate:
		return DiscountTypeRate
	case entity.DiscountTypeFreeShipping:
		return DiscountTypeFreeShipping
	default:
		return DiscountTypeUnknown
	}
}

func (t DiscountType) StoreEntity() entity.DiscountType {
	switch t {
	case DiscountTypeAmount:
		return entity.DiscountTypeAmount
	case DiscountTypeRate:
		return entity.DiscountTypeRate
	case DiscountTypeFreeShipping:
		return entity.DiscountTypeFreeShipping
	default:
		return entity.DiscountTypeUnknown
	}
}

func (t DiscountType) Response() int32 {
	return int32(t)
}

func NewPromotion(promotion *entity.Promotion) *Promotion {
	return &Promotion{
		Promotion: types.Promotion{
			ID:           promotion.ID,
			Title:        promotion.Title,
			Description:  promotion.Description,
			Status:       NewPromotionStatus(promotion.Status).Response(),
			DiscountType: NewDiscountType(promotion.DiscountType).Response(),
			DiscountRate: promotion.DiscountRate,
			Code:         promotion.Code,
			StartAt:      promotion.StartAt.Unix(),
			EndAt:        promotion.EndAt.Unix(),
		},
	}
}

func (p *Promotion) Response() *types.Promotion {
	if p == nil {
		return nil
	}
	return &p.Promotion
}

func NewPromotions(promotions entity.Promotions) Promotions {
	res := make(Promotions, len(promotions))
	for i := range promotions {
		res[i] = NewPromotion(promotions[i])
	}
	return res
}

func (ps Promotions) Response() []*types.Promotion {
	res := make([]*types.Promotion, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}
