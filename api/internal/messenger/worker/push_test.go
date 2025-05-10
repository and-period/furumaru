package worker

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/firebase/messaging"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMultiSendPush(t *testing.T) {
	t.Parallel()

	in := &user.MultiGetAdminDevicesInput{
		AdminIDs: []string{"admin-id"},
	}
	devices := []string{"instance-id"}
	template := &entity.PushTemplate{
		TemplateID:    entity.PushTemplateIDContact,
		TitleTemplate: "件名: {{.Title}}",
		BodyTemplate:  "テンプレートです。",
		ImageURL:      "https://and-period.jp/image.png",
		CreatedAt:     jst.Date(2022, 7, 21, 18, 30, 0, 0),
		UpdatedAt:     jst.Date(2022, 7, 21, 18, 30, 0, 0),
	}
	message := &messaging.Message{
		Title:    "件名: テストお問い合わせ",
		Body:     "テンプレートです。",
		ImageURL: "https://and-period.jp/image.png",
		Data:     map[string]string{"Title": "テストお問い合わせ"},
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
				mocks.user.EXPECT().MultiGetAdminDevices(ctx, in).Return(devices, nil)
				mocks.db.PushTemplate.EXPECT().Get(ctx, entity.PushTemplateIDContact).Return(template, nil)
				mocks.messaging.EXPECT().MultiSend(gomock.Any(), message, devices).Return(int64(1), int64(0), nil)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeReceivedContact,
				UserType:  entity.UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Push: &entity.PushConfig{
					TemplateID: entity.PushTemplateIDContact,
					Data:       map[string]string{"Title": "テストお問い合わせ"},
				},
			},
			expectErr: nil,
		},
		{
			name: "success token empty",
			setup: func(ctx context.Context, mocks *mocks) {
				devices := []string{}
				mocks.user.EXPECT().MultiGetAdminDevices(ctx, in).Return(devices, nil)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeReceivedContact,
				UserType:  entity.UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Push: &entity.PushConfig{
					TemplateID: entity.PushTemplateIDContact,
					Data:       map[string]string{"Title": "テストお問い合わせ"},
				},
			},
			expectErr: nil,
		},
		{
			name: "failed to get tokens",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetAdminDevices(ctx, in).Return(nil, assert.AnError)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeReceivedContact,
				UserType:  entity.UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Push: &entity.PushConfig{
					TemplateID: entity.PushTemplateIDContact,
					Data:       map[string]string{"Title": "テストお問い合わせ"},
				},
			},
			expectErr: assert.AnError,
		},
		{
			name: "failed to get push template",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetAdminDevices(ctx, in).Return(devices, nil)
				mocks.db.PushTemplate.EXPECT().Get(ctx, entity.PushTemplateIDContact).Return(nil, assert.AnError)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeReceivedContact,
				UserType:  entity.UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Push: &entity.PushConfig{
					TemplateID: entity.PushTemplateIDContact,
					Data:       map[string]string{"Title": "テストお問い合わせ"},
				},
			},
			expectErr: assert.AnError,
		},
		{
			name: "failed to multi send push",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetAdminDevices(ctx, in).Return(devices, nil)
				mocks.db.PushTemplate.EXPECT().Get(ctx, entity.PushTemplateIDContact).Return(template, nil)
				mocks.messaging.EXPECT().MultiSend(gomock.Any(), message, devices).Return(int64(0), int64(0), assert.AnError)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeReceivedContact,
				UserType:  entity.UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Push: &entity.PushConfig{
					TemplateID: entity.PushTemplateIDContact,
					Data:       map[string]string{"Title": "テストお問い合わせ"},
				},
			},
			expectErr: assert.AnError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.multiSendPush(ctx, tt.payload)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestFetchTokens(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		payload   *entity.WorkerPayload
		expect    []string
		expectErr error
	}{
		{
			name: "success admins",
			setup: func(ctx context.Context, mocks *mocks) {
				in := &user.MultiGetAdminDevicesInput{AdminIDs: []string{"admin-id"}}
				devices := []string{"instance-id"}
				mocks.user.EXPECT().MultiGetAdminDevices(ctx, in).Return(devices, nil)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeReceivedContact,
				UserType:  entity.UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Push: &entity.PushConfig{
					TemplateID: entity.PushTemplateIDContact,
					Data:       map[string]string{"key": "value"},
				},
			},
			expect:    []string{"instance-id"},
			expectErr: nil,
		},
		{
			name: "success users",
			setup: func(ctx context.Context, mocks *mocks) {
				in := &user.MultiGetUserDevicesInput{UserIDs: []string{"user-id"}}
				devices := []string{"instance-id"}
				mocks.user.EXPECT().MultiGetUserDevices(ctx, in).Return(devices, nil)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeReceivedContact,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Push: &entity.PushConfig{
					TemplateID: entity.PushTemplateIDContact,
					Data:       map[string]string{"key": "value"},
				},
			},
			expect:    []string{"instance-id"},
			expectErr: nil,
		},
		{
			name: "failed to multi get admin devices",
			setup: func(ctx context.Context, mocks *mocks) {
				in := &user.MultiGetAdminDevicesInput{AdminIDs: []string{"admin-id"}}
				mocks.user.EXPECT().MultiGetAdminDevices(ctx, in).Return(nil, assert.AnError)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeReceivedContact,
				UserType:  entity.UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Push: &entity.PushConfig{
					TemplateID: entity.PushTemplateIDContact,
					Data:       map[string]string{"key": "value"},
				},
			},
			expect:    nil,
			expectErr: assert.AnError,
		},
		{
			name: "failed to multi get user devices",
			setup: func(ctx context.Context, mocks *mocks) {
				in := &user.MultiGetUserDevicesInput{UserIDs: []string{"user-id"}}
				mocks.user.EXPECT().MultiGetUserDevices(ctx, in).Return(nil, assert.AnError)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeReceivedContact,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Push: &entity.PushConfig{
					TemplateID: entity.PushTemplateIDContact,
					Data:       map[string]string{"key": "value"},
				},
			},
			expect:    nil,
			expectErr: assert.AnError,
		},
		{
			name:  "failed to unknown user type",
			setup: func(ctx context.Context, mocks *mocks) {},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeReceivedContact,
				UserType:  entity.UserTypeGuest,
				UserIDs:   []string{"user-id"},
				Push: &entity.PushConfig{
					TemplateID: entity.PushTemplateIDContact,
					Data:       map[string]string{"key": "value"},
				},
			},
			expect:    nil,
			expectErr: errUnknownUserType,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			actual, err := worker.fetchTokens(ctx, tt.payload)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}
