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

func TestThread(t *testing.T) {
	assert.NotNil(t, NewThread(nil))
}

func TestThread_ListByContactID(t *testing.T) {
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
		params *database.ListThreadsParams
	}
	type want struct {
		threads entity.Threads
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
				params: &database.ListThreadsParams{
					ContactID: "contact-id",
					Limit:     3,
					Offset:    0,
				},
			},
			want: want{
				threads: threads[:3],
				err:     nil,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &thread{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.ElementsMatch(t, tt.want.threads, actual)
		})
	}
}

func TestThread_Count(t *testing.T) {
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
		params *database.ListThreadsParams
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
				params: &database.ListThreadsParams{
					ContactID: "contact-id",
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

			db := &thread{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestThread_Create(t *testing.T) {
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
				thread: th,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				err = db.DB.Create(&th).Error
				require.NoError(t, err)
			},
			args: args{
				thread: th,
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

			err := delete(ctx, threadTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &thread{db: db, now: now}
			err = db.Create(ctx, tt.args.thread)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestThread_Get(t *testing.T) {
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
		err    error
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
				threadID: "thread-id",
			},
			want: want{
				thread: th,
				err:    nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				threadID: "other-id",
			},
			want: want{
				thread: nil,
				err:    database.ErrNotFound,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &thread{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.threadID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.thread, actual)
		})
	}
}

func TestThread_Update(t *testing.T) {
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

	category := testContactCategory("category-id", "お問い合わせ種別", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)

	c := testContact("contact-id", now())
	err = db.DB.Create(&c).Error
	require.NoError(t, err)

	type args = struct {
		threadID string
		params   *database.UpdateThreadParams
	}
	type want = struct {
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
				thread := testThread("thread-id", "contact-id", now())
				err = db.DB.Create(&thread).Error
				require.NoError(t, err)
			},
			args: args{
				threadID: "thread-id",
				params: &database.UpdateThreadParams{
					Content:  "会話内容",
					UserID:   "ユーザーID",
					UserType: 3,
				},
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

			err := delete(ctx, threadTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &thread{db: db, now: now}
			err = db.Update(ctx, tt.args.threadID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestThread_Delete(t *testing.T) {
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

	category := testContactCategory("category-id", "お問い合わせ種別", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)

	c := testContact("contact-id", now())
	err = db.DB.Create(&c).Error
	require.NoError(t, err)

	type args struct {
		threadID string
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
				thread := testThread("thread-id", "contact-id", now())
				err = db.DB.Create(&thread).Error
				require.NoError(t, err)
			},
			args: args{
				threadID: "thread-id",
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

			err := delete(ctx, threadTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &thread{db: db, now: now}
			err = db.Delete(ctx, tt.args.threadID)
			assert.ErrorIs(t, err, tt.want.err)
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
