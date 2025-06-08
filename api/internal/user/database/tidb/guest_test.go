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

func TestGuest_GetByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	u := testGuestUser("user-id", "test-user@and-period.jp", now())
	err = db.DB.Create(&u).Error
	require.NoError(t, err)
	err = db.DB.Create(&u.Guest).Error
	require.NoError(t, err)

	type args struct {
		email string
	}
	type want struct {
		guest *entity.Guest
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
				email: "test-user@and-period.jp",
			},
			want: want{
				guest: &u.Guest,
				err:   nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				email: "test-other@and-period.jp",
			},
			want: want{
				guest: nil,
				err:   database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &guest{db: db, now: now}
			actual, err := db.GetByEmail(ctx, tt.args.email)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.guest, actual)
		})
	}
}

func TestGuest_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	type args struct {
		user *entity.User
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
				user: testGuestUser("user-id", "test-user@and-period.jp", now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "duplicate user entity",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				u := testGuestUser("user-id", "test-user@and-period.jp", now())
				err := db.DB.Create(&u).Error
				require.NoError(t, err)
				err = db.DB.Create(&u.Guest).Error
				require.NoError(t, err)
			},
			args: args{
				user: testGuestUser("user-id", "test-user@and-period.jp", now()),
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, guestTable, userTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &guest{db: db, now: now}
			err = db.Create(ctx, tt.args.user)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestGuest_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	type args struct {
		userID string
		params *database.UpdateGuestParams
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
				u := testGuestUser("user-id", "test-user@and-period.jp", now())
				err := db.DB.Create(&u).Error
				require.NoError(t, err)
				err = db.DB.Create(&u.Guest).Error
				require.NoError(t, err)
			},
			args: args{
				userID: "user-id",
				params: &database.UpdateGuestParams{
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, guestTable, userTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &guest{db: db, now: now}
			err = db.Update(ctx, tt.args.userID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestGuest_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	type args struct {
		userID string
	}
	type want struct {
		hasErr bool
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
				u := testGuestUser("user-id", "test-user@and-period.jp", now())
				err = db.DB.Create(&u).Error
				require.NoError(t, err)
				err = db.DB.Create(&u.Guest).Error
				require.NoError(t, err)
			},
			args: args{
				userID: "user-id",
			},
			want: want{
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, guestTable, userTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &guest{db: db, now: now}
			err = db.Delete(ctx, tt.args.userID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testGuest(id, email string, now time.Time) *entity.Guest {
	return &entity.Guest{
		UserID:    id,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
