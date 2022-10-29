package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/codes"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	set "github.com/and-period/furumaru/api/pkg/set/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShipping(t *testing.T) {
	assert.NotNil(t, NewShipping(nil))
}

func TestShipping_List(t *testing.T) {
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

	_ = m.dbDelete(ctx, shippingTable)
	shippings := make(entity.Shippings, 3)
	shippings[0] = testShipping("shipping-id01", now())
	shippings[1] = testShipping("shipping-id02", now())
	shippings[2] = testShipping("shipping-id03", now())
	err = m.db.DB.Create(&shippings).Error
	require.NoError(t, err)

	type args struct {
		params *ListShippingsParams
	}
	type want struct {
		shippings entity.Shippings
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
				params: &ListShippingsParams{
					Limit:  20,
					Offset: 1,
				},
			},
			want: want{
				shippings: shippings[1:],
				hasErr:    false,
			},
		},
		{
			name:  "success with sort",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				params: &ListShippingsParams{
					Orders: []*ListShippingsOrder{
						{Key: entity.ShippingOrderByCreatedAt, OrderByASC: true},
						{Key: entity.ShippingOrderByUpdatedAt, OrderByASC: false},
					},
				},
			},
			want: want{
				shippings: shippings,
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

			db := &shipping{db: m.db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			fillIgnoreShippingsField(actual, now())
			assert.Equal(t, tt.want.shippings, actual)
		})
	}
}

func TestShipping_Count(t *testing.T) {
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

	_ = m.dbDelete(ctx, shippingTable)
	shippings := make(entity.Shippings, 3)
	shippings[0] = testShipping("shipping-id01", now())
	shippings[1] = testShipping("shipping-id02", now())
	shippings[2] = testShipping("shipping-id03", now())
	err = m.db.DB.Create(&shippings).Error
	require.NoError(t, err)

	type args struct {
		params *ListShippingsParams
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
				params: &ListShippingsParams{
					Limit:  20,
					Offset: 1,
				},
			},
			want: want{
				total:  3,
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

			db := &shipping{db: m.db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestShipping_MultiGet(t *testing.T) {
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

	_ = m.dbDelete(ctx, shippingTable)
	shippings := make(entity.Shippings, 2)
	shippings[0] = testShipping("shipping-id01", now())
	shippings[1] = testShipping("shipping-id02", now())
	err = m.db.DB.Create(&shippings).Error
	require.NoError(t, err)

	type args struct {
		shippingIDs []string
	}
	type want struct {
		shippings entity.Shippings
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
				shippingIDs: []string{"shipping-id01", "shipping-id02", "shipping-id03"},
			},
			want: want{
				shippings: shippings[:2],
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

			db := &shipping{db: m.db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.shippingIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			fillIgnoreShippingsField(actual, now())
			assert.Equal(t, tt.want.shippings, actual)
		})
	}
}

func TestShipping_Get(t *testing.T) {
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

	_ = m.dbDelete(ctx, shippingTable)
	s := testShipping("shipping-id", now())
	err = m.db.DB.Create(&s).Error
	require.NoError(t, err)

	type args struct {
		shippingID string
	}
	type want struct {
		shipping *entity.Shipping
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
				shippingID: "shipping-id",
			},
			want: want{
				shipping: s,
				hasErr:   false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				shippingID: "other-id",
			},
			want: want{
				shipping: nil,
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

			db := &shipping{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.shippingID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			fillIgnoreShippingField(actual, now())
			assert.Equal(t, tt.want.shipping, actual)
		})
	}
}

func TestShipping_Create(t *testing.T) {
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

	_ = m.dbDelete(ctx, shippingTable)
	s := testShipping("shipping-id", now())

	type args struct {
		shipping *entity.Shipping
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
				shipping: s,
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate key",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				err := m.db.DB.Create(&s).Error
				require.NoError(t, err)
			},
			args: args{
				shipping: s,
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

			err := m.dbDelete(ctx, shippingTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &shipping{db: m.db, now: now}
			err = db.Create(ctx, tt.args.shipping)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestShipping_Update(t *testing.T) {
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

	_ = m.dbDelete(ctx, shippingTable)
	s := testShipping("shipping-id", now())

	type args struct {
		shippingID string
		params     *UpdateShippingParams
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
				err := m.db.DB.Create(&s).Error
				require.NoError(t, err)
			},
			args: args{
				shippingID: "shipping-id",
				params: &UpdateShippingParams{
					Name:               "デフォルト配送設定",
					Box60Rates:         s.Box60Rates,
					Box60Refrigerated:  500,
					Box60Frozen:        800,
					Box80Rates:         s.Box80Rates,
					Box80Refrigerated:  500,
					Box80Frozen:        800,
					Box100Rates:        s.Box100Rates,
					Box100Refrigerated: 500,
					Box100Frozen:       800,
					HasFreeShipping:    true,
					FreeShippingRates:  3000,
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
				shippingID: "shipping-id",
				params:     &UpdateShippingParams{},
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

			err := m.dbDelete(ctx, shippingTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &shipping{db: m.db, now: now}
			err = db.Update(ctx, tt.args.shippingID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestShipping_Delete(t *testing.T) {
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

	_ = m.dbDelete(ctx, shippingTable)
	s := testShipping("shipping-id", now())

	type args struct {
		shippingID string
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
				err := m.db.DB.Create(&s).Error
				require.NoError(t, err)
			},
			args: args{
				shippingID: "shipping-id",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				shippingID: "shipping-id",
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

			err := m.dbDelete(ctx, shippingTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &shipping{db: m.db, now: now}
			err = db.Delete(ctx, tt.args.shippingID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testShipping(id string, now time.Time) *entity.Shipping {
	shikoku := []int64{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int64, 0, 47-len(shikoku))
	for _, val := range codes.PrefectureValues {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := entity.ShippingRates{
		{Number: 1, Name: "四国", Price: 250, Prefectures: shikoku},
		{Number: 2, Name: "その他", Price: 500, Prefectures: others},
	}
	shipping := &entity.Shipping{
		ID:                 id,
		Name:               "デフォルト配送設定",
		Box60Rates:         rates,
		Box60Refrigerated:  500,
		Box60Frozen:        800,
		Box80Rates:         rates,
		Box80Refrigerated:  500,
		Box80Frozen:        800,
		Box100Rates:        rates,
		Box100Refrigerated: 500,
		Box100Frozen:       800,
		HasFreeShipping:    true,
		FreeShippingRates:  3000,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	_ = shipping.FillJSON()
	return shipping
}

func fillIgnoreShippingField(s *entity.Shipping, now time.Time) {
	if s == nil {
		return
	}
	_ = s.FillJSON()
	s.CreatedAt = now
	s.UpdatedAt = now
}

func fillIgnoreShippingsField(ss entity.Shippings, now time.Time) {
	for i := range ss {
		fillIgnoreShippingField(ss[i], now)
	}
}
