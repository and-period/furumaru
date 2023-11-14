package handler

import (
	"net/http"

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

func (h *handler) paymentCancelled(ctx *gin.Context) {
	h.logger.Debug("Received payment cancelled event")
	h.paymentCompleteEvent(ctx, entity.PaymentStatusCanceled)
}

func (h *handler) paymentFailed(ctx *gin.Context) {
	h.logger.Debug("Received payment failed event")
	h.paymentCompleteEvent(ctx, entity.PaymentStatusFailed)
}

func (h *handler) paymentExpired(ctx *gin.Context) {
	h.logger.Debug("Received payment expired event")
	h.paymentCompleteEvent(ctx, entity.PaymentStatusFailed)
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
		h.logger.Error("Failed payment event", zap.Any("request", req), zap.Error(err))
		h.httpError(ctx, err)
		return
	}
	h.logger.Debug("Complete payment event", zap.Any("request", req))
	ctx.Status(http.StatusNoContent)
}
