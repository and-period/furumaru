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

func TestAdminAuth(t *testing.T) {
	assert.NotNil(t, NewAdminAuth(nil))
}

func TestAdminAuth_MultiGet(t *testing.T) {
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

	_ = m.dbDelete(ctx, adminAuthTable)
	auths := make(entity.AdminAuths, 2)
	auths[0] = testAdminAuth("admin-id01", "cognito-id01", now())
	auths[1] = testAdminAuth("admin-id02", "cognito-id02", now())
	err = m.db.DB.Create(&auths).Error
	require.NoError(t, err)

	type args struct {
		adminIDs []string
	}
	type want struct {
		auths  entity.AdminAuths
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
				auths:  auths,
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

			db := &adminAuth{db: m.db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.adminIDs)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreAdminAuthsField(actual, now())
			assert.Equal(t, tt.want.auths, actual)
		})
	}
}

func TestAdminAuth_GetByAdminID(t *testing.T) {
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

	_ = m.dbDelete(ctx, adminAuthTable)
	a := testAdminAuth("admin-id", "cognito-id", now())
	err = m.db.DB.Create(&a).Error
	require.NoError(t, err)

	type args struct {
		adminID string
	}
	type want struct {
		auth   *entity.AdminAuth
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
				auth:   a,
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
				auth:   nil,
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

			db := &adminAuth{db: m.db, now: now}
			actual, err := db.GetByAdminID(ctx, tt.args.adminID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreAdminAuthField(actual, now())
			assert.Equal(t, tt.want.auth, actual)
		})
	}
}

func TestAdminAuth_GetByCognitoID(t *testing.T) {
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

	_ = m.dbDelete(ctx, adminAuthTable)
	a := testAdminAuth("admin-id", "cognito-id", now())
	err = m.db.DB.Create(&a).Error
	require.NoError(t, err)

	type args struct {
		cognitoID string
	}
	type want struct {
		auth   *entity.AdminAuth
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
				cognitoID: "cognito-id",
			},
			want: want{
				auth:   a,
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
				auth:   nil,
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

			db := &adminAuth{db: m.db, now: now}
			actual, err := db.GetByCognitoID(ctx, tt.args.cognitoID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreAdminAuthField(actual, now())
			assert.Equal(t, tt.want.auth, actual)
		})
	}
}

func TestAdminAuth_UpdateDevice(t *testing.T) {
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

	_ = m.dbDelete(ctx, adminAuthTable)

	type args struct {
		adminID string
		device  string
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
				a := testAdminAuth("admin-id", "cognito-id", now())
				err = m.db.DB.Create(&a).Error
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
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				adminID: "admin-id",
				device:  "device",
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

			err := m.dbDelete(ctx, adminAuthTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &adminAuth{db: m.db, now: now}
			err = db.UpdateDevice(ctx, tt.args.adminID, tt.args.device)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testAdminAuth(adminID, cognitoID string, now time.Time) *entity.AdminAuth {
	return &entity.AdminAuth{
		AdminID:   adminID,
		CognitoID: cognitoID,
		Role:      entity.AdminRoleAdministrator,
		Device:    "instance-id",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func fillIgnoreAdminAuthField(a *entity.AdminAuth, now time.Time) {
	if a == nil {
		return
	}
	a.CreatedAt = now
	a.UpdatedAt = now
}

func fillIgnoreAdminAuthsField(as entity.AdminAuths, now time.Time) {
	for i := range as {
		fillIgnoreAdminAuthField(as[i], now)
	}
}
