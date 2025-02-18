package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

func (s *service) CalcCart(ctx context.Context, in *store.CalcCartInput) (*entity.Cart, *entity.OrderPaymentSummary, error) {
	s.logger.Debug("calc cart", zap.Any("input", in))
	if err := s.validator.Struct(in); err != nil {
		return nil, nil, internalError(err)
	}
	shop, err := s.db.Shop.GetByCoordinatorID(ctx, in.CoordinatorID)
	if err != nil {
		return nil, nil, internalError(err)
	}
	var (
		shipping  *entity.Shipping
		cart      *entity.Cart
		promotion *entity.Promotion
	)
	eg, ectx := errgroup.WithContext(ctx)
	// カートの取得
	eg.Go(func() (err error) {
		cart, err = s.getCart(ectx, in.SessionID)
		return
	})
	// 配送設定の取得
	eg.Go(func() (err error) {
		shipping, err = s.getShippingByCoordinatorID(ectx, in.CoordinatorID)
		return
	})
	// プロモーションの取得
	eg.Go(func() (err error) {
		if in.PromotionCode == "" {
			return
		}
		promotion, err = s.db.Promotion.GetByCode(ectx, in.PromotionCode)
		if promotion.IsEnabled(shop.ID) {
			return
		}
		s.logger.Warn("Failed to disable promotion", zap.String("sessionId", in.SessionID), zap.String("code", in.PromotionCode))
		promotion = nil
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, nil, internalError(err)
	}
	// 購入する買い物かごのみ取得
	baskets := cart.Baskets.FilterByCoordinatorID(in.CoordinatorID).FilterByBoxNumber(in.BoxNumber)
	if len(baskets) == 0 {
		return nil, nil, fmt.Errorf("service: no target baskets: %w", exception.ErrInvalidArgument)
	}
	// 在庫一覧を取得
	products, err := s.db.Product.MultiGet(ctx, baskets.ProductIDs())
	if err != nil {
		return nil, nil, internalError(err)
	}
	params := &entity.NewProductOrderPaymentSummaryParams{
		PrefectureCode: in.PrefectureCode,
		Baskets:        baskets,
		Products:       products,
		Shipping:       shipping,
		Promotion:      promotion,
	}
	summary, err := entity.NewProductOrderPaymentSummary(params)
	if err != nil {
		return nil, nil, internalError(err)
	}
	cart.Baskets = baskets
	return cart, summary, nil
}

func (s *service) AddCartItem(ctx context.Context, in *store.AddCartItemInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	product, err := s.db.Product.Get(ctx, in.ProductID)
	if err != nil {
		return internalError(err)
	}
	if product.Status != entity.ProductStatusForSale {
		return fmt.Errorf("service: this product is out of sale: %w", exception.ErrForbidden)
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
	if err := s.refreshCart(ctx, cart); err != nil {
		return internalError(err)
	}
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		params := &entity.AddCartItemActionLogParams{
			SessionID: in.SessionID,
			UserID:    in.UserID,
			UserAgent: in.UserAgent,
			ClientIP:  in.ClientIP,
			ProductID: in.ProductID,
		}
		log := entity.NewAddCartItemActionLog(params)
		if err := s.db.CartActionLog.Create(context.Background(), log); err != nil {
			s.logger.Warn("Failed to create cart action log", zap.Error(err))
		}
	}()
	return nil
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
	if err := s.refreshCart(ctx, cart); err != nil {
		return internalError(err)
	}
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		params := &entity.RemoveCartItemActionLogParams{
			SessionID: in.SessionID,
			UserID:    in.UserID,
			UserAgent: in.UserAgent,
			ClientIP:  in.ClientIP,
			ProductID: in.ProductID,
		}
		log := entity.NewRemoveCartItemActionLog(params)
		if err := s.db.CartActionLog.Create(context.Background(), log); err != nil {
			s.logger.Warn("Failed to create cart action log", zap.Error(err))
		}
	}()
	return nil
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
	if err := cart.Refresh(products.FilterBySales()); err != nil {
		return err
	}
	now := s.now()
	if cart.CreatedAt.IsZero() {
		cart.CreatedAt = now
	}
	cart.UpdatedAt = now
	cart.ExpiredAt = now.Add(s.cartTTL) // 有効期限を延長
	fn := func() error {
		return s.cache.Insert(ctx, cart)
	}
	const maxRetries = 5
	retry := backoff.NewExponentialBackoff(maxRetries)
	return backoff.Retry(ctx, retry, fn, backoff.WithRetryablel(dynamodb.IsRetryable))
}
