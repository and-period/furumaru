package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCoordinator(t *testing.T) {
	assert.NotNil(t, NewCoordinator(nil))
}

func TestCoordinator_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now().Add(time.Hour))
	err = db.DB.Debug().Create(&admins).Error
	require.NoError(t, err)
	coordinators := make(entity.Coordinators, 2)
	coordinators[0] = testCoordinator("admin-id01", now())
	coordinators[0].Admin = *admins[0]
	coordinators[1] = testCoordinator("admin-id02", now())
	coordinators[1].Admin = *admins[1]
	err = db.DB.Table(coordinatorTable).Create(&coordinators).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListCoordinatorsParams
	}
	type want struct {
		coordinators entity.Coordinators
		hasErr       bool
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
				params: &database.ListCoordinatorsParams{
					Name:   "農園",
					Limit:  1,
					Offset: 1,
				},
			},
			want: want{
				coordinators: coordinators[:1],
				hasErr:       false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &coordinator{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.coordinators, actual)
		})
	}
}

func TestCoordinator_Count(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	require.NoError(t, err)
	coordinators := make(entity.Coordinators, 2)
	coordinators[0] = testCoordinator("admin-id01", now())
	coordinators[0].Admin = *admins[0]
	coordinators[1] = testCoordinator("admin-id02", now())
	coordinators[1].Admin = *admins[1]
	err = db.DB.Table(coordinatorTable).Create(&coordinators).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListCoordinatorsParams
	}
	type want struct {
		total  int64
		hasErr bool
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
				params: &database.ListCoordinatorsParams{
					Name:   "農園",
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &coordinator{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestCoordinator_MultiGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	require.NoError(t, err)
	coordinators := make(entity.Coordinators, 2)
	coordinators[0] = testCoordinator("admin-id01", now())
	coordinators[0].Admin = *admins[0]
	coordinators[1] = testCoordinator("admin-id02", now())
	coordinators[1].Admin = *admins[1]
	err = db.DB.Table(coordinatorTable).Create(&coordinators).Error
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
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &coordinator{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.adminIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.coordinators, actual)
		})
	}
}

func TestCoordinator_MultiGetWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin01@and-period.jp", now())
	admins[0].Status = entity.AdminStatusDeactivated
	admins[0].DeletedAt = gorm.DeletedAt{Valid: true, Time: now()}
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	require.NoError(t, err)
	coordinators := make(entity.Coordinators, 2)
	coordinators[0] = testCoordinator("admin-id01", now())
	coordinators[0].Admin = *admins[0]
	coordinators[1] = testCoordinator("admin-id02", now())
	coordinators[1].Admin = *admins[1]
	err = db.DB.Table(coordinatorTable).Create(&coordinators).Error
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
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &coordinator{db: db, now: now}
			actual, err := db.MultiGetWithDeleted(ctx, tt.args.adminIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.coordinators, actual)
		})
	}
}

func TestCoordinator_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
	err = db.DB.Create(&admin).Error
	require.NoError(t, err)
	c := testCoordinator("admin-id", now())
	c.Admin = *admin
	err = db.DB.Table(coordinatorTable).Create(&c).Error
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
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
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
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &coordinator{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.adminID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.coordinator, actual)
		})
	}
}

func TestCoordinator_GetWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin01@and-period.jp", now())
	admins[0].Status = entity.AdminStatusDeactivated
	admins[0].DeletedAt = gorm.DeletedAt{Valid: true, Time: now()}
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	require.NoError(t, err)
	coordinators := make(entity.Coordinators, 2)
	coordinators[0] = testCoordinator("admin-id01", now())
	coordinators[0].Admin = *admins[0]
	coordinators[1] = testCoordinator("admin-id02", now())
	coordinators[1].Admin = *admins[1]
	err = db.DB.Table(coordinatorTable).Create(&coordinators).Error
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
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success to activated",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				adminID: "admin-id01",
			},
			want: want{
				coordinator: coordinators[0],
				hasErr:      false,
			},
		},
		{
			name:  "success to deactivated",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				adminID: "admin-id02",
			},
			want: want{
				coordinator: coordinators[1],
				hasErr:      false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &coordinator{db: db, now: now}
			actual, err := db.GetWithDeleted(ctx, tt.args.adminID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.coordinator, actual)
		})
	}
}

