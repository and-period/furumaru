package service

import (
	"context"
	"errors"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"go.uber.org/zap"
)

func (s *service) NotifyPaymentCompleted(ctx context.Context, in *store.NotifyPaymentCompletedInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.UpdateOrderPaymentParams{
		Status:    in.Status,
		PaymentID: in.PaymentID,
		IssuedAt:  in.IssuedAt,
	}
	err := s.db.Order.UpdatePayment(ctx, in.OrderID, params)
	if errors.Is(err, database.ErrFailedPrecondition) {
		s.logger.Warn("Order is already updated",
			zap.String("orderId", in.OrderID), zap.Int32("status", int32(in.Status)), zap.Time("issuedAt", in.IssuedAt))
		return nil
	}
	if err != nil {
		return internalError(err)
	}
	// 即時決済にも対応できるよう、支払い完了タイミングで通知を投げるように
	if in.Status != entity.PaymentStatusCaptured {
		return nil
	}
	notifyIn := &messenger.NotifyOrderAuthorizedInput{
		OrderID: in.OrderID,
	}
	err = s.messenger.NotifyOrderAuthorized(ctx, notifyIn)
	return internalError(err)
}

func (s *service) NotifyPaymentRefunded(ctx context.Context, in *store.NotifyPaymentRefundedInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.UpdateOrderRefundParams{
		Status:       in.Status,
		RefundType:   in.Type,
		RefundTotal:  in.Total,
		RefundReason: in.Reason,
		IssuedAt:     in.IssuedAt,
	}
	err := s.db.Order.UpdateRefund(ctx, in.OrderID, params)
	if errors.Is(err, database.ErrFailedPrecondition) {
		s.logger.Warn("Order is already updated",
			zap.String("orderId", in.OrderID), zap.Int32("status", int32(in.Status)), zap.Time("issuedAt", in.IssuedAt))
		return nil
	}
	return internalError(err)
}
