package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/set"
)

type Shop struct {
	types.Shop
}

type Shops []*Shop

func NewShop(shop *entity.Shop) *Shop {
	return &Shop{
		Shop: types.Shop{
			ID:             shop.ID,
			Name:           shop.Name,
			CoordinatorID:  shop.CoordinatorID,
			ProducerIDs:    shop.ProducerIDs,
			ProductTypeIDs: shop.ProductTypeIDs,
			BusinessDays:   shop.BusinessDays,
			CreatedAt:      shop.CreatedAt.Unix(),
			UpdatedAt:      shop.UpdatedAt.Unix(),
		},
	}
}

func (s *Shop) GetID() string {
	if s == nil {
		return ""
	}
	return s.ID
}

func (s *Shop) Response() *types.Shop {
	return &s.Shop
}

func NewShops(shops entity.Shops) Shops {
	res := make([]*Shop, len(shops))
	for i, shop := range shops {
		res[i] = NewShop(shop)
	}
	return res
}

func (ss Shops) CoordinatorIDs() []string {
	return set.UniqBy(ss, func(s *Shop) string {
		return s.CoordinatorID
	})
}

func (ss Shops) ProductTypeIDs() []string {
	set := set.NewEmpty[string](len(ss))
	for _, shop := range ss {
		set.Add(shop.ProductTypeIDs...)
	}
	return set.Slice()
}

func (ss Shops) MapByCoordinatorID() map[string]*Shop {
	res := make(map[string]*Shop, len(ss))
	for _, s := range ss {
		res[s.CoordinatorID] = s
	}
	return res
}

func (ss Shops) Response() []*types.Shop {
	res := make([]*types.Shop, len(ss))
	for i, shop := range ss {
		res[i] = shop.Response()
	}
	return res
}
