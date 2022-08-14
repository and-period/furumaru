package worker

import (
	"context"
	"fmt"

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
	altText := fmt.Sprintf("[ふるマル] %s", payload.Report.ReportID)
	container, err := template.Build(payload.Report.Fields())
	if err != nil {
		return err
	}
	w.logger.Debug("Send report", zap.String("reportId", payload.Report.ReportID), zap.Any("message", container))
	sendFn := func() error {
		err := w.line.PushMessage(ctx, linebot.NewFlexMessage(altText, container))
		return exception.InternalError(err)
	}
	retry := backoff.NewExponentialBackoff(w.maxRetries)
	return backoff.Retry(ctx, retry, sendFn, backoff.WithRetryablel(exception.Retryable))
}
