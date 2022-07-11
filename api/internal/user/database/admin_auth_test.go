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
	auth := make(entity.AdminAuths, 1)
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

func testAdminAuth(adminID, cognitoID string, now time.Time) *entity.AdminAuth {
	return &entity.AdminAuth{
		AdminID:   adminID,
		CognitoID: cognitoID,
		Role:      entity.AdminRoleAdministrator,
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
