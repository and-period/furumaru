package service

import (
	"context"

	"github.com/and-period/marche/api/internal/user/database"
	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/and-period/marche/api/pkg/random"
	"github.com/and-period/marche/api/pkg/uuid"
)

func (s *userService) ListShops(ctx context.Context, in *ListShopsInput) (entity.Shops, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	params := &database.ListShopsParams{
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	shops, err := s.db.Shop.List(ctx, params)
	return shops, userError(err)
}

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

func (s *userService) CreateShop(ctx context.Context, in *CreateShopInput) (*entity.Shop, error) {
	const size = 8
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	shopID := uuid.Base58Encode(uuid.New())
	shop := entity.NewShop(shopID, shopID, in.Lastname, in.Firstname, in.LastnameKana, in.FirstnameKana, in.Email)
	if err := s.db.Shop.Create(ctx, shop); err != nil {
		return nil, userError(err)
	}
	password := random.NewStrings(size)
	params := &cognito.AdminCreateUserParams{
		Username: shop.CognitoID,
		Email:    shop.Email,
		Password: password,
	}
	if err := s.shopAuth.AdminCreateUser(ctx, params); err != nil {
		return nil, userError(err)
	}
	// TODO: 販売者登録通知を送信
	return shop, nil
}
