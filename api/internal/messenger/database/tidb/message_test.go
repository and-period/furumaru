package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	assert.NotNil(t, NewMessage(nil))
}

func TestMessage_List(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	messages := make(entity.Messages, 3)
	messages[0] = testMessage("message-id01", now().Add(-time.Hour))
	messages[1] = testMessage("message-id02", now())
	messages[2] = testMessage("message-id03", now().Add(time.Hour))
	err = db.DB.Create(&messages).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListMessagesParams
	}
	type want struct {
		messages entity.Messages
		err      error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListMessagesParams{
					UserType: entity.UserTypeUser,
					UserID:   "user-id",
					Limit:    20,
					Offset:   1,
				},
			},
			want: want{
				messages: messages[1:],
				err:      nil,
			},
		},
		{
			name:  "success with sort",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListMessagesParams{
					Orders: []*database.ListMessagesOrder{
						{Key: entity.MessageOrderByReceivedAt, OrderByASC: false},
					},
				},
			},
			want: want{
				messages: messages,
				err:      nil,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &message{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.ElementsMatch(t, tt.want.messages, actual)
		})
	}
}

func TestMessage_Count(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	messages := make(entity.Messages, 3)
	messages[0] = testMessage("message-id01", now())
	messages[1] = testMessage("message-id02", now())
	messages[2] = testMessage("message-id03", now())
	err = db.DB.Create(&messages).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListMessagesParams
	}
	type want struct {
		total int64
		err   error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListMessagesParams{
					UserType: entity.UserTypeUser,
					UserID:   "user-id",
					Limit:    20,
					Offset:   1,
				},
			},
			want: want{
				total: 3,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &message{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestMessage_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	msg := testMessage("message-id", now())
	err = db.DB.Create(&msg).Error
	require.NoError(t, err)

	type args struct {
		messageID string
	}
	type want struct {
		message *entity.Message
		err     error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				messageID: "message-id",
			},
			want: want{
				message: msg,
				err:     nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				messageID: "other-id",
			},
			want: want{
				message: nil,
				err:     database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &message{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.messageID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.message, actual)
		})
	}
}

func TestMessage_MultiCreate(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	type args struct {
		messages entity.Messages
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				messages: entity.Messages{
					testMessage("message-id", now()),
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				msg := testMessage("message-id", now())
				err = db.DB.Create(&msg).Error
				require.NoError(t, err)
			},
			args: args{
				messages: entity.Messages{
					testMessage("message-id", now()),
				},
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			err := delete(ctx, messageTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &message{db: db, now: now}
			err = db.MultiCreate(ctx, tt.args.messages)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestMessage_UpdateRead(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	type args struct {
		messageID string
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				msg := testMessage("message-id", now())
				msg.Read = false
				err := db.DB.Create(&msg).Error
				require.NoError(t, err)
			},
			args: args{
				messageID: "message-id",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			err := delete(ctx, messageTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &message{db: db, now: now}
			err = db.UpdateRead(ctx, tt.args.messageID)
			assert.ErrorIs(t, err, tt.want.err)
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
