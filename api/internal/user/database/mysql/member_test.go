package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMember_Get(t *testing.T) {
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
		member *entity.Member
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
				member: &u.Member,
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
				member: nil,
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

			db := &member{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.userID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.member, actual)
		})
	}
}

func TestMember_GetByCognitoID(t *testing.T) {
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
		cognitoID string
	}
	type want struct {
		member *entity.Member
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
				cognitoID: "user-id",
			},
			want: want{
				member: &u.Member,
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				cognitoID: "",
			},
			want: want{
				member: nil,
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

			db := &member{db: db, now: now}
			actual, err := db.GetByCognitoID(ctx, tt.args.cognitoID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.member, actual)
		})
	}
}

func TestMember_GetByEmail(t *testing.T) {
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
		email string
	}
	type want struct {
		member *entity.Member
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
				email: "test-user@and-period.jp",
			},
			want: want{
				member: &u.Member,
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				email: "test-other@and-period.jp",
			},
			want: want{
				member: nil,
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

			db := &member{db: db, now: now}
			actual, err := db.GetByEmail(ctx, tt.args.email)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.member, actual)
		})
	}
}

func TestMember_Create(t *testing.T) {
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

	user := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
	user.Member = *testMember("user-id", "test-user@and-period.jp", "+810000000000", now())
	user.Customer = entity.Customer{UserID: "user-id"}

	type args struct {
		user *entity.User
		auth func(ctx context.Context) error
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
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				user: testUser("user-id", "test-user@and-period.jp", "+810000000000", now()),
				auth: func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate user entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err = db.DB.Create(&u).Error
				require.NoError(t, err)
			},
			args: args{
				user: testUser("user-id", "test-user@and-period.jp", "+810000000000", now()),
				auth: func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to duplicate member entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err = db.DB.Create(&u).Error
				require.NoError(t, err)
				err = db.DB.Create(&u.Member).Error
				require.NoError(t, err)
			},
			args: args{
				user: testUser("user-id", "test-user@and-period.jp", "+810000000000", now()),
				auth: func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to duplicate customer entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err = db.DB.Create(&u).Error
				require.NoError(t, err)
				err = db.DB.Create(&u.Member).Error
				require.NoError(t, err)
				err = db.DB.Create(&u.Customer).Error
				require.NoError(t, err)
			},
			args: args{
				user: testUser("user-id", "test-user@and-period.jp", "+810000000000", now()),
				auth: func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name:  "failed to execute external service",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				user: testUser("user-id", "test-user@and-period.jp", "+810000000000", now()),
				auth: func(ctx context.Context) error { return assert.AnError },
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

			err := delete(ctx, customerTable, memberTable, userTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &member{db: db, now: now}
			err = db.Create(ctx, tt.args.user, tt.args.auth)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestMember_UpdateVerified(t *testing.T) {
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
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err = db.DB.Create(&u).Error
				require.NoError(t, err)
				u.Member.VerifiedAt = time.Time{}
				err = db.DB.Create(&u.Member).Error
				require.NoError(t, err)
			},
			args: args{
				userID: "user-id",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				userID: "user-id",
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to verified at is not zero value",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err = db.DB.Create(&u).Error
				require.NoError(t, err)
				u.Member.VerifiedAt = now()
				err = db.DB.Create(&u.Member).Error
				require.NoError(t, err)
			},
			args: args{
				userID: "user-id",
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

			err := delete(ctx, memberTable, userTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &member{db: db, now: now}
			err = db.UpdateVerified(ctx, tt.args.userID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestUser_UpdateAccount(t *testing.T) {
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

	type args struct {
		userID    string
		accountID string
		username  string
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
				user := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err := db.DB.Create(&user).Error
				require.NoError(t, err)
				err = db.DB.Create(&user.Member).Error
				require.NoError(t, err)
			},
			args: args{
				userID:    "user-id",
				accountID: "account-id",
				username:  "username",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate account id",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				user := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err := db.DB.Create(&user).Error
				require.NoError(t, err)
				user.Member.AccountID = ""
				err = db.DB.Create(&user.Member).Error
				require.NoError(t, err)

				other := testUser("other-id", "test-other@and-period.jp", "+81111111111", now())
				err = db.DB.Create(&other).Error
				require.NoError(t, err)
				other.Member.AccountID = "account-id"
				err = db.DB.Create(&other.Member).Error
				require.NoError(t, err)
			},
			args: args{
				userID:    "user-id",
				accountID: "account-id",
				username:  "username",
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

			err := delete(ctx, memberTable, userTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &member{db: db, now: now}
			err = db.UpdateAccount(ctx, tt.args.userID, tt.args.accountID, tt.args.username)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestUser_UpdateEmail(t *testing.T) {
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

	type args struct {
		userID string
		email  string
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
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err = db.DB.Create(&u).Error
				require.NoError(t, err)
				err = db.DB.Create(&u.Member).Error
				require.NoError(t, err)
			},
			args: args{
				userID: "user-id",
				email:  "test-other@and-period.jp",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				userID: "user-id",
				email:  "test-other@and-period.jp",
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to unmatch provider type",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				u.ProviderType = entity.ProviderTypeOAuth
				err = db.DB.Create(&u).Error
				require.NoError(t, err)
				err = db.DB.Create(&u.Member).Error
				require.NoError(t, err)
			},
			args: args{
				userID: "user-id",
				email:  "",
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

			err := delete(ctx, customerTable, memberTable, userTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &member{db: db, now: now}
			err = db.UpdateEmail(ctx, tt.args.userID, tt.args.email)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestMember_Delete(t *testing.T) {
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

	type args struct {
		userID string
		auth   func(ctx context.Context) error
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
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err = db.DB.Create(&u).Error
				require.NoError(t, err)
				err = db.DB.Create(&u.Member).Error
				require.NoError(t, err)
			},
			args: args{
				userID: "user-id",
				auth:   func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to execute external service",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				userID: "user-id",
				auth:   func(ctx context.Context) error { return assert.AnError },
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

			err := delete(ctx, memberTable, userTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &member{db: db, now: now}
			err = db.Delete(ctx, tt.args.userID, tt.args.auth)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testMember(id, email, phoneNumber string, now time.Time) *entity.Member {
	return &entity.Member{
		UserID:       id,
		AccountID:    id,
		CognitoID:    id,
		Username:     id,
		ProviderType: entity.ProviderTypeEmail,
		Email:        email,
		PhoneNumber:  phoneNumber,
		ThumbnailURL: "https://and-period.jp/thumbnail.png",
		CreatedAt:    now,
		UpdatedAt:    now,
		VerifiedAt:   now,
	}
}