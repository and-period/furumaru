package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

// PromotionStatus - プロモーションの状態
type PromotionStatus types.PromotionStatus

// DiscountType - 割引計算方法
type DiscountType types.DiscountType

type Promotion struct {
	types.Promotion
}

type Promotions []*Promotion

func NewPromotionStatus(typ entity.PromotionStatus) PromotionStatus {
	switch typ {
	case entity.PromotionStatusPrivate:
		return PromotionStatus(types.PromotionStatusPrivate)
	case entity.PromotionStatusWaiting:
		return PromotionStatus(types.PromotionStatusWaiting)
	case entity.PromotionStatusEnabled:
		return PromotionStatus(types.PromotionStatusEnabled)
	case entity.PromotionStatusFinished:
		return PromotionStatus(types.PromotionStatusFinished)
	default:
		return PromotionStatus(types.PromotionStatusUnknown)
	}
}

func (s PromotionStatus) Response() types.PromotionStatus {
	return types.PromotionStatus(s)
}

func NewDiscountType(typ entity.DiscountType) DiscountType {
	switch typ {
	case entity.DiscountTypeAmount:
		return DiscountType(types.DiscountTypeAmount)
	case entity.DiscountTypeRate:
		return DiscountType(types.DiscountTypeRate)
	case entity.DiscountTypeFreeShipping:
		return DiscountType(types.DiscountTypeFreeShipping)
	default:
		return DiscountType(types.DiscountTypeUnknown)
	}
}

func (t DiscountType) StoreEntity() entity.DiscountType {
	switch types.DiscountType(t) {
	case types.DiscountTypeAmount:
		return entity.DiscountTypeAmount
	case types.DiscountTypeRate:
		return entity.DiscountTypeRate
	case types.DiscountTypeFreeShipping:
		return entity.DiscountTypeFreeShipping
	default:
		return entity.DiscountTypeUnknown
	}
}

func (t DiscountType) Response() types.DiscountType {
	return types.DiscountType(t)
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
