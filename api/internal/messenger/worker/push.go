package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/and-period/furumaru/api/pkg/firebase/messaging"
)

func (w *worker) multiSendPush(ctx context.Context, payload *entity.WorkerPayload) error {
	tokens, err := w.fetchTokens(ctx, payload)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		slog.Debug("Tokens is empty", slog.String("templateId", string(payload.Push.TemplateID)))
		return nil
	}
	template, err := w.db.PushTemplate.Get(ctx, payload.Push.TemplateID)
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
	slog.Debug("Send push", slog.String("templateId", string(payload.Push.TemplateID)), slog.Any("message", msg))
	sendFn := func() error {
		return w.sendMessaing(ctx, msg, payload.UserType, tokens)
	}
	retry := backoff.NewExponentialBackoff(w.maxRetries)
	return backoff.Retry(ctx, retry, sendFn, backoff.WithRetryablel(w.isRetryable))
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

func (w *worker) sendMessaing(ctx context.Context, msg *messaging.Message, userType entity.UserType, tokens []string) error {
	switch userType {
	case entity.UserTypeAdmin,
		entity.UserTypeAdministrator,
		entity.UserTypeCoordinator,
		entity.UserTypeProducer:
		_, _, err := w.adminMessaging.MultiSend(ctx, msg, tokens...)
		return err
	case entity.UserTypeUser:
		_, _, err := w.userMessagging.MultiSend(ctx, msg, tokens...)
		return err
	default:
		return fmt.Errorf("worker: failed to multi send push: %w", errUnknownUserType)
	}
}
