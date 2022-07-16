package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotifyRegisterAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.NotifyRegisterAdminInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, queue *entity.ReceivedQueue) error {
						expect := &entity.ReceivedQueue{
							ID:        queue.ID, // ignore
							EventType: entity.EventTypeAdminRegister,
							UserType:  entity.UserTypeAdmin,
							UserIDs:   []string{"admin-id"},
							Done:      false,
						}
						assert.Equal(t, expect, queue)
						return nil
					})
				mocks.producer.EXPECT().
					SendMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, b []byte) (string, error) {
						payload := &entity.WorkerPayload{}
						err := json.Unmarshal(b, payload)
						require.NoError(t, err)
						expect := &entity.WorkerPayload{
							QueueID:   payload.QueueID, // ignore
							EventType: entity.EventTypeAdminRegister,
							UserType:  entity.UserTypeAdmin,
							UserIDs:   []string{"admin-id"},
							Email: &entity.MailConfig{
								EmailID: entity.EmailIDAdminRegister,
								Substitutions: map[string]string{
									"サイトURL": "htts://admin.and-period.jp/signin",
									"パスワード":  "!Qaz2wsx",
								},
							},
						}
						assert.Equal(t, expect, payload)
						return "message-id", nil
					})
			},
			input: &messenger.NotifyRegisterAdminInput{
				AdminID:  "admin-id",
				Password: "!Qaz2wsx",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.NotifyRegisterAdminInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", errmock)
			},
			input: &messenger.NotifyRegisterAdminInput{
				AdminID:  "admin-id",
				Password: "!Qaz2wsx",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.NotifyRegisterAdmin(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestNotifyReceivedContact(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 7, 7, 18, 30, 0, 0)
	contact := &entity.Contact{
		ID:          "contact-id",
		Title:       "お問い合わせ件名",
		Content:     "お問い合わせ内容です。",
		Username:    "あんど どっと",
		Email:       "test-user@and-period.jp",
		PhoneNumber: "+819012345678",
		Status:      entity.ContactStatusInprogress,
		Priority:    entity.ContactPriorityMiddle,
		Note:        "対応者のメモです",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.NotifyReceivedContactInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Contact.EXPECT().Get(ctx, "contact-id").Return(contact, nil)
				mocks.db.ReceivedQueue.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, queue *entity.ReceivedQueue) error {
						expect := &entity.ReceivedQueue{
							ID:        queue.ID, // ignore
							EventType: entity.EventTypeUserReceivedContact,
							UserType:  entity.UserTypeGuest,
							UserIDs:   nil,
							Done:      false,
						}
						assert.Equal(t, expect, queue)
						return nil
					})
				mocks.producer.EXPECT().
					SendMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, b []byte) (string, error) {
						payload := &entity.WorkerPayload{}
						err := json.Unmarshal(b, payload)
						require.NoError(t, err)
						assert.Equal(t, now.Unix(), payload.Report.ReceivedAt.Unix())
						expect := &entity.WorkerPayload{
							QueueID:   payload.QueueID, // ignore
							EventType: entity.EventTypeUserReceivedContact,
							UserType:  entity.UserTypeGuest,
							UserIDs:   nil,
							Guest: &entity.Guest{
								Name:  "&. スタッフ",
								Email: "test-user@and-period.jp",
							},
							Email: &entity.MailConfig{
								EmailID: entity.EmailIDUserReceivedContact,
								Substitutions: map[string]string{
									"氏名":      "&. スタッフ",
									"メールアドレス": "test-user@and-period.jp",
									"件名":      "お問い合わせ件名",
									"本文":      "お問い合わせ内容です。",
								},
							},
							Report: &entity.Report{
								ReportID:   entity.ReportIDReceivedContact,
								Overview:   "お問い合わせ件名",
								Detail:     "お問い合わせ内容です。",
								Link:       "htts://admin.and-period.jp/contacts/contact-id",
								ReceivedAt: payload.Report.ReceivedAt,
							},
						}
						assert.Equal(t, expect, payload)
						return "message-id", nil
					})
			},
			input: &messenger.NotifyReceivedContactInput{
				ContactID: "contact-id",
				Username:  "&. スタッフ",
				Email:     "test-user@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.NotifyReceivedContactInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get contact",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Contact.EXPECT().Get(ctx, "contact-id").Return(nil, errmock)
			},
			input: &messenger.NotifyReceivedContactInput{
				ContactID: "contact-id",
				Username:  "&. スタッフ",
				Email:     "test-user@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Contact.EXPECT().Get(ctx, "contact-id").Return(contact, nil)
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", errmock)
			},
			input: &messenger.NotifyReceivedContactInput{
				ContactID: "contact-id",
				Username:  "&. スタッフ",
				Email:     "test-user@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.NotifyReceivedContact(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestSendMessaget(t *testing.T) {
	t.Parallel()
	queue := &entity.ReceivedQueue{
		ID:        "queue-id",
		EventType: entity.EventTypeUnknown,
		UserType:  entity.UserTypeUser,
		UserIDs:   []string{"user-id"},
	}
	tests := []struct {
		name    string
		setup   func(ctx context.Context, mocks *mocks)
		payload *entity.WorkerPayload
		hasErr  bool
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, queue).Return(nil)
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email:     &entity.MailConfig{},
			},
			hasErr: false,
		},
		{
			name: "failed to create received queue",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, queue).Return(errmock)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email:     &entity.MailConfig{},
			},
			hasErr: true,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, queue).Return(nil)
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", errmock)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email:     &entity.MailConfig{},
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.sendMessage(ctx, tt.payload)
			assert.Equal(t, tt.hasErr, err != nil, err)
		}))
	}
}
