package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThread(t *testing.T) {
	assert.NotNil(t, NewThread(nil))
}

func TestThread_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	category := testContactCategory("category-id", "お問い合わせ種別", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)

	c := testContact("contact-id", now())
	err = db.DB.Create(&c).Error
	require.NoError(t, err)

	th := testThread("thread-id", now())
	err = db.DB.Create(&th).Error
	require.NoError(t, err)

	type args struct {
		threadID string
	}
	type want struct {
		thread *entity.Thread
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				threadID: "thread-id",
			},
			want: want{
				thread: th,
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				threadID: "other-id",
			},
			want: want{
				thread: nil,
				hasErr: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &thread{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.threadID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.thread, actual)
		})
	}
}

func testThread(id string, now time.Time) *entity.Thread {
	return &entity.Thread{
		ID:        id,
		ContactID: "contact-id",
		UserID:    "user-id",
		UserType:  1,
		Content:   "content",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
