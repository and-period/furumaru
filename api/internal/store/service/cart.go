package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
)

func (s *service) GetCart(ctx context.Context, in *store.GetCartInput) (*entity.Cart, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	cart, err := s.getCart(ctx, in.SessionID)
	if err != nil {
		return nil, internalError(err)
	}
	if s.now().Sub(cart.UpdatedAt) <= s.cartRefreshInterval {
		return cart, nil
	}
	// 最終更新時間から指定時間経過している場合、カートの中身を整理する
	if err := s.refreshCart(ctx, cart); err != nil {
		return nil, internalError(err)
	}
	return cart, nil
}

func (s *service) AddCartItem(ctx context.Context, in *store.AddCartItemInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	product, err := s.db.Product.Get(ctx, in.ProductID)
	if err != nil {
		return internalError(err)
	}
	if !product.Public {
		return fmt.Errorf("service: this product is not published: %w", exception.ErrForbidden)
	}
	cart, err := s.getCart(ctx, in.SessionID)
	if err != nil {
		return internalError(err)
	}
	if err := cart.Baskets.VerifyQuantity(in.Quantity, product); err != nil {
		return fmt.Errorf("service: out of stock: %w: %s", exception.ErrFailedPrecondition, err.Error())
	}
	// カートに商品を追加し、買い物かご内を整理
	cart.AddItem(in.ProductID, in.Quantity)
	err = s.refreshCart(ctx, cart)
	return internalError(err)
}

func (s *service) RemoveCartItem(ctx context.Context, in *store.RemoveCartItemInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	cart, err := s.getCart(ctx, in.SessionID)
	if err != nil {
		return internalError(err)
	}
	// カートから商品削除し、買い物かご内を整理
	cart.RemoveItem(in.ProductID, in.BoxNumber)
	err = s.refreshCart(ctx, cart)
	return internalError(err)
}

// getCart - カートを取得する もしくは 新規で登録する
func (s *service) getCart(ctx context.Context, sessionID string) (*entity.Cart, error) {
	cart := &entity.Cart{SessionID: sessionID}
	err := s.cache.Get(ctx, cart)
	if err == nil || !errors.Is(err, dynamodb.ErrNotFound) {
		return cart, err
	}
	params := &entity.CartParams{
		SessionID: sessionID,
		Now:       s.now(),
		TTL:       s.cartTTL,
	}
	cart = entity.NewCart(params)
	return cart, s.cache.Insert(ctx, cart)
}

// refreshCart - カートの中身を更新する
func (s *service) refreshCart(ctx context.Context, cart *entity.Cart) error {
	products, err := s.db.Product.MultiGet(ctx, cart.Baskets.ProductIDs())
	if err != nil {
		return err
	}
	if err := cart.Refresh(products.FilterByPublished()); err != nil {
		return err
	}
	now := s.now()
	cart.UpdatedAt = now
	cart.ExpiredAt = now.Add(s.cartTTL) // 有効期限を延長
	return s.cache.Insert(ctx, cart)
}
