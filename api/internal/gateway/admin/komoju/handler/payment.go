package handler

import (
	"net/http"
	"strings"

	"github.com/and-period/furumaru/api/internal/gateway/admin/komoju/request"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *handler) paymentAuthorized(ctx *gin.Context) {
	h.logger.Debug("Received payment authorized event")
	req, err := h.newPaymentRequest(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.NotifyPaymentAuthorizedInput{
		NotifyPaymentPayload: *h.newNotifyPaymentPayload(req, entity.PaymentStatusAuthorized),
	}
	if err := h.store.NotifyPaymentAuthorized(ctx, in); err != nil {
		h.logger.Error("Failed payment authorized event", zap.Any("request", req), zap.Error(err))
		h.httpError(ctx, err)
		return
	}
	h.logger.Debug("Succeeded payment authorized event", zap.Any("request", req))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) paymentCaptured(ctx *gin.Context) {
	h.logger.Debug("Received payment captured event")
	req, err := h.newPaymentRequest(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.NotifyPaymentCapturedInput{
		NotifyPaymentPayload: *h.newNotifyPaymentPayload(req, entity.PaymentStatusCaptured),
	}
	if err := h.store.NotifyPaymentCaptured(ctx, in); err != nil {
		h.logger.Error("Failed payment captured event", zap.Any("request", req), zap.Error(err))
		h.httpError(ctx, err)
		return
	}
	h.logger.Debug("Succeeded payment captured event", zap.Any("request", req))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) paymentFailed(ctx *gin.Context) {
	h.logger.Debug("Received payment failed event")
	req, err := h.newPaymentRequest(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.NotifyPaymentFailedInput{
		NotifyPaymentPayload: *h.newNotifyPaymentPayload(req, entity.PaymentStatusFailed),
	}
	if err := h.store.NotifyPaymentFailed(ctx, in); err != nil {
		h.logger.Error("Failed payment failed event", zap.Any("request", req), zap.Error(err))
		h.httpError(ctx, err)
		return
	}
	h.logger.Debug("Succeeded payment failed event", zap.Any("request", req))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) paymentExpired(ctx *gin.Context) {
	h.logger.Debug("Received payment expired event")
	req, err := h.newPaymentRequest(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.NotifyPaymentFailedInput{
		NotifyPaymentPayload: *h.newNotifyPaymentPayload(req, entity.PaymentStatusExpired),
	}
	if err := h.store.NotifyPaymentFailed(ctx, in); err != nil {
		h.logger.Error("Failed payment expired event", zap.Any("request", req), zap.Error(err))
		h.httpError(ctx, err)
		return
	}
	h.logger.Debug("Succeeded payment expired event", zap.Any("request", req))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) paymentCancelled(ctx *gin.Context) {
	h.logger.Debug("Received payment cancelled event")
	req, err := h.newPaymentRequest(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	var refundTotal int64
	refundReason := make([]string, 0, len(req.Payload.Refunds))
	for _, refund := range req.Payload.Refunds {
		refundTotal += refund.Amount
		refundReason = append(refundReason, refund.Description)
	}
	in := &store.NotifyPaymentRefundedInput{
		NotifyPaymentPayload: *h.newNotifyPaymentPayload(req, entity.PaymentStatusCanceled),
		Type:                 entity.RefundTypeCanceled,
		Reason:               strings.Join(refundReason, "\n"),
		Total:                refundTotal,
	}
	if err := h.store.NotifyPaymentRefunded(ctx, in); err != nil {
		h.logger.Error("Failed payment cancelled event", zap.Any("request", req), zap.Error(err))
		h.httpError(ctx, err)
		return
	}
	h.logger.Debug("Succeeded payment cancelled event", zap.Any("request", req))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) paymentRefunded(ctx *gin.Context) {
	h.logger.Debug("Received payment refunded event")
	req, err := h.newPaymentRequest(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	var refundTotal int64
	refundReason := make([]string, 0, len(req.Payload.Refunds))
	for _, refund := range req.Payload.Refunds {
		refundTotal += refund.Amount
		refundReason = append(refundReason, refund.Description)
	}
	in := &store.NotifyPaymentRefundedInput{
		NotifyPaymentPayload: *h.newNotifyPaymentPayload(req, entity.PaymentStatusRefunded),
		Type:                 entity.RefundTypeRefunded,
		Reason:               strings.Join(refundReason, "\n"),
		Total:                refundTotal,
	}
	if err := h.store.NotifyPaymentRefunded(ctx, in); err != nil {
		h.logger.Error("Failed payment refunded event", zap.Any("request", req), zap.Error(err))
		h.httpError(ctx, err)
		return
	}
	h.logger.Debug("Succeeded payment refunded event", zap.Any("request", req))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) newPaymentRequest(ctx *gin.Context) (*request.PaymentRequest, error) {
	req := &request.PaymentRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.logger.Warn("Failed to parse request", zap.Error(err))
		return nil, err
	}
	return req, nil
}

func (h *handler) newNotifyPaymentPayload(req *request.PaymentRequest, status entity.PaymentStatus) *store.NotifyPaymentPayload {
	return &store.NotifyPaymentPayload{
		OrderID:   req.Payload.ExternalOrderNumber,
		PaymentID: req.Payload.ID,
		IssuedAt:  req.CreatedAt,
		Status:    status,
	}
}
