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

func TestThread_ListByContactID(t *testing.T) {
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

	threads := make(entity.Threads, 3)
	threads[0] = testThread("thread-id01", "contact-id", now())
	threads[1] = testThread("thread-id02", "contact-id", now())
	threads[2] = testThread("thread-id03", "contact-id", now())
	err = db.DB.Create(&threads).Error
	require.NoError(t, err)

	type args struct {
		params *ListThreadsByContactIDParams
	}
	type want struct {
		threads entity.Threads
		hasErr  bool
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
				params: &ListThreadsByContactIDParams{
					ContactID: "contact-id",
					Limit:     3,
					Offset:    0,
				},
			},
			want: want{
				threads: threads[:3],
				hasErr:  false,
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
			actual, err := db.ListByContactID(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.threads, actual)
		})
	}
}

func TestProductType_Count(t *testing.T) {
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

	threads := make(entity.Threads, 3)
	threads[0] = testThread("thread-id01", "contact-id", now())
	threads[1] = testThread("thread-id02", "contact-id", now())
	threads[2] = testThread("thread-id03", "contact-id", now())
	err = db.DB.Create(&threads).Error
	require.NoError(t, err)

	type args struct {
		params *ListThreadsByContactIDParams
	}
	type want struct {
		total  int64
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
				params: &ListThreadsByContactIDParams{
					ContactID: "contact-id",
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

			tt.setup(ctx, t, db)

			db := &thread{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestThread_Create(t *testing.T) {
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

	contact := testContact("contact-id", now())
	err = db.DB.Create(&contact).Error
	require.NoError(t, err)

	th := testThread("thread-id", "contact-id", now())

	type args struct {
		thread *entity.Thread
	}
	type want struct {
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
				thread: th,
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				err = db.DB.Create(&th).Error
				require.NoError(t, err)
			},
			args: args{
				thread: th,
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

			err := delete(ctx, threadTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &thread{db: db, now: now}
			err = db.Create(ctx, tt.args.thread)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
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

	th := testThread("thread-id", "contact-id", now())
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

func TestThread_Update(t *testing.T) {
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

	type args = struct {
		threadID string
		params   *UpdateThreadParams
	}
	type want = struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				thread := testThread("thread-id", "contact-id", now())
				err = db.DB.Create(&thread).Error
				require.NoError(t, err)
			},
			args: args{
				threadID: "thread-id",
				params: &UpdateThreadParams{
					Content:  "会話内容",
					UserID:   "ユーザーID",
					UserType: 3,
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				threadID: "thread-id",
				params:   &UpdateThreadParams{},
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

			err := delete(ctx, threadTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &thread{db: db, now: now}
			err = db.Update(ctx, tt.args.threadID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testThread(id string, contactID string, now time.Time) *entity.Thread {
	return &entity.Thread{
		ID:        id,
		ContactID: contactID,
		UserID:    "user-id",
		UserType:  1,
		Content:   "content",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
