package service

import (
	"context"

	"github.com/and-period/marche/api/internal/user/entity"
)

func (s *userService) MultiGetShops(ctx context.Context, in *MultiGetShopsInput) (entity.Shops, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	shops, err := s.db.Shop.MultiGet(ctx, in.ShopIDs)
	return shops, userError(err)
}

func (s *userService) GetShop(ctx context.Context, in *GetShopInput) (*entity.Shop, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	shop, err := s.db.Shop.Get(ctx, in.ShopID)
	return shop, userError(err)
}
