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

func TestAdminGroup_List(t *testing.T) {
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

	groups := make(entity.AdminGroups, 3)
	groups[0] = testAdminGroup("group-id01", "admin-id", now().Add(-2*time.Hour))
	groups[1] = testAdminGroup("group-id02", "admin-id", now().Add(-time.Hour))
	groups[2] = testAdminGroup("group-id03", "admin-id", now())
	err = db.DB.WithContext(ctx).Table(adminGroupTable).Create(&groups).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdminGroupsParams
	}
	type want struct {
		groups entity.AdminGroups
		err    error
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
				params: &database.ListAdminGroupsParams{},
			},
			want: want{
				groups: groups,
				err:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &adminGroup{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.groups, actual)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdminGroup_Count(t *testing.T) {
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

	groups := make(entity.AdminGroups, 3)
	groups[0] = testAdminGroup("group-id01", "admin-id", now().Add(-2*time.Hour))
	groups[1] = testAdminGroup("group-id02", "admin-id", now().Add(-time.Hour))
	groups[2] = testAdminGroup("group-id03", "admin-id", now())
	err = db.DB.WithContext(ctx).Table(adminGroupTable).Create(&groups).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdminGroupsParams
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
				params: &database.ListAdminGroupsParams{},
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

			db := &adminGroup{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.total, actual)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdminGroup_MultiGet(t *testing.T) {
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

	groups := make(entity.AdminGroups, 3)
	groups[0] = testAdminGroup("group-id01", "admin-id", now().Add(-2*time.Hour))
	groups[1] = testAdminGroup("group-id02", "admin-id", now().Add(-time.Hour))
	groups[2] = testAdminGroup("group-id03", "admin-id", now())
	err = db.DB.WithContext(ctx).Table(adminGroupTable).Create(&groups).Error
	require.NoError(t, err)

	type args struct {
		groupIDs []string
	}
	type want struct {
		groups entity.AdminGroups
		err    error
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
				groupIDs: []string{"group-id01", "group-id02", "group-id03"},
			},
			want: want{
				groups: groups,
				err:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &adminGroup{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.groupIDs)
			assert.Equal(t, tt.want.groups, actual)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdminGroup_Get(t *testing.T) {
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

	type args struct {
		groupID string
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

			db := &adminGroup{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.groupID)
			assert.Equal(t, tt.want.group, actual)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdminGroup_Upsert(t *testing.T) {
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

	type args struct {
		group *entity.AdminGroup
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
				group: testAdminGroup("group-id", "admin-id", now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success to update",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				g := testAdminGroup("group-id", "admin-id", now())
				err := db.DB.WithContext(ctx).Table(adminGroupTable).Create(&g).Error
				require.NoError(t, err)
			},
			args: args{
				group: testAdminGroup("group-id", "admin-id", now()),
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

			err := delete(ctx, adminGroupTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &adminGroup{db: db, now: now}
			err = db.Upsert(ctx, tt.args.group)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdminGroup_Delete(t *testing.T) {
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

	type args struct {
		groupID string
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
				g := testAdminGroup("group-id", "admin-id", now())
				err = db.DB.WithContext(ctx).Table(adminGroupTable).Create(&g).Error
				require.NoError(t, err)
			},
			args: args{
				groupID: "group-id",
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

			err := delete(ctx, adminGroupTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &adminGroup{db: db, now: now}
			err = db.Delete(ctx, tt.args.groupID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testAdminGroup(groupID, adminID string, now time.Time) *entity.AdminGroup {
	return &entity.AdminGroup{
		ID:             groupID,
		Type:           entity.AdminTypeAdministrator,
		Name:           "管理者グループ",
		Description:    "管理者グループです。",
		CreatedAdminID: adminID,
		UpdatedAdminID: adminID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
