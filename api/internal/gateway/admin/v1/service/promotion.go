package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/set"
)

// PromotionStatus - プロモーションの状態
type PromotionStatus types.PromotionStatus

// DiscountType - 割引計算方法
type DiscountType types.DiscountType

// PromotionTargetType - プロモーションの対象
type PromotionTargetType types.PromotionTargetType

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

func NewPromotionTargetType(typ entity.PromotionTargetType) PromotionTargetType {
	switch typ {
	case entity.PromotionTargetTypeAllShop:
		return PromotionTargetType(types.PromotionTargetTypeAllShop)
	case entity.PromotionTargetTypeSpecificShop:
		return PromotionTargetType(types.PromotionTargetTypeSpecificShop)
	default:
		return PromotionTargetType(types.PromotionTargetTypeUnknown)
	}
}

func (t PromotionTargetType) Response() types.PromotionTargetType {
	return types.PromotionTargetType(t)
}

func NewPromotion(promotion *entity.Promotion, aggregate *entity.AggregatedOrderPromotion) *Promotion {
	var usedCount, usedAmount int64
	if aggregate != nil {
		usedCount = aggregate.OrderCount
		usedAmount = aggregate.DiscountTotal
	}
	return &Promotion{
		Promotion: types.Promotion{
			ID:           promotion.ID,
			ShopID:       promotion.ShopID,
			Title:        promotion.Title,
			Description:  promotion.Description,
			Status:       NewPromotionStatus(promotion.Status).Response(),
			Public:       promotion.Public,
			TargetType:   NewPromotionTargetType(promotion.TargetType).Response(),
			DiscountType: NewDiscountType(promotion.DiscountType).Response(),
			DiscountRate: promotion.DiscountRate,
			Code:         promotion.Code,
			UsedCount:    usedCount,
			UsedAmount:   usedAmount,
			StartAt:      promotion.StartAt.Unix(),
			EndAt:        promotion.EndAt.Unix(),
			CreatedAt:    promotion.CreatedAt.Unix(),
			UpdatedAt:    promotion.UpdatedAt.Unix(),
		},
	}
}

func (p *Promotion) Response() *types.Promotion {
	if p == nil {
		return nil
	}
	return &p.Promotion
}

func NewPromotions(promotions entity.Promotions, aggregates map[string]*entity.AggregatedOrderPromotion) Promotions {
	res := make(Promotions, len(promotions))
	for i, p := range promotions {
		res[i] = NewPromotion(promotions[i], aggregates[p.ID])
	}
	return res
}

func (ps Promotions) ShopIDs() []string {
	return set.UniqBy(ps, func(p *Promotion) string {
		return p.ShopID
	})
}

func (ps Promotions) Map() map[string]*Promotion {
	res := make(map[string]*Promotion, len(ps))
	for _, p := range ps {
		res[p.ID] = p
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
