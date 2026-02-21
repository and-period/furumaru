package worker

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSendReport(t *testing.T) {
	t.Parallel()

	template := &entity.ReportTemplate{
		TemplateID: entity.ReportTemplateIDReceivedContact,
		Template:   `{"type":"bubble","body":{"type":"box","contents":[{"type":"text","text":"{{.Overview}}"}]}}`,
		CreatedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
		UpdatedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		payload   *entity.WorkerPayload
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReportTemplate.EXPECT().Get(ctx, entity.ReportTemplateIDReceivedContact).Return(template, nil)
				mocks.line.EXPECT().PushMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, messages ...messaging_api.MessageInterface) error {
						require.Len(t, messages, 1)
						msg, ok := messages[0].(*messaging_api.FlexMessage)
						require.True(t, ok)
						assert.Equal(t, "[ふるマル] received-contact", msg.AltText)
						return nil
					})
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeReceivedContact,
				Report: &entity.ReportConfig{
					TemplateID: entity.ReportTemplateIDReceivedContact,
					Overview:   "お問い合わせ件名",
					Link:       "htts://admin.and-period.jp/contacts/contact-id",
				},
			},
			expectErr: nil,
		},
		{
			name: "failed to get report template",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReportTemplate.EXPECT().Get(ctx, entity.ReportTemplateIDReceivedContact).Return(nil, assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeReceivedContact,
				Report: &entity.ReportConfig{
					TemplateID: entity.ReportTemplateIDReceivedContact,
					Overview:   "お問い合わせ件名",
					Link:       "htts://admin.and-period.jp/contacts/contact-id",
				},
			},
			expectErr: assert.AnError,
		},
		{
			name: "failed to push line message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReportTemplate.EXPECT().Get(ctx, entity.ReportTemplateIDReceivedContact).Return(template, nil)
				mocks.line.EXPECT().PushMessage(ctx, gomock.Any()).Return(assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeReceivedContact,
				Report: &entity.ReportConfig{
					TemplateID: entity.ReportTemplateIDReceivedContact,
					Overview:   "お問い合わせ件名",
					Link:       "htts://admin.and-period.jp/contacts/contact-id",
				},
			},
			expectErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.sendReport(ctx, tt.payload)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
