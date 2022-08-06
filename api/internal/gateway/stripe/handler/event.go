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
	event, err := h.request(ctx.Request)
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

func (h *handler) request(r *http.Request) (*stripe.Event, error) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return h.receiver.Receive(payload, r.Header.Get("Stripe-Signature"))
}

func (h *handler) dispatch(_ context.Context, event *stripe.Event) error {
	if event == nil {
		return fmt.Errorf("%s: %w", errUnknownEvent.Error(), exception.ErrInvalidArgument)
	}
	switch event.Type {
	case "setup_intent.succeeded": // クレジットカード登録成功
		// TODO: デフォルト決済手段として登録
	case "payment_intent.amount_capturable_updated": // 決済実行(オーソリ)
		// TODO: 注文を認証済みに
	case "payment_intent.succeeded": // 決済確定(キャプチャ)
		// TODO: 注文を成功に
	case "payment_intent.canceled": // 決済キャンセル
		// TODO: 注文をキャンセルに
	case "payment_intent.payment_failed": // 決済失敗
		// TODO: 注文を失敗に
	case
		"customer.created",               // 顧客登録
		"setup_intent.created",           // クレジットカード作成成功
		"setup_intent.setup_failed",      // クレジットカード登録失敗
		"payment_created",                // 決済要求
		"payment_intent.requires_action": // 決済実行(未認証)
		h.logger.Info("Received event", zap.String("id", event.ID), zap.String("type", event.Type))
	case
		"payment_method.attached", // 決済手段関連付け
		"customer.updated",        // 顧客更新
		"customer.deleted":        // 顧客削除
		h.logger.Debug("Received event", zap.String("id", event.ID), zap.String("type", event.Type))
	default:
		// エラーを返すと常にリトライが走ってしまうため、ログ出力のみして成功ステータスを返す
		h.logger.Warn("Received not implemented event", zap.String("type", event.Type))
	}
	return nil
}
