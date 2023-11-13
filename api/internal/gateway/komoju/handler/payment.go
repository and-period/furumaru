package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/komoju/request"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *handler) paymentAuthorized(ctx *gin.Context) {
	req := &request.PaymentRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	h.logger.Debug("Received payment authorized event", zap.Any("request", req))
	in := &store.NotifyPaymentCompletedInput{
		OrderID:   req.Payload.ExternalOrderNumber,
		PaymentID: req.Payload.ID,
		Status:    req.Payload.Status,
		IssuedAt:  req.CreatedAt,
	}
	if err := h.store.NotifyPaymentCompleted(ctx, in); err != nil {
		h.logger.Error("Failed payment authorized event", zap.Any("request", req), zap.Error(err))
		h.httpError(ctx, err)
		return
	}
	h.logger.Debug("Complete payment authorized event", zap.Any("request", req))
	ctx.Status(http.StatusNoContent)
}
