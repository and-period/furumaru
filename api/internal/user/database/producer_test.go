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

func TestProducer(t *testing.T) {
	assert.NotNil(t, NewProducer(nil))
}

func TestProducer_List(t *testing.T) {
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

	_ = m.dbDelete(ctx, producerTable)
	producers := make(entity.Producers, 2)
	producers[0] = testProducer("admin-id01", "&.農園", "test-admin01@and-period.jp", now())
	producers[1] = testProducer("admin-id02", "&.水産", "test-admin02@and-period.jp", now())
	err = m.db.DB.Create(&producers).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				params: &ListProducersParams{
					Limit:  1,
					Offset: 1,
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

			tt.setup(ctx, t, m)

			db := &producer{db: m.db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreProducersField(actual, now())
			assert.Equal(t, tt.want.producers, actual)
		})
	}
}

func TestProducer_Get(t *testing.T) {
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

	_ = m.dbDelete(ctx, producerTable)
	p := testProducer("admin-id", "&.農園", "test-admin@and-period.jp", now())
	err = m.db.DB.Create(&p).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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

			tt.setup(ctx, t, m)

			db := &producer{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.producerID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreProducerField(actual, now())
			assert.Equal(t, tt.want.producer, actual)
		})
	}
}

func TestProducer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	type args struct {
		auth     *entity.AdminAuth
		producer *entity.Producer
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
				auth:     testAdminAuth("admin-id", "cognito-id", now()),
				producer: testProducer("admin-id", "&.農園", "test-admin@and-period.jp", now()),
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
				auth:     testAdminAuth("admin-id", "cognito-id", now()),
				producer: testProducer("admin-id", "&.農園", "test-admin@and-period.jp", now()),
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to duplicate entry in producer",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				producer := testProducer("admin-id", "&.農園", "test-admin@and-period.jp", now())
				err = m.db.DB.Create(&producer).Error
				require.NoError(t, err)
			},
			args: args{
				auth:     testAdminAuth("admin-id", "cognito-id", now()),
				producer: testProducer("admin-id", "&.農園", "test-admin@and-period.jp", now()),
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

			err := m.dbDelete(ctx, adminAuthTable, producerTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &producer{db: m.db, now: now}
			err = db.Create(ctx, tt.args.auth, tt.args.producer)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestProducer_UpdateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	type args struct {
		producerID string
		email      string
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
				p := testProducer("admin-id", "&.農園", "test-admin@and-period.jp", now())
				err = m.db.DB.Create(&p).Error
				require.NoError(t, err)
			},
			args: args{
				producerID: "admin-id",
				email:      "test-other@and-period.jp",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				producerID: "admin-id",
				email:      "test-other@and-period.jp",
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

			err := m.dbDelete(ctx, producerTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &producer{db: m.db, now: now}
			err = db.UpdateEmail(ctx, tt.args.producerID, tt.args.email)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testProducer(id, storeName, email string, now time.Time) *entity.Producer {
	return &entity.Producer{
		ID:            id,
		Lastname:      "&.",
		Firstname:     "スタッフ",
		LastnameKana:  "あんどぴりおど",
		FirstnameKana: "すたっふ",
		StoreName:     storeName,
		ThumbnailURL:  "https://and-period.jp/thumbnail.png",
		HeaderURL:     "https://and-period.jp/header.png",
		Email:         email,
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    "東京都",
		City:          "千代田区",
		AddressLine1:  "永田町1-7-1",
		AddressLine2:  "",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func fillIgnoreProducerField(p *entity.Producer, now time.Time) {
	if p == nil {
		return
	}
	p.CreatedAt = now
	p.UpdatedAt = now
}

func fillIgnoreProducersField(ps entity.Producers, now time.Time) {
	for i := range ps {
		fillIgnoreProducerField(ps[i], now)
	}
}
