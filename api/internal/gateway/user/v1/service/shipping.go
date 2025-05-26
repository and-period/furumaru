package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Shipping struct {
	response.Shipping
}

func NewShipping(shipping *entity.Shipping) *Shipping {
	return &Shipping{
		Shipping: response.Shipping{
			ID:                shipping.ID,
			Box60Rates:        NewShippingRates(shipping.Box60Rates).Response(),
			Box60Frozen:       shipping.Box60Frozen,
			Box80Rates:        NewShippingRates(shipping.Box80Rates).Response(),
			Box80Frozen:       shipping.Box80Frozen,
			Box100Rates:       NewShippingRates(shipping.Box100Rates).Response(),
			Box100Frozen:      shipping.Box100Frozen,
			HasFreeShipping:   shipping.HasFreeShipping,
			FreeShippingRates: shipping.FreeShippingRates,
		},
	}
}

func (s *Shipping) Response() *response.Shipping {
	return &s.Shipping
}

type ShippingRate struct {
	response.ShippingRate
}

type ShippingRates []*ShippingRate

func NewShippingRate(rate *entity.ShippingRate) *ShippingRate {
	return &ShippingRate{
		ShippingRate: response.ShippingRate{
			Number:          rate.Number,
			Name:            rate.Name,
			Price:           rate.Price,
			Prefectures:     rate.Prefectures,
			PrefectureCodes: rate.PrefectureCodes,
		},
	}
}

func (r *ShippingRate) Response() *response.ShippingRate {
	return &r.ShippingRate
}

func NewShippingRates(rates []*entity.ShippingRate) ShippingRates {
	shippingRates := make(ShippingRates, 0, len(rates))
	for _, rate := range rates {
		shippingRates = append(shippingRates, NewShippingRate(rate))
	}
	return shippingRates
}

func (rs ShippingRates) Response() []*response.ShippingRate {
	shippingRates := make([]*response.ShippingRate, 0, len(rs))
	for _, rate := range rs {
		shippingRates = append(shippingRates, rate.Response())
	}
	return shippingRates
}
