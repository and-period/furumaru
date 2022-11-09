package service

import (
	"github.com/and-period/furumaru/api/internal/codes"
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

func NewShipping(shipping *entity.Shipping) (*Shipping, error) {
	box60Rates, err := NewShippingRates(shipping.Box60Rates)
	if err != nil {
		return nil, err
	}
	box80Rates, err := NewShippingRates(shipping.Box80Rates)
	if err != nil {
		return nil, err
	}
	box100Rates, err := NewShippingRates(shipping.Box100Rates)
	if err != nil {
		return nil, err
	}
	return &Shipping{
		Shipping: response.Shipping{
			ID:                 shipping.ID,
			Name:               shipping.Name,
			Box60Rates:         box60Rates.Response(),
			Box60Refrigerated:  shipping.Box60Refrigerated,
			Box60Frozen:        shipping.Box60Frozen,
			Box80Rates:         box80Rates.Response(),
			Box80Refrigerated:  shipping.Box80Refrigerated,
			Box80Frozen:        shipping.Box80Frozen,
			Box100Rates:        box100Rates.Response(),
			Box100Refrigerated: shipping.Box100Refrigerated,
			Box100Frozen:       shipping.Box100Frozen,
			HasFreeShipping:    shipping.HasFreeShipping,
			FreeShippingRates:  shipping.FreeShippingRates,
			CreatedAt:          shipping.CreatedAt.Unix(),
			UpdatedAt:          shipping.CreatedAt.Unix(),
		},
	}, nil
}

func (s *Shipping) Response() *response.Shipping {
	return &s.Shipping
}

func NewShippings(shippings entity.Shippings) (Shippings, error) {
	res := make(Shippings, len(shippings))
	for i := range shippings {
		shipping, err := NewShipping(shippings[i])
		if err != nil {
			return nil, err
		}
		res[i] = shipping
	}
	return res, nil
}

func (ss Shippings) Response() []*response.Shipping {
	res := make([]*response.Shipping, len(ss))
	for i := range ss {
		res[i] = ss[i].Response()
	}
	return res
}

func NewShippingRate(rate *entity.ShippingRate) (*ShippingRate, error) {
	prefectures, err := codes.ToPrefectureNames(rate.Prefectures...)
	if err != nil {
		return nil, err
	}
	return &ShippingRate{
		ShippingRate: response.ShippingRate{
			Number:      rate.Number,
			Name:        rate.Name,
			Price:       rate.Price,
			Prefectures: prefectures,
		},
	}, nil
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

func NewShippingRates(rates entity.ShippingRates) (ShippingRates, error) {
	res := make(ShippingRates, len(rates))
	for i := range rates {
		rate, err := NewShippingRate(rates[i])
		if err != nil {
			return nil, err
		}
		res[i] = rate
	}
	return res, nil
}

func (rs ShippingRates) Response() []*response.ShippingRate {
	res := make([]*response.ShippingRate, len(rs))
	for i := range rs {
		res[i] = rs[i].Response()
	}
	return res
}
