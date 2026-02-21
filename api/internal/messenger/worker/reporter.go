package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
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
	slog.Debug("Send report", slog.String("templateId", string(payload.Report.TemplateID)), slog.Any("message", container))
	msg := &messaging_api.FlexMessage{
		AltText:  altText,
		Contents: container,
	}
	sendFn := func() error {
		return w.line.PushMessage(ctx, msg)
	}
	retry := backoff.NewExponentialBackoff(w.maxRetries)
	return backoff.Retry(ctx, retry, sendFn, backoff.WithRetryablel(w.isRetryable))
}
