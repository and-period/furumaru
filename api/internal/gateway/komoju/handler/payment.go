package handler

import (
	"net/http"
	"strings"

	"github.com/and-period/furumaru/api/internal/gateway/komoju/request"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *handler) paymentAuthorized(ctx *gin.Context) {
	h.logger.Debug("Received payment authorized event")
	h.paymentCompleteEvent(ctx, entity.PaymentStatusAuthorized)
}

func (h *handler) paymentCaptured(ctx *gin.Context) {
	h.logger.Debug("Received payment captured event")
	h.paymentCompleteEvent(ctx, entity.PaymentStatusCaptured)
}

func (h *handler) paymentFailed(ctx *gin.Context) {
	h.logger.Debug("Received payment failed event")
	h.paymentCompleteEvent(ctx, entity.PaymentStatusFailed)
}

func (h *handler) paymentExpired(ctx *gin.Context) {
	h.logger.Debug("Received payment expired event")
	h.paymentCompleteEvent(ctx, entity.PaymentStatusFailed)
}

func (h *handler) paymentCancelled(ctx *gin.Context) {
	h.logger.Debug("Received payment cancelled event")
	h.paymentRefundEvent(ctx, entity.PaymentStatusCanceled)
}

func (h *handler) paymentRefunded(ctx *gin.Context) {
	h.logger.Debug("Received payment refunded event")
	h.paymentRefundEvent(ctx, entity.PaymentStatusRefunded)
}

func (h *handler) paymentCompleteEvent(ctx *gin.Context, status entity.PaymentStatus) {
	req := &request.PaymentRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.logger.Warn("Failed to parse request", zap.Error(err))
		h.badRequest(ctx, err)
		return
	}
	in := &store.NotifyPaymentCompletedInput{
		OrderID:   req.Payload.ExternalOrderNumber,
		PaymentID: req.Payload.ID,
		Status:    status,
		IssuedAt:  req.CreatedAt,
	}
	if err := h.store.NotifyPaymentCompleted(ctx, in); err != nil {
		h.logger.Error("Failed payment complete event", zap.Any("request", req), zap.Error(err))
		h.httpError(ctx, err)
		return
	}
	h.logger.Debug("Succeeded payment complete event", zap.Any("request", req))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) paymentRefundEvent(ctx *gin.Context, status entity.PaymentStatus) {
	req := &request.PaymentRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.logger.Warn("Failed to parse request", zap.Error(err))
		h.badRequest(ctx, err)
		return
	}
	var (
		refundTotal int64
		refundType  entity.RefundType
	)
	refundReason := make([]string, 0, len(req.Payload.Refunds))
	for _, refund := range req.Payload.Refunds {
		refundTotal += refund.Amount
		refundReason = append(refundReason, refund.Description)
	}
	switch status {
	case entity.PaymentStatusCanceled:
		refundType = entity.RefundTypeCanceled
	case entity.PaymentStatusRefunded:
		refundType = entity.RefundTypeRefunded
	}
	in := &store.NotifyPaymentRefundedInput{
		OrderID:  req.Payload.ExternalOrderNumber,
		Status:   status,
		Type:     refundType,
		Reason:   strings.Join(refundReason, "\n"),
		Total:    refundTotal,
		IssuedAt: req.CreatedAt,
	}
	if err := h.store.NotifyPaymentRefunded(ctx, in); err != nil {
		h.logger.Error("Failed payment refund event", zap.Any("request", req), zap.Error(err))
		h.httpError(ctx, err)
		return
	}
	h.logger.Debug("Succeeded payment refund event", zap.Any("request", req))
	ctx.Status(http.StatusNoContent)
}
