package service

import "github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"

type Shipping struct {
	response.Shipping
}

type Shippings []*Shipping

func (s *Shipping) Response() *response.Shipping {
	return &s.Shipping
}

func (ss Shippings) Response() []*response.Shipping {
	res := make([]*response.Shipping, len(ss))
	for i := range ss {
		res[i] = ss[i].Response()
	}
	return res
}
