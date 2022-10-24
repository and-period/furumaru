package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
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

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	err = m.dbDelete(ctx, customerTable, memberTable, userTable)
	require.NoError(t, err)
	users := make(entity.Users, 2)
	users[0] = testUser("user-id01", "test-user01@and-period.jp", "+810000000001", now())
	users[1] = testUser("user-id02", "test-user02@and-period.jp", "+810000000002", now())
	err = m.db.DB.Create(&users).Error
	for i := range users {
		err = m.db.DB.Create(&users[i].Member).Error
		err = m.db.DB.Create(&users[i].Customer).Error
	}
	require.NoError(t, err)

	type args struct {
		params *ListUsersParams
	}
	type want struct {
		users  entity.Users
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
				params: &ListUsersParams{
					Limit:  1,
					Offset: 1,
				},
			},
			want: want{
				users:  users[1:],
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

			db := &user{db: m.db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreUsersField(actual, now())
			assert.ElementsMatch(t, tt.want.users, actual)
		})
	}
}

func TestUser_Count(t *testing.T) {
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

	err = m.dbDelete(ctx, customerTable, memberTable, userTable)
	require.NoError(t, err)
	users := make(entity.Users, 2)
	users[0] = testUser("user-id01", "test-user01@and-period.jp", "+810000000001", now())
	users[1] = testUser("user-id02", "test-user02@and-period.jp", "+810000000002", now())
	err = m.db.DB.Create(&users).Error
	for i := range users {
		err = m.db.DB.Create(&users[i].Member).Error
		err = m.db.DB.Create(&users[i].Customer).Error
	}
	require.NoError(t, err)

	type args struct {
		params *ListUsersParams
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
				params: &ListUsersParams{
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, m)

			db := &user{db: m.db, now: now}
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

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	err = m.dbDelete(ctx, customerTable, memberTable, userTable)
	require.NoError(t, err)
	users := make(entity.Users, 2)
	users[0] = testUser("user-id01", "test-user01@and-period.jp", "+810000000001", now())
	users[1] = testUser("user-id02", "test-user02@and-period.jp", "+810000000002", now())
	err = m.db.DB.Create(&users).Error
	for i := range users {
		err = m.db.DB.Create(&users[i].Member).Error
		err = m.db.DB.Create(&users[i].Customer).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, m)

			db := &user{db: m.db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.userIDs)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreUsersField(actual, now())
			assert.ElementsMatch(t, tt.want.users, actual)
		})
	}
}

func TestUser_Get(t *testing.T) {
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

	err = m.dbDelete(ctx, customerTable, memberTable, userTable)
	require.NoError(t, err)
	u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
	err = m.db.DB.Create(&u).Error
	require.NoError(t, err)
	err = m.db.DB.Create(&u.Member).Error
	require.NoError(t, err)
	err = m.db.DB.Create(&u.Customer).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, m)

			db := &user{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.userID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreUserField(actual, now())
			assert.Equal(t, tt.want.user, actual)
		})
	}
}

func testUser(id, email, phoneNumber string, now time.Time) *entity.User {
	return &entity.User{
		ID:         id,
		Registered: true,
		Member:     *testMember(id, email, phoneNumber, now),
		Customer:   *testCustomer(id, now),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func fillIgnoreUserField(u *entity.User, now time.Time) {
	if u == nil {
		return
	}
	u.CreatedAt = now
	u.UpdatedAt = now
	u.Member.CreatedAt = now
	u.Member.UpdatedAt = now
	u.Member.VerifiedAt = now
	u.Customer.CreatedAt = now
	u.Customer.UpdatedAt = now
}

func fillIgnoreUsersField(us entity.Users, now time.Time) {
	for i := range us {
		fillIgnoreUserField(us[i], now)
	}
}
