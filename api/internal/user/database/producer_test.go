package database

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func TestProducer(t *testing.T) {
	assert.NotNil(t, NewProducer(nil))
}

func TestProducer_List(t *testing.T) {
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

	coordinator := testCoordinator("coordinator-id", now())
	coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
	err = db.DB.Create(&coordinator.Admin).Error
	require.NoError(t, err)
	err = db.DB.Create(&coordinator).Error
	require.NoError(t, err)
	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	producers := make(entity.Producers, 2)
	require.NoError(t, err)
	producers[0] = testProducer("admin-id01", "coordinator-id", "&.農園", now())
	producers[0].Admin = *admins[0]
	producers[1] = testProducer("admin-id02", "coordinator-id", "&.水産", now())
	producers[1].Admin = *admins[1]
	err = db.DB.Create(&producers).Error
	require.NoError(t, err)

	type args struct {
		params *ListProducersParams
	}
	type want struct {
		producers entity.Producers
		hasErr    bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &ListProducersParams{
					CoordinatorID: "coordinator-id",
					Limit:         1,
					Offset:        1,
				},
			},
			want: want{
				producers: producers[1:],
				hasErr:    false,
			},
		},
		{
			name:  "success only unrelated",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &ListProducersParams{
					Limit:         1,
					Offset:        1,
					OnlyUnrelated: true,
				},
			},
			want: want{
				producers: entity.Producers{},
				hasErr:    false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &producer{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.producers, actual)
		})
	}
}

func TestProducer_Count(t *testing.T) {
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

	coordinator := testCoordinator("coordinator-id", now())
	coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
	err = db.DB.Create(&coordinator.Admin).Error
	require.NoError(t, err)
	err = db.DB.Create(&coordinator).Error
	require.NoError(t, err)
	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	producers := make(entity.Producers, 2)
	producers[0] = testProducer("admin-id01", "coordinator-id", "&.農園", now())
	producers[0].Admin = *admins[0]
	producers[1] = testProducer("admin-id02", "coordinator-id", "&.水産", now())
	producers[1].Admin = *admins[1]
	err = db.DB.Create(&producers).Error
	require.NoError(t, err)

	type args struct {
		params *ListProducersParams
	}
	type want struct {
		total  int64
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &ListProducersParams{},
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

			tt.setup(ctx, t, db)

			db := &producer{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestProducer_MultiGet(t *testing.T) {
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

	coordinator := testCoordinator("coordinator-id", now())
	coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
	err = db.DB.Create(&coordinator.Admin).Error
	require.NoError(t, err)
	err = db.DB.Create(&coordinator).Error
	require.NoError(t, err)
	admins := make(entity.Admins, 2)
	admins[0] = testAdmin("admin-id01", "cognito-id01", "test-admin01@and-period.jp", now())
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	producers := make(entity.Producers, 2)
	producers[0] = testProducer("admin-id01", "coordinator-id", "&.農園", now())
	producers[0].Admin = *admins[0]
	producers[1] = testProducer("admin-id02", "coordinator-id", "&.水産", now())
	producers[1].Admin = *admins[1]
	err = db.DB.Create(&producers).Error
	require.NoError(t, err)

	type args struct {
		producerIDs []string
	}
	type want struct {
		producers entity.Producers
		hasErr    bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				producerIDs: []string{"admin-id01", "admin-id02"},
			},
			want: want{
				producers: producers,
				hasErr:    false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &producer{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.producerIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.producers, actual)
		})
	}
}

func TestProducer_Get(t *testing.T) {
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

	coordinator := testCoordinator("coordinator-id", now())
	coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
	err = db.DB.Create(&coordinator.Admin).Error
	require.NoError(t, err)
	err = db.DB.Create(&coordinator).Error
	require.NoError(t, err)
	admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
	err = db.DB.Create(&admin).Error
	require.NoError(t, err)
	p := testProducer("admin-id", "coordinator-id", "&.農園", now())
	p.Admin = *admin
	err = db.DB.Create(&p).Error
	require.NoError(t, err)

	type args struct {
		producerID string
	}
	type want struct {
		producer *entity.Producer
		hasErr   bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				producerID: "admin-id",
			},
			want: want{
				producer: p,
				hasErr:   false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				producerID: "",
			},
			want: want{
				producer: nil,
				hasErr:   true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &producer{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.producerID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.producer, actual)
		})
	}
}

