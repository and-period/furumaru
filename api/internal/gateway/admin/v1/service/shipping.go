package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Shipping struct {
	types.Shipping
	ShopID        string
	coordinatorID string
}

type Shippings []*Shipping

type ShippingRate struct {
	types.ShippingRate
}

type ShippingRates []*ShippingRate

func NewShipping(shipping *entity.Shipping) *Shipping {
	return &Shipping{
		Shipping: types.Shipping{
			ID:                shipping.ID,
			Name:              shipping.Name,
			IsDefault:         shipping.IsDefault(),
			Box60Rates:        NewShippingRates(shipping.Box60Rates).Response(),
			Box60Frozen:       shipping.Box60Frozen,
			Box80Rates:        NewShippingRates(shipping.Box80Rates).Response(),
			Box80Frozen:       shipping.Box80Frozen,
			Box100Rates:       NewShippingRates(shipping.Box100Rates).Response(),
			Box100Frozen:      shipping.Box100Frozen,
			HasFreeShipping:   shipping.HasFreeShipping,
			FreeShippingRates: shipping.FreeShippingRates,
			CreatedAt:         shipping.CreatedAt.Unix(),
			UpdatedAt:         shipping.CreatedAt.Unix(),
		},
		ShopID:        shipping.ShopID,
		coordinatorID: shipping.CoordinatorID,
	}
}

func (s *Shipping) Response() *types.Shipping {
	return &s.Shipping
}

func NewShippings(shippings entity.Shippings) Shippings {
	res := make(Shippings, len(shippings))
	for i := range shippings {
		res[i] = NewShipping(shippings[i])
	}
	return res
}

func (ss Shippings) Response() []*types.Shipping {
	res := make([]*types.Shipping, len(ss))
	for i := range ss {
		res[i] = ss[i].Response()
	}
	return res
}

func NewShippingRate(rate *entity.ShippingRate) *ShippingRate {
	return &ShippingRate{
		ShippingRate: types.ShippingRate{
			Number:          rate.Number,
			Name:            rate.Name,
			Price:           rate.Price,
			PrefectureCodes: rate.PrefectureCodes,
		},
	}
}

func (r *ShippingRate) Response() *types.ShippingRate {
	return &r.ShippingRate
}

func NewShippingRates(rates entity.ShippingRates) ShippingRates {
	res := make(ShippingRates, len(rates))
	for i := range rates {
		res[i] = NewShippingRate(rates[i])
	}
	return res
}

func (rs ShippingRates) Response() []*types.ShippingRate {
	res := make([]*types.ShippingRate, len(rs))
	for i := range rs {
		res[i] = rs[i].Response()
	}
	return res
}
