package mysql

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

func TestAdministrator(t *testing.T) {
	assert.NotNil(t, NewAdministrator(nil))
}

func TestAdministrator_List(t *testing.T) {
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

	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	require.NoError(t, err)
	administrators := make(entity.Administrators, 2)
	administrators[0] = testAdministrator("admin-id01", now())
	administrators[0].Admin = *admins[0]
	administrators[1] = testAdministrator("admin-id02", now())
	administrators[1].Admin = *admins[1]
	err = db.DB.Create(&administrators).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdministratorsParams
	}
	type want struct {
		admins entity.Administrators
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
				params: &database.ListAdministratorsParams{
					Limit:  1,
					Offset: 1,
				},
			},
			want: want{
				admins: administrators[1:],
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

			db := &administrator{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.admins, actual)
		})
	}
}

func TestAdministrator_Count(t *testing.T) {
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

	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	require.NoError(t, err)
	administrators := make(entity.Administrators, 2)
	administrators[0] = testAdministrator("admin-id01", now())
	administrators[0].Admin = *admins[0]
	administrators[1] = testAdministrator("admin-id02", now())
	administrators[1].Admin = *admins[1]
	err = db.DB.Create(&administrators).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdministratorsParams
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
				params: &database.ListAdministratorsParams{},
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

			tt.setup(ctx, t, db)

			db := &administrator{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestAdministrator_MultiGet(t *testing.T) {
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

	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	require.NoError(t, err)
	administrators := make(entity.Administrators, 2)
	administrators[0] = testAdministrator("admin-id01", now())
	administrators[0].Admin = *admins[0]
	administrators[1] = testAdministrator("admin-id02", now())
	administrators[1].Admin = *admins[1]
	err = db.DB.Create(&administrators).Error
	require.NoError(t, err)

	type args struct {
		adminIDs []string
	}
	type want struct {
		admins entity.Administrators
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
				admins: administrators,
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

			db := &administrator{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.adminIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.admins, actual)
		})
	}
}

func TestAdministrator_Get(t *testing.T) {
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

	admin := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
	err = db.DB.Create(&admin).Error
	require.NoError(t, err)
	a := testAdministrator("admin-id", now())
	a.Admin = *admin
	err = db.DB.Create(&a).Error
	require.NoError(t, err)

	type args struct {
		adminID string
	}
	type want struct {
		admin  *entity.Administrator
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &administrator{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.adminID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.admin, actual)
		})
	}
}

func TestAdministrator_Create(t *testing.T) {
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

	a := testAdministrator("admin-id", now())
	a.Admin = *testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())

	type args struct {
		administrator *entity.Administrator
		auth          func(ctx context.Context) error
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
				administrator: a,
				auth:          func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry in admin auth",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				auth := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
				err = db.DB.Create(&auth).Error
				require.NoError(t, err)
			},
			args: args{
				administrator: a,
				auth:          func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name:  "failed to execute external service",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				administrator: a,
				auth:          func(ctx context.Context) error { return assert.AnError },
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

			err := delete(ctx, administratorTable, adminTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &administrator{db: db, now: now}
			err = db.Create(ctx, tt.args.administrator, tt.args.auth)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestAdministrator_Update(t *testing.T) {
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
		administratorID string
		params          *database.UpdateAdministratorParams
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
				admin := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				administrator := testAdministrator("admin-id", now())
				err = db.DB.Create(&administrator).Error
				require.NoError(t, err)
			},
			args: args{
				administratorID: "admin-id",
				params: &database.UpdateAdministratorParams{
					Lastname:      "&.",
					Firstname:     "スタッフ",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "すたっふ",
					PhoneNumber:   "+819012345678",
				},
			},
			want: want{
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, administratorTable, adminTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &administrator{db: db, now: now}
			err = db.Update(ctx, tt.args.administratorID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestAdministrator_Delete(t *testing.T) {
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
		administratorID string
		auth            func(ctx context.Context) error
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
				admin := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				administrator := testAdministrator("admin-id", now())
				err = db.DB.Create(&administrator).Error
				require.NoError(t, err)
			},
			args: args{
				administratorID: "admin-id",
				auth:            func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to execute external service",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				admin := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				administrator := testAdministrator("admin-id", now())
				err = db.DB.Create(&administrator).Error
				require.NoError(t, err)
			},
			args: args{
				administratorID: "admin-id",
				auth:            func(ctx context.Context) error { return assert.AnError },
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

			err := delete(ctx, administratorTable, adminTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &administrator{db: db, now: now}
			err = db.Delete(ctx, tt.args.administratorID, tt.args.auth)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testAdministrator(id string, now time.Time) *entity.Administrator {
	return &entity.Administrator{
		AdminID:     id,
		PhoneNumber: "+819012345678",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