func TestProducer_Create(t *testing.T) {
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

	p := testProducer("admin-id", "coordinator-id", "&.農園", now())
	p.Admin = *testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())

	type args struct {
		producer *entity.Producer
		auth     func(ctx context.Context) error
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
			},
			args: args{
				producer: p,
				auth:     func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found coordinator",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				producer: p,
				auth:     func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to duplicate entry in admin auth",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				p := testProducer("admin-id", "coordinator-id", "&.農園", now())
				err = db.DB.Create(&p).Error
				require.NoError(t, err)
			},
			args: args{
				producer: p,
				auth:     func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to execute external service",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
			},
			args: args{
				producer: p,
				auth:     func(ctx context.Context) error { return assert.AnError },
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

			err := delete(ctx, producerTable, coordinatorTable, adminTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &producer{db: db, now: now}
			err = db.Create(ctx, tt.args.producer, tt.args.auth)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProducer_Update(t *testing.T) {
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

	type args struct {
		producerID string
		params     *UpdateProducerParams
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				p := testProducer("admin-id", "coordinator-id", "&.農園", now())
				err = db.DB.Create(&p).Error
				require.NoError(t, err)
			},
			args: args{
				producerID: "admin-id",
				params: &UpdateProducerParams{
					Lastname:      "&.",
					Firstname:     "スタッフ",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "すたっふ",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					HeaderURL:     "https://and-period.jp/header.png",
					PhoneNumber:   "+819012345678",
					PostalCode:    "1000014",
					Prefecture:    "東京都",
					City:          "千代田区",
					AddressLine1:  "永田町1-7-1",
					AddressLine2:  "",
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				producerID: "admin-id",
				params:     &UpdateProducerParams{},
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

			err := delete(ctx, producerTable, coordinatorTable, adminTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &producer{db: db, now: now}
			err = db.Update(ctx, tt.args.producerID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProducer_UpdateThumbnails(t *testing.T) {
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

	type args struct {
		producerID string
		thumbnails common.Images
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				p := testProducer("admin-id", "coordinator-id", "&.農園", now())
				err = db.DB.Create(&p).Error
				require.NoError(t, err)
			},
			args: args{
				producerID: "admin-id",
				thumbnails: common.Images{
					{
						Size: common.ImageSizeSmall,
						URL:  "https://and-period.jp/thumbnail_240.png",
					},
					{
						Size: common.ImageSizeMedium,
						URL:  "https://and-period.jp/thumbnail_675.png",
					},
					{
						Size: common.ImageSizeLarge,
						URL:  "https://and-period.jp/thumbnail_900.png",
					},
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				producerID: "admin-id",
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed precondition for thumbnail url is empty",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				p := testProducer("admin-id", "coordinator-id", "&.農園", now())
				p.ThumbnailURL = ""
				err = db.DB.Create(&p).Error
				require.NoError(t, err)
			},
			args: args{
				producerID: "admin-id",
				thumbnails: common.Images{
					{
						Size: common.ImageSizeSmall,
						URL:  "https://and-period.jp/thumbnail_240.png",
					},
					{
						Size: common.ImageSizeMedium,
						URL:  "https://and-period.jp/thumbnail_675.png",
					},
					{
						Size: common.ImageSizeLarge,
						URL:  "https://and-period.jp/thumbnail_900.png",
					},
				},
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

			err := delete(ctx, producerTable, adminTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &producer{db: db, now: now}
			err = db.UpdateThumbnails(ctx, tt.args.producerID, tt.args.thumbnails)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProducer_UpdateHeaders(t *testing.T) {
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

	type args struct {
		ProducerID string
		headers    common.Images
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				p := testProducer("admin-id", "coordinator-id", "&.農園", now())
				err = db.DB.Create(&p).Error
				require.NoError(t, err)
			},
			args: args{
				ProducerID: "admin-id",
				headers: common.Images{
					{
						Size: common.ImageSizeSmall,
						URL:  "https://and-period.jp/header_240.png",
					},
					{
						Size: common.ImageSizeMedium,
						URL:  "https://and-period.jp/header_675.png",
					},
					{
						Size: common.ImageSizeLarge,
						URL:  "https://and-period.jp/header_900.png",
					},
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				ProducerID: "admin-id",
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed precondition for header url is empty",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				p := testProducer("admin-id", "coordinator-id", "&.農園", now())
				p.HeaderURL = ""
				err = db.DB.Create(&p).Error
				require.NoError(t, err)
			},
			args: args{
				ProducerID: "admin-id",
				headers: common.Images{
					{
						Size: common.ImageSizeSmall,
						URL:  "https://and-period.jp/header_240.png",
					},
					{
						Size: common.ImageSizeMedium,
						URL:  "https://and-period.jp/header_675.png",
					},
					{
						Size: common.ImageSizeLarge,
						URL:  "https://and-period.jp/header_900.png",
					},
				},
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

			err := delete(ctx, producerTable, adminTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &producer{db: db, now: now}
			err = db.UpdateHeaders(ctx, tt.args.ProducerID, tt.args.headers)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProducer_UpdateRelationship(t *testing.T) {
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

	type args struct {
		coordinatorID string
		producerIDs   []string
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name: "success to relate",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				p := testProducer("admin-id", "", "&.農園", now())
				err = db.DB.Create(&p).Error
				require.NoError(t, err)
			},
			args: args{
				coordinatorID: "coordinator-id",
				producerIDs:   []string{"admin-id"},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "success to unrelate",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				p := testProducer("admin-id", "coordinator-id", "&.農園", now())
				err = db.DB.Create(&p).Error
				require.NoError(t, err)
			},
			args: args{
				coordinatorID: "",
				producerIDs:   []string{"admin-id"},
			},
			want: want{
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, producerTable, coordinatorTable, adminTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &producer{db: db, now: now}
			err = db.UpdateRelationship(ctx, tt.args.coordinatorID, tt.args.producerIDs...)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProducer_Delete(t *testing.T) {
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

	type args struct {
		producerID string
		auth       func(ctx context.Context) error
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				p := testProducer("admin-id", "coordinator-id", "&.農園", now())
				err = db.DB.Create(&p).Error
				require.NoError(t, err)
			},
			args: args{
				producerID: "admin-id",
				auth:       func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				producerID: "admin-id",
				auth:       func(ctx context.Context) error { return nil },
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to execute external service",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				p := testProducer("admin-id", "coordinator-id", "&.農園", now())
				err = db.DB.Create(&p).Error
				require.NoError(t, err)
			},
			args: args{
				producerID: "admin-id",
				auth:       func(ctx context.Context) error { return assert.AnError },
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

			err := delete(ctx, producerTable, coordinatorTable, adminTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &producer{db: db, now: now}
			err = db.Delete(ctx, tt.args.producerID, tt.args.auth)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testProducer(id, coordinatorID, storeName string, now time.Time) *entity.Producer {
	p := &entity.Producer{
		AdminID:       id,
		CoordinatorID: coordinatorID,
		StoreName:     storeName,
		ThumbnailURL:  "https://and-period.jp/thumbnail.png",
		Thumbnails:    common.Images{},
		HeaderURL:     "https://and-period.jp/header.png",
		Headers:       common.Images{},
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    "東京都",
		City:          "千代田区",
		AddressLine1:  "永田町1-7-1",
		AddressLine2:  "",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	fillProducerJSON(p)
	return p
}

func fillProducerJSON(p *entity.Producer) {
	thumbnails, _ := json.Marshal(p.Thumbnails)
	headers, _ := json.Marshal(p.Headers)
	p.ThumbnailsJSON = datatypes.JSON(thumbnails)
	p.HeadersJSON = datatypes.JSON(headers)
}
