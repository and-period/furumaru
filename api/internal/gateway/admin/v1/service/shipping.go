package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Shipping struct {
	response.Shipping
}

type Shippings []*Shipping

type ShippingRate struct {
	response.ShippingRate
}

type ShippingRates []*ShippingRate

func NewShipping(shipping *entity.Shipping) *Shipping {
	return &Shipping{
		Shipping: response.Shipping{
			ID:                 shipping.ID,
			CoordinatorID:      shipping.CoordinatorID,
			Name:               shipping.Name,
			IsDefault:          shipping.IsDefault,
			Box60Rates:         NewShippingRates(shipping.Box60Rates).Response(),
			Box60Refrigerated:  shipping.Box60Refrigerated,
			Box60Frozen:        shipping.Box60Frozen,
			Box80Rates:         NewShippingRates(shipping.Box80Rates).Response(),
			Box80Refrigerated:  shipping.Box80Refrigerated,
			Box80Frozen:        shipping.Box80Frozen,
			Box100Rates:        NewShippingRates(shipping.Box100Rates).Response(),
			Box100Refrigerated: shipping.Box100Refrigerated,
			Box100Frozen:       shipping.Box100Frozen,
			HasFreeShipping:    shipping.HasFreeShipping,
			FreeShippingRates:  shipping.FreeShippingRates,
			CreatedAt:          shipping.CreatedAt.Unix(),
			UpdatedAt:          shipping.CreatedAt.Unix(),
		},
	}
}

func (s *Shipping) Response() *response.Shipping {
	return &s.Shipping
}

func NewShippings(shippings entity.Shippings) Shippings {
	res := make(Shippings, len(shippings))
	for i := range shippings {
		res[i] = NewShipping(shippings[i])
	}
	return res
}

func (ss Shippings) Response() []*response.Shipping {
	res := make([]*response.Shipping, len(ss))
	for i := range ss {
		res[i] = ss[i].Response()
	}
	return res
}

func NewShippingRate(rate *entity.ShippingRate) *ShippingRate {
	return &ShippingRate{
		ShippingRate: response.ShippingRate{
			Number:          rate.Number,
			Name:            rate.Name,
			Price:           rate.Price,
			PrefectureCodes: rate.PrefectureCodes,
		},
	}
}

func (r *ShippingRate) Response() *response.ShippingRate {
	return &r.ShippingRate
}

func (ss Shippings) Map() map[string]*Shipping {
	res := make(map[string]*Shipping, len(ss))
	for _, s := range ss {
		res[s.ID] = s
	}
	return res
}

func NewShippingRates(rates entity.ShippingRates) ShippingRates {
	res := make(ShippingRates, len(rates))
	for i := range rates {
		res[i] = NewShippingRate(rates[i])
	}
	return res
}

func (rs ShippingRates) Response() []*response.ShippingRate {
	res := make([]*response.ShippingRate, len(rs))
	for i := range rs {
		res[i] = rs[i].Response()
	}
	return res
}
