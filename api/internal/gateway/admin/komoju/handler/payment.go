package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/and-period/furumaru/api/internal/gateway/admin/komoju/request"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/gin-gonic/gin"
)

func (h *handler) paymentAuthorized(ctx *gin.Context) {
	slog.Debug("Received payment authorized event")
	req, err := h.newPaymentRequest(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.NotifyPaymentAuthorizedInput{
		NotifyPaymentPayload: *h.newNotifyPaymentPayload(req, entity.PaymentStatusAuthorized),
	}
	if err := h.store.NotifyPaymentAuthorized(ctx, in); err != nil {
		slog.Error("Failed payment authorized event", slog.Any("request", req), log.Error(err))
		h.httpError(ctx, err)
		return
	}
	slog.Debug("Succeeded payment authorized event", slog.Any("request", req))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) paymentCaptured(ctx *gin.Context) {
	slog.Debug("Received payment captured event")
	req, err := h.newPaymentRequest(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.NotifyPaymentCapturedInput{
		NotifyPaymentPayload: *h.newNotifyPaymentPayload(req, entity.PaymentStatusCaptured),
	}
	if err := h.store.NotifyPaymentCaptured(ctx, in); err != nil {
		slog.Error("Failed payment captured event", slog.Any("request", req), log.Error(err))
		h.httpError(ctx, err)
		return
	}
	slog.Debug("Succeeded payment captured event", slog.Any("request", req))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) paymentFailed(ctx *gin.Context) {
	slog.Debug("Received payment failed event")
	req, err := h.newPaymentRequest(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.NotifyPaymentFailedInput{
		NotifyPaymentPayload: *h.newNotifyPaymentPayload(req, entity.PaymentStatusFailed),
	}
	if err := h.store.NotifyPaymentFailed(ctx, in); err != nil {
		slog.Error("Failed payment failed event", slog.Any("request", req), log.Error(err))
		h.httpError(ctx, err)
		return
	}
	slog.Debug("Succeeded payment failed event", slog.Any("request", req))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) paymentExpired(ctx *gin.Context) {
	slog.Debug("Received payment expired event")
	req, err := h.newPaymentRequest(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.NotifyPaymentFailedInput{
		NotifyPaymentPayload: *h.newNotifyPaymentPayload(req, entity.PaymentStatusExpired),
	}
	if err := h.store.NotifyPaymentFailed(ctx, in); err != nil {
		slog.Error("Failed payment expired event", slog.Any("request", req), log.Error(err))
		h.httpError(ctx, err)
		return
	}
	slog.Debug("Succeeded payment expired event", slog.Any("request", req))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) paymentCancelled(ctx *gin.Context) {
	slog.Debug("Received payment cancelled event")
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
		slog.Error("Failed payment cancelled event", slog.Any("request", req), log.Error(err))
		h.httpError(ctx, err)
		return
	}
	slog.Debug("Succeeded payment cancelled event", slog.Any("request", req))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) paymentRefunded(ctx *gin.Context) {
	slog.Debug("Received payment refunded event")
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
		slog.Error("Failed payment refunded event", slog.Any("request", req), log.Error(err))
		h.httpError(ctx, err)
		return
	}
	slog.Debug("Succeeded payment refunded event", slog.Any("request", req))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) newPaymentRequest(ctx *gin.Context) (*request.PaymentRequest, error) {
	req := &request.PaymentRequest{}
	if err := ctx.BindJSON(req); err != nil {
		slog.Warn("Failed to parse request", log.Error(err))
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
