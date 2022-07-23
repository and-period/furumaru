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

func TestMessage_List(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, messageTable)
	messages := make(entity.Messages, 3)
	messages[0] = testMessage("message-id01", now().Add(-time.Hour))
	messages[1] = testMessage("message-id02", now())
	messages[2] = testMessage("message-id03", now().Add(time.Hour))
	err = m.db.DB.Create(&messages).Error
	require.NoError(t, err)

	type args struct {
		params *ListMessagesParams
	}
	type want struct {
		messages entity.Messages
		hasErr   bool
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
				params: &ListMessagesParams{
					UserType: entity.UserTypeUser,
					UserID:   "user-id",
					Limit:    20,
					Offset:   1,
				},
			},
			want: want{
				messages: messages[0:2],
				hasErr:   false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, m)

			db := &message{db: m.db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			fillIgnoreMessagesField(actual, now())
			assert.ElementsMatch(t, tt.want.messages, actual)
		})
	}
}

func TestMessage_Count(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, messageTable)
	messages := make(entity.Messages, 3)
	messages[0] = testMessage("message-id01", now())
	messages[1] = testMessage("message-id02", now())
	messages[2] = testMessage("message-id03", now())
	err = m.db.DB.Create(&messages).Error
	require.NoError(t, err)

	type args struct {
		params *ListMessagesParams
	}
	type want struct {
		total  int64
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
				params: &ListMessagesParams{
					UserType: entity.UserTypeUser,
					UserID:   "user-id",
					Limit:    20,
					Offset:   1,
				},
			},
			want: want{
				total:  3,
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, m)

			db := &message{db: m.db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestMessage_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, messageTable)
	msg := testMessage("message-id", now())
	err = m.db.DB.Create(&msg).Error
	require.NoError(t, err)

	type args struct {
		messageID string
	}
	type want struct {
		message *entity.Message
		hasErr  bool
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
				messageID: "message-id",
			},
			want: want{
				message: msg,
				hasErr:  false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				messageID: "other-id",
			},
			want: want{
				message: nil,
				hasErr:  true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, m)

			db := &message{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.messageID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreMessageField(actual, now())
			assert.Equal(t, tt.want.message, actual)
		})
	}
}

func TestMessage_MultiCreate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 6, 26, 19, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, messageTable)

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

			err := m.dbDelete(ctx, messageTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &message{db: m.db, now: now}
			err = db.MultiCreate(ctx, tt.args.messages)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestMessage_UpdateRead(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 6, 26, 19, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, messageTable)

	type args struct {
		messageID string
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
			name: "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				msg := testMessage("message-id", now())
				msg.Read = false
				err := m.db.DB.Create(&msg).Error
				require.NoError(t, err)
			},
			args: args{
				messageID: "message-id",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				messageID: "message-id",
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "already updated",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				msg := testMessage("message-id", now())
				msg.Read = true
				err := m.db.DB.Create(&msg).Error
				require.NoError(t, err)
			},
			args: args{
				messageID: "message-id",
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

			err := m.dbDelete(ctx, messageTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &message{db: m.db, now: now}
			err = db.UpdateRead(ctx, tt.args.messageID)
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

func fillIgnoreMessageField(m *entity.Message, now time.Time) {
	if m == nil {
		return
	}
	m.ReceivedAt = m.ReceivedAt.In(now.Location())
	m.CreatedAt = m.ReceivedAt.In(now.Location())
	m.UpdatedAt = m.ReceivedAt.In(now.Location())
}

func fillIgnoreMessagesField(ms entity.Messages, now time.Time) {
	for i := range ms {
		fillIgnoreMessageField(ms[i], now)
	}
}
