package service

import (
	"context"
	"errors"

	"github.com/and-period/marche/api/internal/messenger/entity"
	"github.com/and-period/marche/api/pkg/backoff"
	"github.com/and-period/marche/api/pkg/mailer"
)

func (s *messengerService) sendInfoMail(
	ctx context.Context, msg *entity.MailConfig, ps ...*mailer.Personalization,
) error {
	sendFn := func() error {
		return s.mailer.MultiSendFromInfo(ctx, msg.EmailID, ps)
	}
	retry := backoff.NewExponentialBackoff(s.maxRetries)
	return backoff.Retry(ctx, retry, sendFn, backoff.WithRetryablel(s.retryableSendMail))
}

func (s *messengerService) retryableSendMail(err error) bool {
	return errors.Is(err, mailer.ErrTimeout) ||
		errors.Is(err, mailer.ErrUnavailable) ||
		errors.Is(err, mailer.ErrInternal)
}
