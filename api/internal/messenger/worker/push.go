package worker

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/and-period/furumaru/api/pkg/firebase/messaging"
	"go.uber.org/zap"
)

func (w *worker) multiSendPush(ctx context.Context, payload *entity.WorkerPayload) error {
	tokens, err := w.fetchTokens(ctx, payload)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		w.logger.Debug("Tokens is empty", zap.String("pushId", payload.Push.PushID))
		return nil
	}
	template, err := w.db.PushTemplate.Get(ctx, payload.Push.PushID)
	if err != nil {
		return err
	}
	title, body, err := template.Build(payload.Push.Data)
	if err != nil {
		return err
	}
	msg := &messaging.Message{
		Title:    title,
		Body:     body,
		ImageURL: template.ImageURL,
		Data:     payload.Push.Data,
	}
	w.logger.Debug("Send push", zap.String("pushId", payload.Push.PushID), zap.Any("message", msg))
	sendFn := func() error {
		_, _, err := w.messaging.MultiSend(ctx, msg, tokens...)
		return exception.InternalError(err)
	}
	retry := backoff.NewExponentialBackoff(w.maxRetries)
	return backoff.Retry(ctx, retry, sendFn, backoff.WithRetryablel(exception.Retryable))
}

func (w *worker) fetchTokens(ctx context.Context, payload *entity.WorkerPayload) ([]string, error) {
	switch payload.UserType {
	case entity.UserTypeAdmin,
		entity.UserTypeAdministrator,
		entity.UserTypeCoordinator,
		entity.UserTypeProducer:
		return w.fetchAdminTokens(ctx, payload.UserIDs)
	case entity.UserTypeUser:
		return w.fetchUserTokens(ctx, payload.UserIDs)
	default:
		return nil, fmt.Errorf("worker: failed to multi send push: %w", errUnknownUserType)
	}
}

func (w *worker) fetchAdminTokens(ctx context.Context, adminIDs []string) ([]string, error) {
	in := &user.MultiGetAdminDevicesInput{
		AdminIDs: adminIDs,
	}
	return w.user.MultiGetAdminDevices(ctx, in)
}

func (w *worker) fetchUserTokens(ctx context.Context, userIDs []string) ([]string, error) {
	in := &user.MultiGetUserDevicesInput{
		UserIDs: userIDs,
	}
	return w.user.MultiGetUserDevices(ctx, in)
}
