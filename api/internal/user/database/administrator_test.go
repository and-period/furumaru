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

func TestAdministrator(t *testing.T) {
	assert.NotNil(t, NewAdministrator(nil))
}

func TestAdministrator_List(t *testing.T) {
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

	_ = m.dbDelete(ctx, administratorTable)
	admins := make(entity.Administrators, 2)
	admins[0] = testAdministrator("admin-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdministrator("admin-id02", "test-admin02@and-period.jp", now())
	err = m.db.DB.Create(&admins).Error
	require.NoError(t, err)

	type args struct {
		params *ListAdministratorsParams
	}
	type want struct {
		admins entity.Administrators
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
				params: &ListAdministratorsParams{
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

			db := &administrator{db: m.db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreAdministratorsField(actual, now())
			assert.Equal(t, tt.want.admins, actual)
		})
	}
}

func TestAdministrator_Count(t *testing.T) {
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

	_ = m.dbDelete(ctx, administratorTable)
	admins := make(entity.Administrators, 2)
	admins[0] = testAdministrator("admin-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdministrator("admin-id02", "test-admin02@and-period.jp", now())
	err = m.db.DB.Create(&admins).Error
	require.NoError(t, err)

	type args struct {
		params *ListAdministratorsParams
	}
	type want struct {
		total  int64
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
				params: &ListAdministratorsParams{},
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

			tt.setup(ctx, t, m)

			db := &administrator{db: m.db, now: now}
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

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, administratorTable)
	admins := make(entity.Administrators, 2)
	admins[0] = testAdministrator("admin-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdministrator("admin-id02", "test-admin02@and-period.jp", now())
	err = m.db.DB.Create(&admins).Error
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

			db := &administrator{db: m.db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.adminIDs)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreAdministratorsField(actual, now())
			assert.Equal(t, tt.want.admins, actual)
		})
	}
}

func TestAdministrator_Get(t *testing.T) {
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

	_ = m.dbDelete(ctx, administratorTable)
	a := testAdministrator("admin-id", "test-admin@and-period.jp", now())
	err = m.db.DB.Create(&a).Error
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

			db := &administrator{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.adminID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreAdministratorField(actual, now())
			assert.Equal(t, tt.want.admin, actual)
		})
	}
}

func TestAdministrator_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	type args struct {
		auth  *entity.AdminAuth
		admin *entity.Administrator
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
				auth:  testAdminAuth("admin-id", "cognito-id", now()),
				admin: testAdministrator("admin-id", "test-admin@and-period.jp", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry in admin auth",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				auth := testAdminAuth("admin-id", "cognito-id", now())
				err = m.db.DB.Create(&auth).Error
				require.NoError(t, err)
			},
			args: args{
				auth:  testAdminAuth("admin-id", "cognito-id", now()),
				admin: testAdministrator("admin-id", "test-admin@and-period.jp", now()),
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to duplicate entry in administrator",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				admin := testAdministrator("admin-id", "test-admin@and-period.jp", now())
				err = m.db.DB.Create(&admin).Error
				require.NoError(t, err)
			},
			args: args{
				auth:  testAdminAuth("admin-id", "cognito-id", now()),
				admin: testAdministrator("admin-id", "test-admin@and-period.jp", now()),
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

			err := m.dbDelete(ctx, adminAuthTable, administratorTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &administrator{db: m.db, now: now}
			err = db.Create(ctx, tt.args.auth, tt.args.admin)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestAdministrator_UpdateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	type args struct {
		administratorID string
		email           string
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
				p := testAdministrator("admin-id", "test-admin@and-period.jp", now())
				err = m.db.DB.Create(&p).Error
				require.NoError(t, err)
			},
			args: args{
				administratorID: "admin-id",
				email:           "test-other@and-period.jp",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				administratorID: "admin-id",
				email:           "test-other@and-period.jp",
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

			err := m.dbDelete(ctx, administratorTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &administrator{db: m.db, now: now}
			err = db.UpdateEmail(ctx, tt.args.administratorID, tt.args.email)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testAdministrator(id, email string, now time.Time) *entity.Administrator {
	return &entity.Administrator{
		ID:            id,
		Lastname:      "&.",
		Firstname:     "スタッフ",
		LastnameKana:  "あんどぴりおど",
		FirstnameKana: "すたっふ",
		Email:         email,
		PhoneNumber:   "+819012345678",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func fillIgnoreAdministratorField(a *entity.Administrator, now time.Time) {
	if a == nil {
		return
	}
	a.CreatedAt = now
	a.UpdatedAt = now
}

func fillIgnoreAdministratorsField(as entity.Administrators, now time.Time) {
	for i := range as {
		fillIgnoreAdministratorField(as[i], now)
	}
}
