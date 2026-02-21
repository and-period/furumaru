package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"go.uber.org/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminRole_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	roles := make(entity.AdminRoles, 3)
	roles[0] = testAdminRole("role-id01", now())
	roles[1] = testAdminRole("role-id02", now())
	roles[2] = testAdminRole("role-id03", now())
	err = db.DB.WithContext(ctx).Table(adminRoleTable).Create(&roles).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdminRolesParams
	}
	type want struct {
		roles entity.AdminRoles
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
				params: &database.ListAdminRolesParams{},
			},
			want: want{
				roles: roles,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &adminRole{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.roles, actual)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdminRole_Count(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	roles := make(entity.AdminRoles, 3)
	roles[0] = testAdminRole("role-id01", now())
	roles[1] = testAdminRole("role-id02", now())
	roles[2] = testAdminRole("role-id03", now())
	err = db.DB.WithContext(ctx).Table(adminRoleTable).Create(&roles).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdminRolesParams
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
				params: &database.ListAdminRolesParams{},
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
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &adminRole{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.total, actual)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdminRole_MultiGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	roles := make(entity.AdminRoles, 3)
	roles[0] = testAdminRole("role-id01", now())
	roles[1] = testAdminRole("role-id02", now())
	roles[2] = testAdminRole("role-id03", now())
	err = db.DB.WithContext(ctx).Table(adminRoleTable).Create(&roles).Error
	require.NoError(t, err)

	type args struct {
		roleIDs []string
	}
	type want struct {
		roles entity.AdminRoles
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
				roleIDs: []string{"role-id01", "role-id02", "role-id03"},
			},
			want: want{
				roles: roles,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &adminRole{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.roleIDs)
			assert.Equal(t, tt.want.roles, actual)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdminRole_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	r := testAdminRole("role-id", now())
	err = db.DB.WithContext(ctx).Table(adminRoleTable).Create(&r).Error
	require.NoError(t, err)

	type args struct {
		roleID string
	}
	type want struct {
		role *entity.AdminRole
		err  error
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
				roleID: "role-id",
			},
			want: want{
				role: r,
				err:  nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				roleID: "",
			},
			want: want{
				role: nil,
				err:  database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &adminRole{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.roleID)
			assert.Equal(t, tt.want.role, actual)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testAdminRole(roleID string, now time.Time) *entity.AdminRole {
	return &entity.AdminRole{
		ID:          roleID,
		Name:        "スポット編集者",
		Description: "スポットの編集ができる権限",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
