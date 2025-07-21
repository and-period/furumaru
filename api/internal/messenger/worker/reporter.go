package worker

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

func (w *worker) sendReport(ctx context.Context, payload *entity.WorkerPayload) error {
	template, err := w.db.ReportTemplate.Get(ctx, payload.Report.TemplateID)
	if err != nil {
		return err
	}
	altText := fmt.Sprintf("[ふるマル] %s", payload.Report.TemplateID)
	container, err := template.Build(payload.Report.Fields())
	if err != nil {
		return err
	}
	w.logger.Debug(
		"Send report",
		zap.String("templateId", string(payload.Report.TemplateID)),
		zap.Any("message", container),
	)
	sendFn := func() error {
		return w.line.PushMessage(ctx, linebot.NewFlexMessage(altText, container))
	}
	retry := backoff.NewExponentialBackoff(w.maxRetries)
	return backoff.Retry(ctx, retry, sendFn, backoff.WithRetryablel(w.isRetryable))
}
