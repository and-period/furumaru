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

	_ = m.dbDelete(ctx, coordinatorTable, adminTable)
	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = m.db.DB.Debug().Create(&admins).Error
	coordinators := make(entity.Coordinators, 2)
	coordinators[0] = testCoordinator("admin-id01", now())
	coordinators[0].Admin = *admins[0]
	coordinators[1] = testCoordinator("admin-id02", now())
	coordinators[1].Admin = *admins[1]
	err = m.db.DB.Debug().Create(&coordinators).Error
	require.NoError(t, err)

	type args struct {
		params *ListCoordinatorsParams
	}
	type want struct {
		coordinators entity.Coordinators
		hasErr       bool
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
				coordinators: coordinators[1:],
				hasErr:       false,
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
			assert.Equal(t, tt.want.coordinators, actual)
		})
	}
}

func TestCoordinator_Count(t *testing.T) {
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

	_ = m.dbDelete(ctx, coordinatorTable, adminTable)
	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = m.db.DB.Create(&admins).Error
	coordinators := make(entity.Coordinators, 2)
	coordinators[0] = testCoordinator("admin-id01", now())
	coordinators[0].Admin = *admins[0]
	coordinators[1] = testCoordinator("admin-id02", now())
	coordinators[1].Admin = *admins[1]
	err = m.db.DB.Create(&coordinators).Error
	require.NoError(t, err)

	type args struct {
		params *ListCoordinatorsParams
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
				params: &ListCoordinatorsParams{
					Limit:  1,
					Offset: 1,
				},
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

			db := &coordinator{db: m.db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestCoordinator_MultiGet(t *testing.T) {
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

	_ = m.dbDelete(ctx, coordinatorTable, adminTable)
	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = m.db.DB.Create(&admins).Error
	coordinators := make(entity.Coordinators, 2)
	coordinators[0] = testCoordinator("admin-id01", now())
	coordinators[0].Admin = *admins[0]
	coordinators[1] = testCoordinator("admin-id02", now())
	coordinators[1].Admin = *admins[1]
	err = m.db.DB.Create(&coordinators).Error
	require.NoError(t, err)

	type args struct {
		adminIDs []string
	}
	type want struct {
		coordinators entity.Coordinators
		hasErr       bool
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
				coordinators: coordinators,
				hasErr:       false,
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
			actual, err := db.MultiGet(ctx, tt.args.adminIDs)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreCoordinatorsField(actual, now())
			assert.Equal(t, tt.want.coordinators, actual)
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

	_ = m.dbDelete(ctx, coordinatorTable, adminTable)
	admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
	err = m.db.DB.Create(&admin).Error
	require.NoError(t, err)
	c := testCoordinator("admin-id", now())
	c.Admin = *admin
	err = m.db.DB.Create(&c).Error
	require.NoError(t, err)

	type args struct {
		adminID string
	}
	type want struct {
		coordinator *entity.Coordinator
		hasErr      bool
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
				coordinator: c,
				hasErr:      false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				adminID: "",
			},
			want: want{
				coordinator: nil,
				hasErr:      true,
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
			assert.Equal(t, tt.want.coordinator, actual)
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
		admin       *entity.Admin
		coordinator *entity.Coordinator
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
				admin:       testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now()),
				coordinator: testCoordinator("admin-id", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry in admin auth",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = m.db.DB.Create(&admin).Error
				require.NoError(t, err)
				coordinator := testCoordinator("admin-id", now())
				err = m.db.DB.Create(&coordinator).Error
				require.NoError(t, err)
			},
			args: args{
				admin:       testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now()),
				coordinator: testCoordinator("admin-id", now()),
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

			err := m.dbDelete(ctx, coordinatorTable, adminTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &coordinator{db: m.db, now: now}
			err = db.Create(ctx, tt.args.admin, tt.args.coordinator)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestCoordinator_Update(t *testing.T) {
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
		params        *UpdateCoordinatorParams
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
				_ = m.dbDelete(ctx, coordinatorTable, adminTable)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = m.db.DB.Create(&admin).Error
				require.NoError(t, err)
				coordinator := testCoordinator("admin-id", now())
				err = m.db.DB.Create(&coordinator).Error
				require.NoError(t, err)
			},
			args: args{
				coordinatorID: "admin-id",
				params: &UpdateCoordinatorParams{
					Lastname:         "&.",
					Firstname:        "スタッフ",
					LastnameKana:     "あんどぴりおど",
					FirstnameKana:    "すたっふ",
					CompanyName:      "&.株式会社",
					ThumbnailURL:     "https://and-period.jp/thumbnail.png",
					HeaderURL:        "https://and-period.jp/header.png",
					TwitterAccount:   "twitter-id",
					InstagramAccount: "instagram-id",
					FacebookAccount:  "facebook-id",
					PhoneNumber:      "+819012345678",
					PostalCode:       "1000014",
					Prefecture:       "東京都",
					City:             "千代田区",
					AddressLine1:     "永田町1-7-1",
					AddressLine2:     "",
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				coordinatorID: "admin-id",
				params:        &UpdateCoordinatorParams{},
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

			err := m.dbDelete(ctx, coordinatorTable, adminTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &coordinator{db: m.db, now: now}
			err = db.Update(ctx, tt.args.coordinatorID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testCoordinator(id string, now time.Time) *entity.Coordinator {
	return &entity.Coordinator{
		AdminID:          id,
		PhoneNumber:      "+819012345678",
		CompanyName:      "&.株式会社",
		StoreName:        "&.農園",
		ThumbnailURL:     "https://and-period.jp/thumbnail.png",
		HeaderURL:        "https://and-period.jp/header.png",
		TwitterAccount:   "twitter-id",
		InstagramAccount: "instagram-id",
		FacebookAccount:  "facebook-id",
		PostalCode:       "1000014",
		Prefecture:       "東京都",
		City:             "千代田区",
		AddressLine1:     "永田町1-7-1",
		AddressLine2:     "",
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func fillIgnoreCoordinatorField(a *entity.Coordinator, now time.Time) {
	if a == nil {
		return
	}
	a.CreatedAt = now
	a.UpdatedAt = now
	a.Admin.CreatedAt = now
	a.Admin.UpdatedAt = now
}

func fillIgnoreCoordinatorsField(as entity.Coordinators, now time.Time) {
	for i := range as {
		fillIgnoreCoordinatorField(as[i], now)
	}
}
