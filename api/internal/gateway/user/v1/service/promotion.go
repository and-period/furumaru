package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
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
			DiscountType: NewDiscountType(promotion.DiscountType).Response(),
			DiscountRate: promotion.DiscountRate,
			Code:         promotion.Code,
			StartAt:      promotion.StartAt.Unix(),
			EndAt:        promotion.EndAt.Unix(),
		},
	}
}

func (p *Promotion) Response() *response.Promotion {
	if p == nil {
		return nil
	}
	return &p.Promotion
}
