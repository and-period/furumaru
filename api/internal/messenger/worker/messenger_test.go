package worker

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessenger(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 7, 21, 18, 30, 0, 0)
	template := &entity.MessageTemplate{
		TemplateID:    entity.MessageIDNotification,
		TitleTemplate: "件名: {{.Title}}",
		BodyTemplate:  "テンプレートです。",
		CreatedAt:     jst.Date(2022, 7, 21, 18, 30, 0, 0),
		UpdatedAt:     jst.Date(2022, 7, 21, 18, 30, 0, 0),
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
				mocks.db.MessageTemplate.EXPECT().Get(ctx, entity.MessageIDNotification).Return(template, nil)
				mocks.db.Message.EXPECT().
					MultiCreate(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, messages entity.Messages) error {
						require.Len(t, messages, 1)
						expect := entity.Messages{
							{
								ID:         messages[0].ID, // ignore
								UserType:   entity.UserTypeUser,
								UserID:     "user-id",
								Type:       entity.MessageTypeNotification,
								Title:      "件名: メッセージタイトル",
								Body:       "テンプレートです。",
								Link:       "https://and-period.jp",
								ReceivedAt: now,
							},
						}
						assert.Equal(t, expect, messages)
						return nil
					})
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeNotification,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Message: &entity.MessageConfig{
					MessageID:   entity.MessageIDNotification,
					MessageType: entity.MessageTypeNotification,
					Title:       "メッセージタイトル",
					Link:        "https://and-period.jp",
					ReceivedAt:  now,
				},
			},
			expectErr: nil,
		},
		{
			name: "failed to get message template",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.MessageTemplate.EXPECT().Get(ctx, entity.MessageIDNotification).Return(nil, errmock)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeNotification,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Message: &entity.MessageConfig{
					MessageID:   entity.MessageIDNotification,
					MessageType: entity.MessageTypeNotification,
					Title:       "メッセージタイトル",
					Link:        "https://and-period.jp",
					ReceivedAt:  now,
				},
			},
			expectErr: errmock,
		},
		{
			name: "failed to multi create messages",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.MessageTemplate.EXPECT().Get(ctx, entity.MessageIDNotification).Return(template, nil)
				mocks.db.Message.EXPECT().MultiCreate(ctx, gomock.Any()).Return(errmock)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeNotification,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Message: &entity.MessageConfig{
					MessageID:   entity.MessageIDNotification,
					MessageType: entity.MessageTypeNotification,
					Title:       "メッセージタイトル",
					Link:        "https://and-period.jp",
					ReceivedAt:  now,
				},
			},
			expectErr: errmock,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.messenger(ctx, tt.payload)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
