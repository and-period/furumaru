package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
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
	response.Promotion
}

type Promotions []*Promotion

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
		Promotion: response.Promotion{
			ID:           promotion.ID,
			Title:        promotion.Title,
			Description:  promotion.Description,
			Public:       promotion.Public,
			PublishedAt:  promotion.PublishedAt.Unix(),
			DiscountType: NewDiscountType(promotion.DiscountType).Response(),
			DiscountRate: promotion.DiscountRate,
			Code:         promotion.Code,
			StartAt:      promotion.StartAt.Unix(),
			EndAt:        promotion.EndAt.Unix(),
			CreatedAt:    promotion.CreatedAt.Unix(),
			UpdatedAt:    promotion.UpdatedAt.Unix(),
		},
	}
}

func (p *Promotion) Response() *response.Promotion {
	return &p.Promotion
}

func NewPromotions(promotions entity.Promotions) Promotions {
	res := make(Promotions, len(promotions))
	for i := range promotions {
		res[i] = NewPromotion(promotions[i])
	}
	return res
}

func (ps Promotions) Map() map[string]*Promotion {
	res := make(map[string]*Promotion, len(ps))
	for _, p := range ps {
		res[p.ID] = p
	}
	return res
}

func (ps Promotions) Response() []*response.Promotion {
	res := make([]*response.Promotion, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}
