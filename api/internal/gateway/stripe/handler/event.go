package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v73"
	"go.uber.org/zap"
)

func (h *handler) Event(ctx *gin.Context) {
	event, err := h.request(ctx.Writer, ctx.Request)
	if err != nil {
		badRequest(ctx, err)
		return
	}

	h.logger.Debug("Received Event", zap.Any("event", event))

	if err := h.dispatch(ctx, event); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) request(rw http.ResponseWriter, r *http.Request) (*stripe.Event, error) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return h.receiver.Receive(payload, r.Header.Get("Stripe-Signature"))
}

func (h *handler) dispatch(ctx context.Context, event *stripe.Event) error {
	if event == nil {
		return fmt.Errorf("%s: %w", errUnknownEvent.Error(), exception.ErrInvalidArgument)
	}
	switch event.Type {
	case "payment_intent.payment_failed":
	default:
		// エラーを返すと常にリトライが走ってしまうため、ログ出力のみして成功ステータスを返す
		h.logger.Warn("Received not implemented event", zap.Any("type", event.Type))
	}
	return nil
}
