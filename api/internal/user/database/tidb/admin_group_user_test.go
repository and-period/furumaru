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

func TestAdminGroupUser_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin1@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin2@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	require.NoError(t, err)

	g := testAdminGroup("group-id", "admin-id01", now())
	err = db.DB.WithContext(ctx).Table(adminGroupTable).Create(&g).Error
	require.NoError(t, err)

	role := testAdminRole("role-id", now())
	err = db.DB.WithContext(ctx).Table(adminRoleTable).Create(&role).Error
	require.NoError(t, err)

	gr := testAdminGroupRole("group-id", "role-id", "admin-id01", now())
	err = db.DB.WithContext(ctx).Table(adminGroupRoleTable).Create(&gr).Error
	require.NoError(t, err)

	users := make(entity.AdminGroupUsers, 2)
	users[0] = testAdminGroupUser(
		"group-id",
		"admin-id01",
		now().Add(time.Hour),
		now().Add(-time.Hour),
	)
	users[1] = testAdminGroupUser("group-id", "admin-id02", now().Add(time.Hour), now())
	err = db.DB.WithContext(ctx).Table(adminGroupUserTable).Create(&users).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdminGroupUsersParams
	}
	type want struct {
		groupUsers entity.AdminGroupUsers
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
				params: &database.ListAdminGroupUsersParams{},
			},
			want: want{
				groupUsers: users,
				err:        nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &adminGroupUser{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.groupUsers, actual)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdminGroupUser_Count(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin1@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin2@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	require.NoError(t, err)

	g := testAdminGroup("group-id", "admin-id01", now())
	err = db.DB.WithContext(ctx).Table(adminGroupTable).Create(&g).Error
	require.NoError(t, err)

	role := testAdminRole("role-id", now())
	err = db.DB.WithContext(ctx).Table(adminRoleTable).Create(&role).Error
	require.NoError(t, err)

	gr := testAdminGroupRole("group-id", "role-id", "admin-id01", now())
	err = db.DB.WithContext(ctx).Table(adminGroupRoleTable).Create(&gr).Error
	require.NoError(t, err)

	users := make(entity.AdminGroupUsers, 2)
	users[0] = testAdminGroupUser(
		"group-id",
		"admin-id01",
		now().Add(time.Hour),
		now().Add(-time.Hour),
	)
	users[1] = testAdminGroupUser("group-id", "admin-id02", now().Add(time.Hour), now())
	err = db.DB.WithContext(ctx).Table(adminGroupUserTable).Create(&users).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdminGroupUsersParams
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
				params: &database.ListAdminGroupUsersParams{},
			},
			want: want{
				total: 2,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &adminGroupUser{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.total, actual)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdminGroupUser_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
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

	u := testAdminGroupUser("group-id", "admin-id", now().Add(time.Hour), now())
	err = db.DB.WithContext(ctx).Table(adminGroupUserTable).Create(&u).Error
	require.NoError(t, err)

	type args struct {
		groupID string
		adminID string
	}
	type want struct {
		groupUser *entity.AdminGroupUser
		err       error
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
				adminID: "admin-id",
			},
			want: want{
				groupUser: u,
				err:       nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				groupID: "",
				adminID: "",
			},
			want: want{
				groupUser: nil,
				err:       database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &adminGroupUser{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.groupID, tt.args.adminID)
			assert.Equal(t, tt.want.groupUser, actual)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdminGroupUser_Upsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	admin := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
	err = db.DB.WithContext(ctx).Table(adminTable).Create(&admin).Error
	require.NoError(t, err)

	group := testAdminGroup("group-id", "admin-id", now())
	err = db.DB.WithContext(ctx).Table(adminGroupTable).Create(&group).Error
	require.NoError(t, err)

	type args struct {
		user *entity.AdminGroupUser
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
				user: testAdminGroupUser("group-id", "admin-id", now().Add(time.Hour), now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success to udpate",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				u := testAdminGroupUser("group-id", "admin-id", now(), now())
				err := db.DB.WithContext(ctx).Table(adminGroupUserTable).Create(&u).Error
				require.NoError(t, err)
			},
			args: args{
				user: testAdminGroupUser("group-id", "admin-id", now().Add(time.Hour), now()),
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, adminGroupUserTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &adminGroupUser{db: db, now: now}
			err = db.Upsert(ctx, tt.args.user)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdminGroupUser_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	admin := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
	err = db.DB.WithContext(ctx).Table(adminTable).Create(&admin).Error
	require.NoError(t, err)

	group := testAdminGroup("group-id", "admin-id", now())
	err = db.DB.WithContext(ctx).Table(adminGroupTable).Create(&group).Error
	require.NoError(t, err)

	type args struct {
		groupID string
		adminID string
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
				u := testAdminGroupUser("group-id", "admin-id", now().Add(time.Hour), now())
				err := db.DB.WithContext(ctx).Table(adminGroupUserTable).Create(&u).Error
				require.NoError(t, err)
			},
			args: args{
				groupID: "group-id",
				adminID: "admin-id",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, adminGroupUserTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &adminGroupUser{db: db, now: now}
			err = db.Delete(ctx, tt.args.groupID, tt.args.adminID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testAdminGroupUser(groupID, adminID string, expiredAt, now time.Time) *entity.AdminGroupUser {
	return &entity.AdminGroupUser{
		GroupID:        groupID,
		AdminID:        adminID,
		CreatedAdminID: adminID,
		UpdatedAdminID: adminID,
		ExpiredAt:      expiredAt,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
