package worker

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

func (w *worker) reporter(ctx context.Context, payload *entity.WorkerPayload) error {
	template, err := w.db.ReportTemplate.Get(ctx, payload.Report.ReportID)
	if err != nil {
		return err
	}
	msg, err := template.Build(payload.Report.Fields())
	if err != nil {
		return err
	}
	sendFn := func() error {
		w.logger.Debug("Send report",
			zap.String("reportId", payload.Report.ReportID), zap.String("message", msg))
		return w.line.PushMessage(ctx, linebot.NewTextMessage(msg))
	}
	retry := backoff.NewExponentialBackoff(w.maxRetries)
	return backoff.Retry(ctx, retry, sendFn, backoff.WithRetryablel(exception.Retryable))
}
