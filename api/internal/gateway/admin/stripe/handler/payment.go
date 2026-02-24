package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/gin-gonic/gin"
	lib "github.com/stripe/stripe-go/v82"
)

func (h *handler) paymentIntentAuthorized(ctx *gin.Context, event *lib.Event) {
	slog.Debug("Received payment_intent.amount_capturable_updated event")
	pi, err := parsePaymentIntent(event)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.NotifyPaymentAuthorizedInput{
		NotifyPaymentPayload: store.NotifyPaymentPayload{
			OrderID:   pi.Metadata["order_id"],
			PaymentID: pi.ID,
			IssuedAt:  time.Unix(event.Created, 0),
			Status:    entity.PaymentStatusAuthorized,
		},
	}
	if err := h.store.NotifyPaymentAuthorized(ctx, in); err != nil {
		slog.Error("Failed payment_intent.amount_capturable_updated event",
			slog.String("paymentIntentId", pi.ID), log.Error(err))
		h.httpError(ctx, err)
		return
	}
	slog.Debug("Succeeded payment_intent.amount_capturable_updated event",
		slog.String("paymentIntentId", pi.ID))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) paymentIntentSucceeded(ctx *gin.Context, event *lib.Event) {
	slog.Debug("Received payment_intent.succeeded event")
	pi, err := parsePaymentIntent(event)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.NotifyPaymentCapturedInput{
		NotifyPaymentPayload: store.NotifyPaymentPayload{
			OrderID:   pi.Metadata["order_id"],
			PaymentID: pi.ID,
			IssuedAt:  time.Unix(event.Created, 0),
			Status:    entity.PaymentStatusCaptured,
		},
	}
	if err := h.store.NotifyPaymentCaptured(ctx, in); err != nil {
		slog.Error("Failed payment_intent.succeeded event",
			slog.String("paymentIntentId", pi.ID), log.Error(err))
		h.httpError(ctx, err)
		return
	}
	slog.Debug("Succeeded payment_intent.succeeded event",
		slog.String("paymentIntentId", pi.ID))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) paymentIntentFailed(ctx *gin.Context, event *lib.Event) {
	slog.Debug("Received payment_intent.payment_failed event")
	pi, err := parsePaymentIntent(event)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.NotifyPaymentFailedInput{
		NotifyPaymentPayload: store.NotifyPaymentPayload{
			OrderID:   pi.Metadata["order_id"],
			PaymentID: pi.ID,
			IssuedAt:  time.Unix(event.Created, 0),
			Status:    entity.PaymentStatusFailed,
		},
	}
	if err := h.store.NotifyPaymentFailed(ctx, in); err != nil {
		slog.Error("Failed payment_intent.payment_failed event",
			slog.String("paymentIntentId", pi.ID), log.Error(err))
		h.httpError(ctx, err)
		return
	}
	slog.Debug("Succeeded payment_intent.payment_failed event",
		slog.String("paymentIntentId", pi.ID))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) paymentIntentCanceled(ctx *gin.Context, event *lib.Event) {
	slog.Debug("Received payment_intent.canceled event")
	pi, err := parsePaymentIntent(event)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.NotifyPaymentRefundedInput{
		NotifyPaymentPayload: store.NotifyPaymentPayload{
			OrderID:   pi.Metadata["order_id"],
			PaymentID: pi.ID,
			IssuedAt:  time.Unix(event.Created, 0),
			Status:    entity.PaymentStatusCanceled,
		},
		Type: entity.RefundTypeCanceled,
	}
	if pi.CancellationReason != "" {
		in.Reason = string(pi.CancellationReason)
	}
	if err := h.store.NotifyPaymentRefunded(ctx, in); err != nil {
		slog.Error("Failed payment_intent.canceled event",
			slog.String("paymentIntentId", pi.ID), log.Error(err))
		h.httpError(ctx, err)
		return
	}
	slog.Debug("Succeeded payment_intent.canceled event",
		slog.String("paymentIntentId", pi.ID))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) chargeRefunded(ctx *gin.Context, event *lib.Event) {
	slog.Debug("Received charge.refunded event")
	charge, err := parseCharge(event)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	orderID := charge.Metadata["order_id"]
	paymentID := ""
	if charge.PaymentIntent != nil {
		paymentID = charge.PaymentIntent.ID
	}
	var refundTotal int64
	refundReasons := make([]string, 0)
	if charge.Refunds != nil {
		for _, r := range charge.Refunds.Data {
			refundTotal += r.Amount
			if r.Reason != "" {
				refundReasons = append(refundReasons, string(r.Reason))
			}
		}
	}
	in := &store.NotifyPaymentRefundedInput{
		NotifyPaymentPayload: store.NotifyPaymentPayload{
			OrderID:   orderID,
			PaymentID: paymentID,
			IssuedAt:  time.Unix(event.Created, 0),
			Status:    entity.PaymentStatusRefunded,
		},
		Type:   entity.RefundTypeRefunded,
		Reason: strings.Join(refundReasons, "\n"),
		Total:  refundTotal,
	}
	if err := h.store.NotifyPaymentRefunded(ctx, in); err != nil {
		slog.Error("Failed charge.refunded event",
			slog.String("chargeId", charge.ID), log.Error(err))
		h.httpError(ctx, err)
		return
	}
	slog.Debug("Succeeded charge.refunded event",
		slog.String("chargeId", charge.ID))
	ctx.Status(http.StatusNoContent)
}

func parsePaymentIntent(event *lib.Event) (*lib.PaymentIntent, error) {
	pi := &lib.PaymentIntent{}
	if err := json.Unmarshal(event.Data.Raw, pi); err != nil {
		return nil, err
	}
	return pi, nil
}

func parseCharge(event *lib.Event) (*lib.Charge, error) {
	charge := &lib.Charge{}
	if err := json.Unmarshal(event.Data.Raw, charge); err != nil {
		return nil, err
	}
	return charge, nil
}
