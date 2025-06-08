package tidb

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

func TestAdmin(t *testing.T) {
	assert.NotNil(t, NewAdmin(nil))
}

func TestAdmin_MultiGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin1@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin2@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	require.NoError(t, err)

	type args struct {
		adminIDs []string
	}
	type want struct {
		admins entity.Admins
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
				adminIDs: []string{"admin-id01", "admin-id02"},
			},
			want: want{
				admins: admins,
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &admin{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.adminIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.admins, actual)
		})
	}
}

func TestAdmin_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	a := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
	err = db.DB.Create(&a).Error
	require.NoError(t, err)

	type args struct {
		adminID string
	}
	type want struct {
		admin  *entity.Admin
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
				adminID: "admin-id",
			},
			want: want{
				admin:  a,
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				adminID: "",
			},
			want: want{
				admin:  nil,
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &admin{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.adminID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.admin, actual)
		})
	}
}

func TestAdmin_GetByCognitoID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	a := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
	err = db.DB.Create(&a).Error
	require.NoError(t, err)

	type args struct {
		cognitoID string
	}
	type want struct {
		admin  *entity.Admin
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
				cognitoID: "cognito-id",
			},
			want: want{
				admin:  a,
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
				admin:  nil,
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &admin{db: db, now: now}
			actual, err := db.GetByCognitoID(ctx, tt.args.cognitoID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.admin, actual)
		})
	}
}

func TestAdmin_GetByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	a := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
	err = db.DB.Create(&a).Error
	require.NoError(t, err)

	type args struct {
		email string
	}
	type want struct {
		admin  *entity.Admin
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
				email: "test-admin@and-period.jp",
			},
			want: want{
				admin:  a,
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				email: "",
			},
			want: want{
				admin:  nil,
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &admin{db: db, now: now}
			actual, err := db.GetByEmail(ctx, tt.args.email)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.admin, actual)
		})
	}
}

func TestAdmin_UpdateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	type args struct {
		adminID string
		email   string
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
				a := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
				err = db.DB.Create(&a).Error
				require.NoError(t, err)
			},
			args: args{
				adminID: "admin-id",
				email:   "other@and-period.jp",
			},
			want: want{
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, adminTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &admin{db: db, now: now}
			err = db.UpdateEmail(ctx, tt.args.adminID, tt.args.email)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestAdmin_UpdateDevice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	type args struct {
		adminID string
		device  string
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
				a := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
				err = db.DB.Create(&a).Error
				require.NoError(t, err)
			},
			args: args{
				adminID: "admin-id",
				device:  "device",
			},
			want: want{
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, adminTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &admin{db: db, now: now}
			err = db.UpdateDevice(ctx, tt.args.adminID, tt.args.device)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testAdmin(adminID, cognitoID, email string, now time.Time) *entity.Admin {
	return &entity.Admin{
		ID:            adminID,
		CognitoID:     cognitoID,
		Type:          entity.AdminTypeAdministrator,
		GroupIDs:      []string{},
		Status:        entity.AdminStatusActivated,
		Lastname:      "&.",
		Firstname:     "スタッフ",
		LastnameKana:  "あんどぴりおど",
		FirstnameKana: "すたっふ",
		Email:         email,
		Device:        "instance-id",
		FirstSignInAt: now,
		LastSignInAt:  now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