func TestCoordinator_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	admin := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
	c := testCoordinator("admin-id", now())
	c.Admin = *admin

	ishop := testShop("shop-id", "admin-id", []string{}, []string{}, now())
	shop := &ishop.Shop

	type args struct {
		coordinator *entity.Coordinator
		shop        *entity.Shop
		auth        func(ctx context.Context) error
	}
	type want struct {
		hasErr bool
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
				coordinator: c,
				shop:        shop,
				auth:        func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry in admin auth",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				internal := testCoordinator("admin-id", now())
				err = db.DB.Table(coordinatorTable).Create(&internal).Error
				require.NoError(t, err)
			},
			args: args{
				coordinator: c,
				shop:        shop,
				auth:        func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name:  "failed to execute external service",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				coordinator: c,
				shop:        shop,
				auth:        func(ctx context.Context) error { return assert.AnError },
			},
			want: want{
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, shopTable, coordinatorTable, adminTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &coordinator{db: db, now: now}
			err = db.Create(ctx, tt.args.coordinator, tt.args.shop, tt.args.auth)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestCoordinator_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	type args struct {
		coordinatorID string
		params        *database.UpdateCoordinatorParams
	}
	type want struct {
		hasErr bool
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
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				coordinator := testCoordinator("admin-id", now())
				err = db.DB.Table(coordinatorTable).Create(&coordinator).Error
				require.NoError(t, err)
			},
			args: args{
				coordinatorID: "admin-id",
				params: &database.UpdateCoordinatorParams{
					Lastname:          "&.",
					Firstname:         "スタッフ",
					LastnameKana:      "あんどぴりおど",
					FirstnameKana:     "すたっふ",
					PhoneNumber:       "+819012345678",
					Username:          "&.農園",
					Profile:           "紹介文です。",
					ThumbnailURL:      "https://and-period.jp/thumbnail.png",
					HeaderURL:         "https://and-period.jp/header.png",
					PromotionVideoURL: "https://and-period.jp/promotion.mp4",
					BonusVideoURL:     "https://and-period.jp/bonus.mp4",
					InstagramID:       "instagram-id",
					FacebookID:        "facebook-id",
					PostalCode:        "1000014",
					PrefectureCode:    codes.PrefectureValues["tokyo"],
					City:              "千代田区",
					AddressLine1:      "永田町1-7-1",
					AddressLine2:      "",
				},
			},
			want: want{
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, coordinatorTable, adminTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &coordinator{db: db, now: now}
			err = db.Update(ctx, tt.args.coordinatorID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestCoordinator_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	type args struct {
		coordinatorID string
		auth          func(ctx context.Context) error
	}
	type want struct {
		hasErr bool
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
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				coordinator := testCoordinator("admin-id", now())
				err = db.DB.Table(coordinatorTable).Create(&coordinator).Error
				require.NoError(t, err)
			},
			args: args{
				coordinatorID: "admin-id",
				auth:          func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to execute external service",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				coordinator := testCoordinator("admin-id", now())
				err = db.DB.Table(coordinatorTable).Create(&coordinator).Error
				require.NoError(t, err)
			},
			args: args{
				coordinatorID: "admin-id",
				auth:          func(ctx context.Context) error { return assert.AnError },
			},
			want: want{
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, coordinatorTable, adminTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &coordinator{db: db, now: now}
			err = db.Delete(ctx, tt.args.coordinatorID, tt.args.auth)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testCoordinator(id string, now time.Time) *entity.Coordinator {
	return &entity.Coordinator{
		AdminID:           id,
		PhoneNumber:       "+819012345678",
		Username:          "&.農園",
		Profile:           "紹介文です。",
		ThumbnailURL:      "https://and-period.jp/thumbnail.png",
		HeaderURL:         "https://and-period.jp/header.png",
		PromotionVideoURL: "https://and-period.jp/promotion.mp4",
		BonusVideoURL:     "https://and-period.jp/bonus.mp4",
		InstagramID:       "instagram-id",
		FacebookID:        "facebook-id",
		PostalCode:        "1000014",
		Prefecture:        "東京都",
		PrefectureCode:    13,
		City:              "千代田区",
		AddressLine1:      "永田町1-7-1",
		AddressLine2:      "",
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}
