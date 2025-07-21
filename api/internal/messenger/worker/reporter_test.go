package worker

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendReport(t *testing.T) {
	t.Parallel()

	template := &entity.ReportTemplate{
		TemplateID: entity.ReportTemplateIDReceivedContact,
		Template:   `{"type":"bubble","body":{"type":"box","contents":[{"type":"text","text":"{{.Overview}}"}]}}`,
		CreatedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
		UpdatedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
	}
	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type: linebot.FlexComponentTypeBox,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "お問い合わせ件名",
				},
			},
		},
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
				mocks.db.ReportTemplate.EXPECT().
					Get(ctx, entity.ReportTemplateIDReceivedContact).
					Return(template, nil)
				mocks.line.EXPECT().PushMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, messages ...linebot.SendingMessage) error {
						require.Len(t, messages, 1)
						msg, ok := messages[0].(*linebot.FlexMessage)
						require.True(t, ok)
						assert.Equal(t, "[ふるマル] received-contact", msg.AltText)
						assert.Equal(t, container, msg.Contents)
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
				mocks.db.ReportTemplate.EXPECT().
					Get(ctx, entity.ReportTemplateIDReceivedContact).
					Return(nil, assert.AnError)
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
				mocks.db.ReportTemplate.EXPECT().
					Get(ctx, entity.ReportTemplateIDReceivedContact).
					Return(template, nil)
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
		t.Run(
			tt.name,
			testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
				err := worker.sendReport(ctx, tt.payload)
				assert.ErrorIs(t, err, tt.expectErr)
			}),
		)
	}
}
