package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	assert.NotNil(t, NewUser(nil))
}

func TestUser_List(t *testing.T) {
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

	users := make(entity.Users, 2)
	users[0] = testUser("user-id01", "test-user01@and-period.jp", "+810000000001", now())
	users[1] = testUser("user-id02", "test-user02@and-period.jp", "+810000000002", now().Add(time.Hour))
	err = db.DB.Create(&users).Error
	for i := range users {
		err = db.DB.Create(&users[i].Member).Error
	}
	require.NoError(t, err)

	type args struct {
		params *database.ListUsersParams
	}
	type want struct {
		users  entity.Users
		hasErr bool
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
				params: &database.ListUsersParams{
					Limit:  1,
					Offset: 1,
				},
			},
			want: want{
				users:  users[:1],
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &user{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.users, actual)
		})
	}
}

func TestUser_Count(t *testing.T) {
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

	users := make(entity.Users, 2)
	users[0] = testUser("user-id01", "test-user01@and-period.jp", "+810000000001", now())
	users[1] = testUser("user-id02", "test-user02@and-period.jp", "+810000000002", now())
	err = db.DB.Create(&users).Error
	for i := range users {
		err = db.DB.Create(&users[i].Member).Error
	}
	require.NoError(t, err)

	type args struct {
		params *database.ListUsersParams
	}
	type want struct {
		total  int64
		hasErr bool
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
				params: &database.ListUsersParams{
					Limit:  1,
					Offset: 1,
				},
			},
			want: want{
				total:  2,
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &user{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestUser_MultiGet(t *testing.T) {
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

	users := make(entity.Users, 2)
	users[0] = testUser("user-id01", "test-user01@and-period.jp", "+810000000001", now())
	users[1] = testUser("user-id02", "test-user02@and-period.jp", "+810000000002", now())
	err = db.DB.Create(&users).Error
	for i := range users {
		err = db.DB.Create(&users[i].Member).Error
	}
	require.NoError(t, err)

	type args struct {
		userIDs []string
	}
	type want struct {
		users  entity.Users
		hasErr bool
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
				userIDs: []string{"user-id01", "user-id02"},
			},
			want: want{
				users:  users,
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &user{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.userIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.users, actual)
		})
	}
}

func TestUser_Get(t *testing.T) {
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

	u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
	err = db.DB.Create(&u).Error
	require.NoError(t, err)
	err = db.DB.Create(&u.Member).Error
	require.NoError(t, err)

	type args struct {
		userID string
	}
	type want struct {
		user   *entity.User
		hasErr bool
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
				userID: "user-id",
			},
			want: want{
				user:   u,
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				userID: "",
			},
			want: want{
				user:   nil,
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &user{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.userID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.user, actual)
		})
	}
}

func testUser(id, email, phoneNumber string, now time.Time) *entity.User {
	return &entity.User{
		ID:         id,
		Registered: true,
		Status:     entity.UserStatusVerified,
		Member:     *testMember(id, email, phoneNumber, now),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func testGuestUser(id, email string, now time.Time) *entity.User {
	return &entity.User{
		ID:         id,
		Registered: false,
		Status:     entity.UserStatusGuest,
		Guest:      *testGuest(id, email, now),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
