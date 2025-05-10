package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListMessages(t *testing.T) {
	t.Parallel()

	now := jst.Date(20222, 7, 7, 18, 30, 0, 0)
	params := &database.ListMessagesParams{
		Limit:    20,
		Offset:   0,
		UserType: entity.UserTypeUser,
		UserID:   "user-id",
		Orders: []*database.ListMessagesOrder{
			{Key: entity.MessageOrderByReceivedAt, OrderByASC: false},
		},
	}
	messages := entity.Messages{
		{
			ID:         "message-id",
			UserType:   entity.UserTypeUser,
			UserID:     "user-id",
			Type:       entity.MessageTypeNotification,
			Title:      "メッセージタイトル",
			Body:       "メッセージの内容です。",
			Link:       "https://and-period.jp",
			Read:       false,
			ReceivedAt: now,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *messenger.ListMessagesInput
		expect      entity.Messages
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Message.EXPECT().List(gomock.Any(), params).Return(messages, nil)
				mocks.db.Message.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &messenger.ListMessagesInput{
				Limit:    20,
				Offset:   0,
				UserType: entity.UserTypeUser,
				UserID:   "user-id",
				Orders: []*messenger.ListMessagesOrder{
					{Key: entity.MessageOrderByReceivedAt, OrderByASC: false},
				},
			},
			expect:      messages,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.ListMessagesInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list messages",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Message.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Message.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &messenger.ListMessagesInput{
				Limit:    20,
				Offset:   0,
				UserType: entity.UserTypeUser,
				UserID:   "user-id",
				Orders: []*messenger.ListMessagesOrder{
					{Key: entity.MessageOrderByReceivedAt, OrderByASC: false},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count messagses",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Message.EXPECT().List(gomock.Any(), params).Return(messages, nil)
				mocks.db.Message.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &messenger.ListMessagesInput{
				Limit:    20,
				Offset:   0,
				UserType: entity.UserTypeUser,
				UserID:   "user-id",
				Orders: []*messenger.ListMessagesOrder{
					{Key: entity.MessageOrderByReceivedAt, OrderByASC: false},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListMessages(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestGetMessage(t *testing.T) {
	t.Parallel()

	now := jst.Date(20222, 7, 7, 18, 30, 0, 0)
	message := &entity.Message{
		ID:         "message-id",
		UserType:   entity.UserTypeUser,
		UserID:     "user-id",
		Type:       entity.MessageTypeNotification,
		Title:      "メッセージタイトル",
		Body:       "メッセージの内容です。",
		Link:       "https://and-period.jp",
		Read:       true,
		ReceivedAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.GetMessageInput
		expect    *entity.Message
		expectErr error
	}{
		{
			name: "success from system",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Message.EXPECT().Get(ctx, "message-id").Return(message, nil)
			},
			input: &messenger.GetMessageInput{
				MessageID: "message-id",
				UserType:  entity.UserTypeNone,
				UserID:    "",
			},
			expect:    message,
			expectErr: nil,
		},
		{
			name: "success first read",
			setup: func(ctx context.Context, mocks *mocks) {
				message := &entity.Message{
					ID:       "message-id",
					UserType: entity.UserTypeUser,
					UserID:   "user-id",
					Read:     false,
				}
				mocks.db.Message.EXPECT().Get(ctx, "message-id").Return(message, nil)
				mocks.db.Message.EXPECT().UpdateRead(gomock.Any(), "message-id").Return(nil)
			},
			input: &messenger.GetMessageInput{
				MessageID: "message-id",
				UserType:  entity.UserTypeUser,
				UserID:    "user-id",
			},
			expect: &entity.Message{
				ID:       "message-id",
				UserType: entity.UserTypeUser,
				UserID:   "user-id",
				Read:     true,
			},
			expectErr: nil,
		},
		{
			name: "success second after read",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Message.EXPECT().Get(ctx, "message-id").Return(message, nil)
			},
			input: &messenger.GetMessageInput{
				MessageID: "message-id",
				UserType:  entity.UserTypeUser,
				UserID:    "user-id",
			},
			expect:    message,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.GetMessageInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Message.EXPECT().Get(ctx, "message-id").Return(nil, assert.AnError)
			},
			input: &messenger.GetMessageInput{
				MessageID: "message-id",
				UserType:  entity.UserTypeUser,
				UserID:    "user-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get someone else",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Message.EXPECT().Get(ctx, "message-id").Return(message, nil)
			},
			input: &messenger.GetMessageInput{
				MessageID: "message-id",
				UserType:  entity.UserTypeUser,
				UserID:    "other-id",
			},
			expect:    nil,
			expectErr: exception.ErrForbidden,
		},
		{
			name: "failed to update read",
			setup: func(ctx context.Context, mocks *mocks) {
				message := &entity.Message{
					ID:       "message-id",
					UserType: entity.UserTypeUser,
					UserID:   "user-id",
					Read:     false,
				}
				mocks.db.Message.EXPECT().Get(ctx, "message-id").Return(message, nil)
				mocks.db.Message.EXPECT().UpdateRead(gomock.Any(), "message-id").Return(assert.AnError)
			},
			input: &messenger.GetMessageInput{
				MessageID: "message-id",
				UserType:  entity.UserTypeUser,
				UserID:    "user-id",
			},
			expect: &entity.Message{
				ID:       "message-id",
				UserType: entity.UserTypeUser,
				UserID:   "user-id",
				Read:     true,
			},
			expectErr: nil,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetMessage(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}
