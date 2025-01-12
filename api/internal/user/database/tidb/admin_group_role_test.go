package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAdminGroupRole_List(t *testing.T) {
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
	err = db.DB.WithContext(ctx).Table(adminTable).Create(&admin).Error
	require.NoError(t, err)

	g := testAdminGroup("group-id", "admin-id", now())
	err = db.DB.WithContext(ctx).Table(adminGroupTable).Create(&g).Error
	require.NoError(t, err)

	roles := make(entity.AdminRoles, 3)
	roles[0] = testAdminRole("role-id01", now())
	roles[1] = testAdminRole("role-id02", now())
	roles[2] = testAdminRole("role-id03", now())
	err = db.DB.WithContext(ctx).Table(adminRoleTable).Create(&roles).Error
	require.NoError(t, err)

	groupRoles := make(entity.AdminGroupRoles, 3)
	groupRoles[0] = testAdminGroupRole("group-id", roles[0].ID, "admin-id", now().Add(-2*time.Hour))
	groupRoles[1] = testAdminGroupRole("group-id", roles[1].ID, "admin-id", now().Add(-time.Hour))
	groupRoles[2] = testAdminGroupRole("group-id", roles[2].ID, "admin-id", now())
	err = db.DB.WithContext(ctx).Table(adminGroupRoleTable).Create(&groupRoles).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdminGroupRolesParams
	}
	type want struct {
		groupRoles entity.AdminGroupRoles
		err        error
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
				params: &database.ListAdminGroupRolesParams{},
			},
			want: want{
				groupRoles: groupRoles,
				err:        nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &adminGroupRole{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			require.Equal(t, tt.want.groupRoles, actual)
			require.Equal(t, tt.want.err, err)
		})
	}
}

func TestAdminGroupRole_Count(t *testing.T) {
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
	err = db.DB.WithContext(ctx).Table(adminTable).Create(&admin).Error
	require.NoError(t, err)

	g := testAdminGroup("group-id", "admin-id", now())
	err = db.DB.WithContext(ctx).Table(adminGroupTable).Create(&g).Error
	require.NoError(t, err)

	roles := make(entity.AdminRoles, 3)
	roles[0] = testAdminRole("role-id01", now())
	roles[1] = testAdminRole("role-id02", now())
	roles[2] = testAdminRole("role-id03", now())
	err = db.DB.WithContext(ctx).Table(adminRoleTable).Create(&roles).Error
	require.NoError(t, err)

	groupRoles := make(entity.AdminGroupRoles, 3)
	groupRoles[0] = testAdminGroupRole("group-id", roles[0].ID, "admin-id", now().Add(-2*time.Hour))
	groupRoles[1] = testAdminGroupRole("group-id", roles[1].ID, "admin-id", now().Add(-time.Hour))
	groupRoles[2] = testAdminGroupRole("group-id", roles[2].ID, "admin-id", now())
	err = db.DB.WithContext(ctx).Table(adminGroupRoleTable).Create(&groupRoles).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdminGroupRolesParams
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
				params: &database.ListAdminGroupRolesParams{},
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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &adminGroupRole{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			require.Equal(t, tt.want.total, actual)
			require.Equal(t, tt.want.err, err)
		})
	}
}

func TestAdminGroupRole_Get(t *testing.T) {
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
	err = db.DB.WithContext(ctx).Table(adminTable).Create(&admin).Error
	require.NoError(t, err)

	g := testAdminGroup("group-id", "admin-id", now())
	err = db.DB.WithContext(ctx).Table(adminGroupTable).Create(&g).Error
	require.NoError(t, err)

	role := testAdminRole("role-id", now())
	err = db.DB.WithContext(ctx).Table(adminRoleTable).Create(&role).Error
	require.NoError(t, err)

	gr := testAdminGroupRole("group-id", "role-id", "admin-id", now())
	err = db.DB.WithContext(ctx).Table(adminGroupRoleTable).Create(&gr).Error
	require.NoError(t, err)

	type args struct {
		groupID string
		roleID  string
	}
	type want struct {
		group *entity.AdminGroup
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
				groupID: "group-id",
				roleID:  "role-id",
			},
			want: want{
				group: g,
				err:   nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				groupID: "",
				roleID:  "",
			},
			want: want{
				group: nil,
				err:   database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &adminGroupRole{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.groupID, tt.args.roleID)
			require.Equal(t, tt.want.group, actual)
			require.Equal(t, tt.want.err, err)
		})
	}
}

func TestAdminGroupRole_Upsert(t *testing.T) {
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
	err = db.DB.WithContext(ctx).Table(adminTable).Create(&admin).Error
	require.NoError(t, err)

	group := testAdminGroup("group-id", "admin-id", now())
	err = db.DB.WithContext(ctx).Table(adminGroupTable).Create(&group).Error
	require.NoError(t, err)

	role := testAdminRole("role-id", now())
	err = db.DB.WithContext(ctx).Table(adminRoleTable).Create(&role).Error
	require.NoError(t, err)

	type args struct {
		groupRole *entity.AdminGroupRole
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
			name:  "success to create",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				groupRole: testAdminGroupRole("group-id", "role-id", "admin-id", now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success to update",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				gr := testAdminGroupRole("group-id", "role-id", "admin-id", now())
				err := db.DB.WithContext(ctx).Table(adminGroupRoleTable).Create(&gr).Error
				require.NoError(t, err)
			},
			args: args{
				groupRole: testAdminGroupRole("group-id", "role-id", "admin-id", now()),
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, adminGroupRoleTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &adminGroupRole{db: db, now: now}
			err = db.Upsert(ctx, tt.args.groupRole)
			require.Equal(t, tt.want.err, err)
		})
	}
}

func TestAdminGroupRole_Delete(t *testing.T) {
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
	err = db.DB.WithContext(ctx).Table(adminTable).Create(&admin).Error
	require.NoError(t, err)

	group := testAdminGroup("group-id", "admin-id", now())
	err = db.DB.WithContext(ctx).Table(adminGroupTable).Create(&group).Error
	require.NoError(t, err)

	role := testAdminRole("role-id", now())
	err = db.DB.WithContext(ctx).Table(adminRoleTable).Create(&role).Error
	require.NoError(t, err)

	type args struct {
		groupID string
		roleID  string
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
				gr := testAdminGroupRole("group-id", "role-id", "admin-id", now())
				err = db.DB.WithContext(ctx).Table(adminGroupRoleTable).Create(&gr).Error
				require.NoError(t, err)
			},
			args: args{
				groupID: "group-id",
				roleID:  "role-id",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, adminGroupRoleTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &adminGroupRole{db: db, now: now}
			err = db.Delete(ctx, tt.args.groupID, tt.args.roleID)
			require.Equal(t, tt.want.err, err)
		})
	}
}

func testAdminGroupRole(groupID, roleID, adminID string, now time.Time) *entity.AdminGroupRole {
	return &entity.AdminGroupRole{
		GroupID:        groupID,
		RoleID:         roleID,
		CreatedAdminID: adminID,
		UpdatedAdminID: adminID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
