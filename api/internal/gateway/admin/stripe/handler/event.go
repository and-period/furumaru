package handler

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"

	lib "github.com/stripe/stripe-go/v82"

	"github.com/gin-gonic/gin"
)

const (
	eventTypePaymentIntentAmountCapturableUpdated = "payment_intent.amount_capturable_updated"
	eventTypePaymentIntentSucceeded               = "payment_intent.succeeded"
	eventTypePaymentIntentPaymentFailed           = "payment_intent.payment_failed"
	eventTypePaymentIntentCanceled                = "payment_intent.canceled"
	eventTypeChargeRefunded                       = "charge.refunded"
)

func (h *handler) Event(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	// Stripe Webhook署名検証
	signature := ctx.GetHeader("Stripe-Signature")
	event, err := h.receiver.Receive(body, signature)
	if err != nil {
		slog.Warn("Failed to verify stripe webhook signature", slog.String("error", err.Error()))
		ctx.Status(http.StatusUnauthorized)
		return
	}

	switch event.Type {
	case eventTypePaymentIntentAmountCapturableUpdated:
		h.paymentIntentAuthorized(ctx, event)
	case eventTypePaymentIntentSucceeded:
		h.paymentIntentSucceeded(ctx, event)
	case eventTypePaymentIntentPaymentFailed:
		h.paymentIntentFailed(ctx, event)
	case eventTypePaymentIntentCanceled:
		h.paymentIntentCanceled(ctx, event)
	case eventTypeChargeRefunded:
		h.chargeRefunded(ctx, event)
	default:
		h.unexpected(ctx, event)
	}
}

func (h *handler) unexpected(ctx *gin.Context, event *lib.Event) {
	req, err := httputil.DumpRequest(ctx.Request, false)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	slog.Debug("Received unexpected stripe event",
		slog.String("eventType", string(event.Type)),
		slog.Any("request", req),
	)
	ctx.Status(http.StatusNoContent)
}
