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

func TestMember_Get(t *testing.T) {
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

	err = m.dbDelete(ctx, memberTable, userTable)
	require.NoError(t, err)
	u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
	err = m.db.DB.Create(&u).Error
	require.NoError(t, err)
	err = m.db.DB.Create(&u.Member).Error
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
				member: &u.Member,
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

			tt.setup(ctx, t, m)

			db := &member{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.userID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreMemberField(actual, now())
			assert.Equal(t, tt.want.member, actual)
		})
	}
}

func TestMember_GetByCognitoID(t *testing.T) {
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

	err = m.dbDelete(ctx, memberTable, userTable)
	require.NoError(t, err)
	u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
	err = m.db.DB.Create(&u).Error
	require.NoError(t, err)
	err = m.db.DB.Create(&u.Member).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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

			tt.setup(ctx, t, m)

			db := &member{db: m.db, now: now}
			actual, err := db.GetByCognitoID(ctx, tt.args.cognitoID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreMemberField(actual, now())
			assert.Equal(t, tt.want.member, actual)
		})
	}
}

func TestMember_GetByEmail(t *testing.T) {
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

	err = m.dbDelete(ctx, memberTable, userTable)
	require.NoError(t, err)
	u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
	err = m.db.DB.Create(&u).Error
	require.NoError(t, err)
	err = m.db.DB.Create(&u.Member).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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

			tt.setup(ctx, t, m)

			db := &member{db: m.db, now: now}
			actual, err := db.GetByEmail(ctx, tt.args.email)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreMemberField(actual, now())
			assert.Equal(t, tt.want.member, actual)
		})
	}
}

func TestMember_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	type args struct {
		user   *entity.User
		member *entity.Member
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
				user:   testUser("user-id", "test-user@and-period.jp", "+810000000000", now()),
				member: testMember("user-id", "test-user@and-period.jp", "+810000000000", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "success when user already created",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err = m.db.DB.Create(&u).Error
				require.NoError(t, err)
			},
			args: args{
				user:   testUser("user-id", "test-user@and-period.jp", "+810000000000", now()),
				member: testMember("user-id", "test-user@and-period.jp", "+810000000000", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err = m.db.DB.Create(&u).Error
				require.NoError(t, err)
				err = m.db.DB.Create(&u.Member).Error
				require.NoError(t, err)
			},
			args: args{
				user:   testUser("user-id", "test-user@and-period.jp", "+810000000000", now()),
				member: testMember("user-id", "test-user@and-period.jp", "+810000000000", now()),
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

			err := m.dbDelete(ctx, customerTable, memberTable, userTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &member{db: m.db, now: now}
			err = db.Create(ctx, tt.args.user, tt.args.member)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestMember_UpdateVerified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	type args struct {
		userID string
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
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err = m.db.DB.Create(&u).Error
				require.NoError(t, err)
				u.Member.VerifiedAt = time.Time{}
				err = m.db.DB.Create(&u.Member).Error
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				userID: "user-id",
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to verified at is not zero value",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err = m.db.DB.Create(&u).Error
				require.NoError(t, err)
				u.Member.VerifiedAt = now()
				err = m.db.DB.Create(&u.Member).Error
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

			err := m.dbDelete(ctx, memberTable, userTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &member{db: m.db, now: now}
			err = db.UpdateVerified(ctx, tt.args.userID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestUser_UpdateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				user := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err := m.db.DB.Create(&user).Error
				require.NoError(t, err)
				err = m.db.DB.Create(&user.Member).Error
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
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				userID:    "user-id",
				accountID: "account-id",
				username:  "username",
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to duplicate account id",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				user := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err := m.db.DB.Create(&user).Error
				require.NoError(t, err)
				user.Member.AccountID = ""
				err = m.db.DB.Create(&user.Member).Error
				require.NoError(t, err)

				other := testUser("other-id", "test-other@and-period.jp", "+81111111111", now())
				err = m.db.DB.Create(&other).Error
				require.NoError(t, err)
				other.Member.AccountID = "account-id"
				err = m.db.DB.Create(&other.Member).Error
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

			err := m.dbDelete(ctx, memberTable, userTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &member{db: m.db, now: now}
			err = db.UpdateAccount(ctx, tt.args.userID, tt.args.accountID, tt.args.username)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestUser_UpdateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	type args struct {
		userID string
		email  string
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
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err = m.db.DB.Create(&u).Error
				require.NoError(t, err)
				err = m.db.DB.Create(&u.Member).Error
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				u.ProviderType = entity.ProviderTypeOAuth
				err = m.db.DB.Create(&u).Error
				require.NoError(t, err)
				err = m.db.DB.Create(&u.Member).Error
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

			err := m.dbDelete(ctx, customerTable, memberTable, userTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &member{db: m.db, now: now}
			err = db.UpdateEmail(ctx, tt.args.userID, tt.args.email)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestMember_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	type args struct {
		userID string
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
				u := testUser("user-id", "test-user@and-period.jp", "+810000000000", now())
				err = m.db.DB.Create(&u).Error
				require.NoError(t, err)
				err = m.db.DB.Create(&u.Member).Error
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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

			err := m.dbDelete(ctx, memberTable, userTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &member{db: m.db, now: now}
			err = db.Delete(ctx, tt.args.userID)
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

func fillIgnoreMemberField(m *entity.Member, now time.Time) {
	if m == nil {
		return
	}
	m.CreatedAt = now
	m.UpdatedAt = now
	m.VerifiedAt = now
}
