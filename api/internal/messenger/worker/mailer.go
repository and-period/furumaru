package worker

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"go.uber.org/zap"
)

func (w *worker) sendMail(ctx context.Context, emailID string, ps ...*mailer.Personalization) error {
	if len(ps) == 0 {
		w.logger.Debug("Personalizations is empty", zap.String("emailId", emailID))
		return nil
	}
	sendFn := func() error {
		w.logger.Debug("Send email", zap.String("emailId", emailID), zap.Any("personalizations", ps))
		if err := w.mailer.MultiSendFromInfo(ctx, emailID, ps); err != nil {
			return exception.InternalError(err)
		}
		return nil
	}
	retry := backoff.NewExponentialBackoff(w.maxRetries)
	return backoff.Retry(ctx, retry, sendFn, backoff.WithRetryablel(exception.Retryable))
}
