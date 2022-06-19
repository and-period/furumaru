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

func TestCoordinator(t *testing.T) {
	assert.NotNil(t, NewCoordinator(nil))
}

func TestCoordinator_List(t *testing.T) {
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

	_ = m.dbDelete(ctx, coordinatorTable)
	admins := make(entity.Coordinators, 2)
	admins[0] = testCoordinator("admin-id01", "test-admin01@and-period.jp", now())
	admins[1] = testCoordinator("admin-id02", "test-admin02@and-period.jp", now())
	err = m.db.DB.Create(&admins).Error
	require.NoError(t, err)

	type args struct {
		params *ListCoordinatorsParams
	}
	type want struct {
		admins entity.Coordinators
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
				params: &ListCoordinatorsParams{
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

			db := &coordinator{db: m.db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreCoordinatorsField(actual, now())
			assert.Equal(t, tt.want.admins, actual)
		})
	}
}

func TestCoordinator_Get(t *testing.T) {
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

	_ = m.dbDelete(ctx, coordinatorTable)
	a := testCoordinator("admin-id", "test-admin@and-period.jp", now())
	err = m.db.DB.Create(&a).Error
	require.NoError(t, err)

	type args struct {
		adminID string
	}
	type want struct {
		admin  *entity.Coordinator
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

			db := &coordinator{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.adminID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreCoordinatorField(actual, now())
			assert.Equal(t, tt.want.admin, actual)
		})
	}
}

func TestCoordinator_Create(t *testing.T) {
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
		admin *entity.Coordinator
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
				admin: testCoordinator("admin-id", "test-admin@and-period.jp", now()),
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
				admin: testCoordinator("admin-id", "test-admin@and-period.jp", now()),
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to duplicate entry in coordinator",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				admin := testCoordinator("admin-id", "test-admin@and-period.jp", now())
				err = m.db.DB.Create(&admin).Error
				require.NoError(t, err)
			},
			args: args{
				auth:  testAdminAuth("admin-id", "cognito-id", now()),
				admin: testCoordinator("admin-id", "test-admin@and-period.jp", now()),
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

			err := m.dbDelete(ctx, adminAuthTable, coordinatorTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &coordinator{db: m.db, now: now}
			err = db.Create(ctx, tt.args.auth, tt.args.admin)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestCoordinator_UpdateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	type args struct {
		coordinatorID string
		email         string
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
				p := testCoordinator("admin-id", "test-admin@and-period.jp", now())
				err = m.db.DB.Create(&p).Error
				require.NoError(t, err)
			},
			args: args{
				coordinatorID: "admin-id",
				email:         "test-other@and-period.jp",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				coordinatorID: "admin-id",
				email:         "test-other@and-period.jp",
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

			err := m.dbDelete(ctx, coordinatorTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &coordinator{db: m.db, now: now}
			err = db.UpdateEmail(ctx, tt.args.coordinatorID, tt.args.email)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testCoordinator(id, email string, now time.Time) *entity.Coordinator {
	return &entity.Coordinator{
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

func fillIgnoreCoordinatorField(a *entity.Coordinator, now time.Time) {
	if a == nil {
		return
	}
	a.CreatedAt = now
	a.UpdatedAt = now
}

func fillIgnoreCoordinatorsField(as entity.Coordinators, now time.Time) {
	for i := range as {
		fillIgnoreCoordinatorField(as[i], now)
	}
}
