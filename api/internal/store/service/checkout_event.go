package service

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

func (s *service) NotifyPaymentAuthorized(ctx context.Context, in *store.NotifyPaymentAuthorizedInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	if in.Status != entity.PaymentStatusAuthorized {
		return fmt.Errorf("this status is not authorized. status=%d: %w", in.Status, exception.ErrInvalidArgument)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	// KOMOJU側の仕様として、即時決済出ないものは支払い待ちステータスでAuthorizedの通知が来るためスキップさせる
	if order.OrderPayment.IsDeferredPayment() {
		s.logger.Info("Order is authorized but deferred payment", zap.String("orderId", in.OrderID))
		return nil
	}
	params := &database.UpdateOrderAuthorizedParams{
		PaymentID: in.PaymentID,
		IssuedAt:  in.IssuedAt,
	}
	err = s.db.Order.UpdateAuthorized(ctx, in.OrderID, params)
	if err != nil && !errors.Is(err, database.ErrFailedPrecondition) {
		return internalError(err)
	}
	_, err = s.komoju.Payment.Capture(ctx, in.PaymentID)
	if err != nil && komoju.IsRetryable(err) {
		return fmt.Errorf("store: failed to capture payment. err=%s: %w", err.Error(), exception.ErrInternal)
	}
	if err != nil {
		s.logger.Warn("Failed to capture payment", zap.String("paymentID", in.PaymentID), zap.Error(err))
		return nil // Webhookがリトライしないように成功として返す
	}
	// 即時決済の場合、買い物かごからの削除処理も追加で行う
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		if !order.OrderPayment.IsImmediatePayment() {
			return
		}
		if err := s.removeCartItemByOrder(ctx, order); err != nil {
			s.logger.Error("Failed to remove cart item by order", zap.String("orderId", in.OrderID), zap.Error(err))
		}
	}()
	return nil
}

func (s *service) NotifyPaymentCaptured(ctx context.Context, in *store.NotifyPaymentCapturedInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	if in.Status != entity.PaymentStatusCaptured {
		return fmt.Errorf("this status is not captured. status=%d: %w", in.Status, exception.ErrInvalidArgument)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	params := &database.UpdateOrderCapturedParams{
		PaymentID: in.PaymentID,
		IssuedAt:  in.IssuedAt,
	}
	err = s.db.Order.UpdateCaptured(ctx, in.OrderID, params)
	if errors.Is(err, database.ErrFailedPrecondition) {
		s.logger.Warn("Order can't be captured", zap.String("orderId", in.OrderID), zap.Time("issuedAt", in.IssuedAt))
		return nil
	}
	if err != nil {
		return internalError(err)
	}
	if err := s.notifyPaymentCompleted(ctx, order); err != nil {
		return internalError(err)
	}
	return nil
}

func (s *service) NotifyPaymentFailed(ctx context.Context, in *store.NotifyPaymentFailedInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	if !slices.Contains(entity.PaymentFailedStatuses, in.Status) {
		return fmt.Errorf("this status is not failed. status=%d: %w", in.Status, exception.ErrInvalidArgument)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	params := &database.UpdateOrderFailedParams{
		PaymentID: in.PaymentID,
		Status:    in.Status,
		IssuedAt:  in.IssuedAt,
	}
	err = s.db.Order.UpdateFailed(ctx, in.OrderID, params)
	if errors.Is(err, database.ErrFailedPrecondition) {
		s.logger.Warn("Order can't be failed", zap.String("orderId", in.OrderID), zap.Time("issuedAt", in.IssuedAt))
		return nil
	}
	if err != nil {
		return internalError(err)
	}
	s.waitGroup.Add(1)
	// 確保していた商品在庫の開放
	go func() {
		defer s.waitGroup.Done()
		s.increaseProductInventories(context.Background(), order.OrderItems)
	}()
	return nil
}

func (s *service) NotifyPaymentRefunded(ctx context.Context, in *store.NotifyPaymentRefundedInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.UpdateOrderRefundedParams{
		Status:       in.Status,
		RefundType:   in.Type,
		RefundTotal:  in.Total,
		RefundReason: in.Reason,
		IssuedAt:     in.IssuedAt,
	}
	err := s.db.Order.UpdateRefunded(ctx, in.OrderID, params)
	if errors.Is(err, database.ErrFailedPrecondition) {
		s.logger.Warn("Order can't be refunded", zap.String("orderId", in.OrderID), zap.Time("issuedAt", in.IssuedAt))
		return nil
	}
	return internalError(err)
}

func (s *service) notifyPaymentCompleted(ctx context.Context, order *entity.Order) error {
	notifyIn := &messenger.NotifyOrderCapturedInput{
		OrderID: order.ID,
	}
	if err := s.messenger.NotifyOrderCaptured(ctx, notifyIn); err != nil {
		return fmt.Errorf("failed to notify order captured: %w", err)
	}
	return nil
}

func (s *service) removeCartItemByOrder(ctx context.Context, order *entity.Order) error {
	cart, err := s.getCart(ctx, order.SessionID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}
	if len(cart.Baskets) == 0 {
		s.logger.Debug("Cart is empty", zap.String("sessionID", order.SessionID))
		return nil
	}
	revisionIDs := order.OrderItems.ProductRevisionIDs()
	products, err := s.db.Product.MultiGetByRevision(ctx, revisionIDs)
	if err != nil {
		return fmt.Errorf("failed to get products: %w", err)
	}
	productMap := products.MapByRevision()
	for _, item := range order.OrderItems {
		product, ok := productMap[item.ProductRevisionID]
		if !ok {
			continue
		}
		cart.DecreaseItem(product.ID, item.Quantity)
	}
	return s.refreshCart(ctx, cart)
}

func (s *service) increaseProductInventories(ctx context.Context, items entity.OrderItems) {
	sem := semaphore.NewWeighted(3)
	for _, item := range items {
		if err := sem.Acquire(ctx, 1); err != nil {
			s.logger.Error("Failed to acquire semaphore", zap.Any("item", item), zap.Error(err))
			return
		}
		s.waitGroup.Add(1)
		go func(item *entity.OrderItem) {
			defer func() {
				sem.Release(1)
				s.waitGroup.Done()
			}()
			err := s.decreaseProductInventory(ctx, item.ProductRevisionID, -item.Quantity)
			if err != nil {
				s.logger.Error("Failed to increase product inventory", zap.Any("item", item), zap.Error(err))
			}
		}(item)
	}
}
