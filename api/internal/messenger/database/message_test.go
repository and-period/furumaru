package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	assert.NotNil(t, NewMessage(nil))
}

func TestMessage_MultiCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 6, 26, 19, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	type args struct {
		messages entity.Messages
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				messages: entity.Messages{
					testMessage("message-id", now()),
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				msg := testMessage("message-id", now())
				err = m.db.DB.Create(&msg).Error
				require.NoError(t, err)
			},
			args: args{
				messages: entity.Messages{
					testMessage("message-id", now()),
				},
			},
			want: want{
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := m.dbDelete(ctx, notificationTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &message{db: m.db, now: now}
			err = db.MultiCreate(ctx, tt.args.messages)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testMessage(id string, now time.Time) *entity.Message {
	return &entity.Message{
		ID:         id,
		UserType:   entity.UserTypeUser,
		UserID:     "user-id",
		Type:       entity.MessageTypeNotification,
		Title:      "メッセージ件名",
		Body:       "メッセージ内容です。",
		Link:       "https://and-period.jp",
		Read:       false,
		ReceivedAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
