package mysql

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
	"gorm.io/gorm"
)

func TestProducer(t *testing.T) {
	assert.NotNil(t, newProducer(nil))
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
	producers[0] = testProducer("admin-id01", "coordinator-id", now())
	producers[0].Admin = *admins[0]
	producers[1] = testProducer("admin-id02", "coordinator-id", now())
	producers[1].Admin = *admins[1]
	err = db.DB.Create(&producers).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListProducersParams
	}
	type want struct {
		producers entity.Producers
		hasErr    bool
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
				params: &database.ListProducersParams{
					CoordinatorID: "coordinator-id",
					Name:          "農園",
					Limit:         1,
					Offset:        1,
				},
			},
			want: want{
				producers: producers[1:],
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
	producers[0] = testProducer("admin-id01", "coordinator-id", now())
	producers[0].Admin = *admins[0]
	producers[1] = testProducer("admin-id02", "coordinator-id", now())
	producers[1].Admin = *admins[1]
	err = db.DB.Create(&producers).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListProducersParams
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
				params: &database.ListProducersParams{
					Name: "農園",
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
	producers[0] = testProducer("admin-id01", "coordinator-id", now())
	producers[0].Admin = *admins[0]
	producers[1] = testProducer("admin-id02", "coordinator-id", now())
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
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
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

func TestProducer_MultiGetWithDeleted(t *testing.T) {
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
	admins[0].Status = entity.AdminStatusDeactivated
	admins[0].DeletedAt = gorm.DeletedAt{Valid: true, Time: now()}
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	producers := make(entity.Producers, 2)
	producers[0] = testProducer("admin-id01", "coordinator-id", now())
	producers[0].Admin = *admins[0]
	producers[1] = testProducer("admin-id02", "coordinator-id", now())
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
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
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
			actual, err := db.MultiGetWithDeleted(ctx, tt.args.producerIDs)
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
	p := testProducer("admin-id", "coordinator-id", now())
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
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
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
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
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

func TestProducer_GetWithDeleted(t *testing.T) {
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
	admins[0].Status = entity.AdminStatusDeactivated
	admins[0].DeletedAt = gorm.DeletedAt{Valid: true, Time: now()}
	admins[1] = testAdmin("admin-id02", "cognito-id02", "test-admin02@and-period.jp", now())
	err = db.DB.Create(&admins).Error
	producers := make(entity.Producers, 2)
	producers[0] = testProducer("admin-id01", "coordinator-id", now())
	producers[0].Admin = *admins[0]
	producers[1] = testProducer("admin-id02", "coordinator-id", now())
	producers[1].Admin = *admins[1]
	err = db.DB.Create(&producers).Error
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
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success to activated",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				producerID: "admin-id01",
			},
			want: want{
				producer: producers[0],
				hasErr:   false,
			},
		},
		{
			name:  "success to deactivated",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				producerID: "admin-id02",
			},
			want: want{
				producer: producers[1],
				hasErr:   false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
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
			actual, err := db.GetWithDeleted(ctx, tt.args.producerID)
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

	p := testProducer("admin-id", "coordinator-id", now())
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
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
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
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
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
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				p := testProducer("admin-id", "coordinator-id", now())
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
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
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
		params     *database.UpdateProducerParams
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
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				p := testProducer("admin-id", "coordinator-id", now())
				err = db.DB.Create(&p).Error
				require.NoError(t, err)
			},
			args: args{
				producerID: "admin-id",
				params: &database.UpdateProducerParams{
					Lastname:       "&.",
					Firstname:      "スタッフ",
					LastnameKana:   "あんどぴりおど",
					FirstnameKana:  "すたっふ",
					ThumbnailURL:   "https://and-period.jp/thumbnail.png",
					HeaderURL:      "https://and-period.jp/header.png",
					Email:          "test-admin01@and-period.jp",
					PhoneNumber:    "+819012345678",
					PostalCode:     "1000014",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
				},
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
			err = db.Update(ctx, tt.args.producerID, tt.args.params)
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
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				p := testProducer("admin-id", "coordinator-id", now())
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
			name: "failed to execute external service",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				coordinator := testCoordinator("coordinator-id", now())
				coordinator.Admin = *testAdmin("coordinator-id", "coordinator-id", "test-coordinator@and-period.jp", now())
				err = db.DB.Create(&coordinator.Admin).Error
				require.NoError(t, err)
				err = db.DB.Create(&coordinator).Error
				require.NoError(t, err)
				admin := testAdmin("admin-id", "cognito-id", "test-admin01@and-period.jp", now())
				err = db.DB.Create(&admin).Error
				require.NoError(t, err)
				p := testProducer("admin-id", "coordinator-id", now())
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

func TestProducer_AggregateByCoordinatorID(t *testing.T) {
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
	producers[0] = testProducer("admin-id01", "coordinator-id", now())
	producers[0].Admin = *admins[0]
	producers[1] = testProducer("admin-id02", "coordinator-id", now())
	producers[1].Admin = *admins[1]
	err = db.DB.Create(&producers).Error
	require.NoError(t, err)

	type args struct {
		coordinatorIDs []string
	}
	type want struct {
		total  map[string]int64
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
				coordinatorIDs: []string{"coordinator-id"},
			},
			want: want{
				total: map[string]int64{
					"coordinator-id": 2,
				},
				hasErr: false,
			},
		},
		{
			name:  "empty",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				coordinatorIDs: []string{},
			},
			want: want{
				total:  map[string]int64{},
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
			actual, err := db.AggregateByCoordinatorID(ctx, tt.args.coordinatorIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func testProducer(id, coordinatorID string, now time.Time) *entity.Producer {
	return &entity.Producer{
		AdminID:           id,
		CoordinatorID:     coordinatorID,
		Username:          "&.農園",
		ThumbnailURL:      "https://and-period.jp/thumbnail.png",
		HeaderURL:         "https://and-period.jp/header.png",
		PromotionVideoURL: "https://and-period.jp/promotion.mp4",
		BonusVideoURL:     "https://and-period.jp/bonus.mp4",
		InstagramID:       "instagram-id",
		FacebookID:        "facebook-id",
		PhoneNumber:       "+819012345678",
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
