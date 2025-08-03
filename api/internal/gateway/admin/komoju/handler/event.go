package handler

import (
	"log/slog"
	"net/http"
	"net/http/httputil"

	"github.com/and-period/furumaru/api/internal/gateway/admin/komoju/request"
	"github.com/gin-gonic/gin"
)

type EventType string

const (
	EventTypePing              EventType = "ping"
	EventTypePaymentAuthorized EventType = "payment.authorized"
	EventTypePaymentCaptured   EventType = "payment.captured"
	EventTypePaymentCancelled  EventType = "payment.cancelled"
	EventTypePaymentRefunded   EventType = "payment.refunded"
	EventTypePaymentFailed     EventType = "payment.failed"
	EventTypePaymentExpired    EventType = "payment.expired"
)

func (h *handler) Event(ctx *gin.Context) {
	event := ctx.GetHeader("X-Komoju-Event")
	switch EventType(event) {
	case EventTypePing:
		h.ping(ctx)
	case EventTypePaymentAuthorized:
		h.paymentAuthorized(ctx)
	case EventTypePaymentCaptured:
		h.paymentCaptured(ctx)
	case EventTypePaymentCancelled:
		h.paymentCancelled(ctx)
	case EventTypePaymentRefunded:
		h.paymentRefunded(ctx)
	case EventTypePaymentFailed:
		h.paymentFailed(ctx)
	case EventTypePaymentExpired:
		h.paymentExpired(ctx)
	default:
		h.unexpected(ctx, event)
	}
}

func (h *handler) ping(ctx *gin.Context) {
	req := &request.PingRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	slog.Debug("Received ping event", slog.Any("request", req))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) unexpected(ctx *gin.Context, event string) {
	req, err := httputil.DumpRequest(ctx.Request, true)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	slog.Debug("Received unexpected event", slog.String("event", event), slog.Any("request", req))
	ctx.Status(http.StatusNoContent)
}
