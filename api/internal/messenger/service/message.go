package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/backoff"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListMessages(ctx context.Context, in *messenger.ListMessagesInput) (entity.Messages, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	params := &database.ListMessagesParams{
		Limit:    int(in.Limit),
		Offset:   int(in.Offset),
		UserType: in.UserType,
		UserID:   in.UserID,
	}
	var (
		messages entity.Messages
		total    int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		messages, err = s.db.Message.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Message.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	return messages, total, nil
}

func (s *service) GetMessage(ctx context.Context, in *messenger.GetMessageInput) (*entity.Message, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	message, err := s.db.Message.Get(ctx, in.MessageID)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	if message.Read {
		return message, nil
	}
	s.waitGroup.Add(1)
	go func(message *entity.Message) {
		defer s.waitGroup.Done()
		if err := s.updateMessageRead(context.Background(), message.ID); err != nil {
			s.logger.Error("Failed to update message read", zap.String("messageId", message.ID), zap.Error(err))
		}
	}(message)
	message.Read = true
	return message, nil
}

func (s *service) updateMessageRead(ctx context.Context, messageID string) error {
	updateFn := func() error {
		return s.db.Message.UpdateRead(ctx, messageID)
	}
	const maxRetries = 3
	retry := backoff.NewExponentialBackoff(maxRetries)
	return backoff.Retry(ctx, retry, updateFn, backoff.WithRetryablel(exception.Retryable))
}
