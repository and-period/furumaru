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

func TestAdmin(t *testing.T) {
	assert.NotNil(t, NewAdmin(nil))
}

func TestAdmin_List(t *testing.T) {
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

	_ = m.dbDelete(ctx, adminTable)
	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "test-admin02@and-period.jp", now())
	err = m.db.DB.Create(&admins).Error
	require.NoError(t, err)

	type args struct {
		params *ListAdminsParams
	}
	type want struct {
		admins entity.Admins
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
				params: &ListAdminsParams{
					Roles: []entity.AdminRole{
						entity.AdminRoleAdministrator,
						entity.AdminRoleProducer,
					},
					Limit:  1,
					Offset: 1,
				},
			},
			want: want{
				admins: admins[1:],
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

			db := &admin{db: m.db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreAdminsField(actual, now())
			assert.ElementsMatch(t, tt.want.admins, actual)
		})
	}
}

func TestAdmin_MultiGet(t *testing.T) {
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

	_ = m.dbDelete(ctx, adminTable)
	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "test-admin02@and-period.jp", now())
	err = m.db.DB.Create(&admins).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, m)

			db := &admin{db: m.db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.adminIDs)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreAdminsField(actual, now())
			assert.ElementsMatch(t, tt.want.admins, actual)
		})
	}
}

func TestAdmin_Get(t *testing.T) {
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

	_ = m.dbDelete(ctx, adminTable)
	a := testAdmin("admin-id", "test-admin@and-period.jp", now())
	err = m.db.DB.Create(&a).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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

			tt.setup(ctx, t, m)

			db := &admin{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.adminID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreAdminField(actual, now())
			assert.Equal(t, tt.want.admin, actual)
		})
	}
}

func TestAdmin_GetByCognitoID(t *testing.T) {
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

	_ = m.dbDelete(ctx, adminTable)
	a := testAdmin("admin-id", "test-admin@and-period.jp", now())
	err = m.db.DB.Create(&a).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				cognitoID: "admin-id",
			},
			want: want{
				admin:  a,
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

			tt.setup(ctx, t, m)

			db := &admin{db: m.db, now: now}
			actual, err := db.GetByCognitoID(ctx, tt.args.cognitoID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreAdminField(actual, now())
			assert.Equal(t, tt.want.admin, actual)
		})
	}
}

func TestAdmin_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	type args struct {
		admin *entity.Admin
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
				admin: testAdmin("admin-id", "test-admin@and-period.jp", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				admin := testAdmin("admin-id", "test-admin@and-period.jp", now())
				err = m.db.DB.Create(&admin).Error
				require.NoError(t, err)
			},
			args: args{
				admin: testAdmin("admin-id", "test-admin@and-period.jp", now()),
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

			err := m.dbDelete(ctx, adminTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &admin{db: m.db, now: now}
			err = db.Create(ctx, tt.args.admin)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestAdmin_UpdateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	type args struct {
		adminID string
		email   string
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
				a := testAdmin("admin-id", "test-admin@and-period.jp", now())
				err = m.db.DB.Create(&a).Error
				require.NoError(t, err)
			},
			args: args{
				adminID: "admin-id",
				email:   "test-other@and-period.jp",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				adminID: "admin-id",
				email:   "test-other@and-period.jp",
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

			err := m.dbDelete(ctx, adminTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &admin{db: m.db, now: now}
			err = db.UpdateEmail(ctx, tt.args.adminID, tt.args.email)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testAdmin(id, email string, now time.Time) *entity.Admin {
	return &entity.Admin{
		ID:            id,
		CognitoID:     id,
		Email:         email,
		Lastname:      "&.",
		Firstname:     "スタッフ",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "すたっふ",
		ThumbnailURL:  "https://and-period.jp/thumbnail.png",
		Role:          entity.AdminRoleAdministrator,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func fillIgnoreAdminField(a *entity.Admin, now time.Time) {
	if a == nil {
		return
	}
	a.CreatedAt = now
	a.UpdatedAt = now
}

func fillIgnoreAdminsField(as entity.Admins, now time.Time) {
	for i := range as {
		fillIgnoreAdminField(as[i], now)
	}
}
