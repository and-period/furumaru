package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListShippingsByShopID(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 7, 15, 18, 30, 0, 0)
	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for val := range codes.PrefectureNames {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := entity.ShippingRates{
		{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
	}
	shippings := entity.Shippings{
		{
			ID:            "shipping-id",
			CoordinatorID: "coordinator-id",
			CreatedAt:     now,
			UpdatedAt:     now,
			ShippingRevision: entity.ShippingRevision{
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
		},
	}
	params := &database.ListShippingsParams{
		ShopID: "shop-id",
		Limit:  200,
		Offset: 0,
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListShippingsByShopIDInput
		expect      entity.Shippings
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().List(gomock.Any(), params).Return(shippings, nil)
				mocks.db.Shipping.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListShippingsByShopIDInput{
				ShopID: "shop-id",
				Limit:  200,
				Offset: 0,
			},
			expect:      shippings,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.ListShippingsByShopIDInput{
				ShopID: "",
				Limit:  200,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list shippings",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Shipping.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListShippingsByShopIDInput{
				ShopID: "shop-id",
				Limit:  200,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count shippings",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().List(gomock.Any(), params).Return(shippings, nil)
				mocks.db.Shipping.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListShippingsByShopIDInput{
				ShopID: "shop-id",
				Limit:  200,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListShippingsByShopID(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestListShippingsByCoordinatorIDs(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 7, 15, 18, 30, 0, 0)
	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for val := range codes.PrefectureNames {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := entity.ShippingRates{
		{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
	}
	shippings := entity.Shippings{
		{
			ID:            "shipping-id",
			CoordinatorID: "coordinator-id",
			CreatedAt:     now,
			UpdatedAt:     now,
			ShippingRevision: entity.ShippingRevision{
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.ListShippingsByCoordinatorIDsInput
		expect    entity.Shippings
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().ListByCoordinatorIDs(gomock.Any(), []string{"coordinator-id"}).Return(shippings, nil)
			},
			input: &store.ListShippingsByCoordinatorIDsInput{
				CoordinatorIDs: []string{"coordinator-id"},
			},
			expect:    shippings,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.ListShippingsByCoordinatorIDsInput{
				CoordinatorIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list shippings",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().ListByCoordinatorIDs(gomock.Any(), []string{"coordinator-id"}).Return(nil, assert.AnError)
			},
			input: &store.ListShippingsByCoordinatorIDsInput{
				CoordinatorIDs: []string{"coordinator-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.ListShippingsByCoordinatorIDs(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestMutiGetShippingsByRevision(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 7, 16, 18, 30, 0, 0)
	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for val := range codes.PrefectureNames {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := entity.ShippingRates{
		{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
	}
	shippings := entity.Shippings{
		{
			ID:            "shipping-id",
			CoordinatorID: "coordinator-id",
			CreatedAt:     now,
			UpdatedAt:     now,
			ShippingRevision: entity.ShippingRevision{
				ID:                1,
				ShippingID:        "shipping-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetShippingsByRevisionInput
		expect    entity.Shippings
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().MultiGetByRevision(ctx, []int64{1}).Return(shippings, nil)
			},
			input: &store.MultiGetShippingsByRevisionInput{
				ShippingRevisionIDs: []int64{1},
			},
			expect:    shippings,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetShippingsByRevisionInput{
				ShippingRevisionIDs: []int64{0},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multiGet",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().MultiGetByRevision(ctx, []int64{1}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetShippingsByRevisionInput{
				ShippingRevisionIDs: []int64{1},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetShippingsByRevision(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetShipping(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 7, 16, 18, 30, 0, 0)
	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for val := range codes.PrefectureNames {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := entity.ShippingRates{
		{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
	}
	shipping := &entity.Shipping{
		ID:            "shipping-id",
		CoordinatorID: "",
		CreatedAt:     now,
		UpdatedAt:     now,
		ShippingRevision: entity.ShippingRevision{
			ID:                1,
			ShippingID:        "shipping-id",
			Box60Rates:        rates,
			Box60Frozen:       800,
			Box80Rates:        rates,
			Box80Frozen:       800,
			Box100Rates:       rates,
			Box100Frozen:      800,
			HasFreeShipping:   true,
			FreeShippingRates: 3000,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetShippingInput
		expect    *entity.Shipping
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Get(ctx, "shipping-id").Return(shipping, nil)
			},
			input: &store.GetShippingInput{
				ShippingID: "shipping-id",
			},
			expect:    shipping,
			expectErr: nil,
		},
		{
			name: "failed to get shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Get(ctx, "shipping-id").Return(nil, assert.AnError)
			},
			input: &store.GetShippingInput{
				ShippingID: "shipping-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetShipping(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestGetDefaultShipping(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 7, 16, 18, 30, 0, 0)
	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for val := range codes.PrefectureNames {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := entity.ShippingRates{
		{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
	}
	shipping := &entity.Shipping{
		ID:            entity.DefaultShippingID,
		CoordinatorID: "",
		CreatedAt:     now,
		UpdatedAt:     now,
		ShippingRevision: entity.ShippingRevision{
			ID:                1,
			ShippingID:        entity.DefaultShippingID,
			Box60Rates:        rates,
			Box60Frozen:       800,
			Box80Rates:        rates,
			Box80Frozen:       800,
			Box100Rates:       rates,
			Box100Frozen:      800,
			HasFreeShipping:   true,
			FreeShippingRates: 3000,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetDefaultShippingInput
		expect    *entity.Shipping
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().GetDefault(ctx).Return(shipping, nil)
			},
			input:     &store.GetDefaultShippingInput{},
			expect:    shipping,
			expectErr: nil,
		},
		{
			name: "failed to get shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().GetDefault(ctx).Return(nil, assert.AnError)
			},
			input:     &store.GetDefaultShippingInput{},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetDefaultShipping(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestGetShippingByCoordinatorID(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 7, 16, 18, 30, 0, 0)
	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for val := range codes.PrefectureNames {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := entity.ShippingRates{
		{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
	}
	shipping := &entity.Shipping{
		ID:            "shipping-id",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
		ShippingRevision: entity.ShippingRevision{
			ID:                1,
			ShippingID:        "shipping-id",
			Box60Rates:        rates,
			Box60Frozen:       800,
			Box80Rates:        rates,
			Box80Frozen:       800,
			Box100Rates:       rates,
			Box100Frozen:      800,
			HasFreeShipping:   true,
			FreeShippingRates: 3000,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetShippingByCoordinatorIDInput
		expect    *entity.Shipping
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(shipping, nil)
			},
			input: &store.GetShippingByCoordinatorIDInput{
				CoordinatorID: "coordinator-id",
			},
			expect:    shipping,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetShippingByCoordinatorIDInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(nil, assert.AnError)
			},
			input: &store.GetShippingByCoordinatorIDInput{
				CoordinatorID: "coordinator-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetShippingByCoordinatorID(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateShipping(t *testing.T) {
	t.Parallel()

	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for val := range codes.PrefectureNames {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := []*store.CreateShippingRate{
		{Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Name: "その他", Price: 500, PrefectureCodes: others},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateShippingInput
		expectErr error
	}{
		{
			name: "success first create",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().
					GetByCoordinatorID(ctx, "coordinator-id").
					Return(nil, database.ErrNotFound)
				mocks.db.Shipping.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, shipping *entity.Shipping) error {
						expect := &entity.Shipping{
							ID:            shipping.ID, // ignore
							ShopID:        "shop-id",
							CoordinatorID: "coordinator-id",
							ShippingRevision: entity.ShippingRevision{
								ShippingID: shipping.ID,
								Box60Rates: entity.ShippingRates{
									{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
									{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
								},
								Box60Frozen: 800,
								Box80Rates: entity.ShippingRates{
									{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
									{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
								},
								Box80Frozen: 800,
								Box100Rates: entity.ShippingRates{
									{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
									{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
								},
								Box100Frozen:      800,
								HasFreeShipping:   true,
								FreeShippingRates: 3000,
							},
							InUse: true,
						}
						assert.Equal(t, expect, shipping)
						return nil
					})
			},
			input: &store.CreateShippingInput{
				ShopID:            "shop-id",
				CoordinatorID:     "coordinator-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expectErr: nil,
		},
		{
			name: "success second create",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().
					GetByCoordinatorID(ctx, "coordinator-id").
					Return(&entity.Shipping{}, nil)
				mocks.db.Shipping.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, shipping *entity.Shipping) error {
						assert.False(t, shipping.InUse)
						return nil
					})
			},
			input: &store.CreateShippingInput{
				ShopID:            "shop-id",
				CoordinatorID:     "coordinator-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateShippingInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid box 60 rates",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: &store.CreateShippingInput{
				ShopID:            "shop-id",
				CoordinatorID:     "coordinator-id",
				Box60Rates:        []*store.CreateShippingRate{},
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid box 80 rates",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.CreateShippingInput{
				ShopID:            "shop-id",
				CoordinatorID:     "coordinator-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        []*store.CreateShippingRate{},
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid box 100 rates",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.CreateShippingInput{
				ShopID:        "shop-id",
				CoordinatorID: "coordinator-id",
				Box60Rates:    rates,
				Box60Frozen:   800,
				Box80Rates:    rates,
				Box80Frozen:   800,
				Box100Rates:   []*store.CreateShippingRate{},
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get by coordinator id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(nil, assert.AnError)
			},
			input: &store.CreateShippingInput{
				ShopID:            "shop-id",
				CoordinatorID:     "coordinator-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(nil, database.ErrNotFound)
				mocks.db.Shipping.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateShippingInput{
				ShopID:            "shop-id",
				CoordinatorID:     "coordinator-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateShipping(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateShipping(t *testing.T) {
	t.Parallel()

	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for val := range codes.PrefectureNames {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := []*store.UpdateShippingRate{
		{Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Name: "その他", Price: 500, PrefectureCodes: others},
	}
	params := &database.UpdateShippingParams{
		Box60Rates: entity.ShippingRates{
			{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
			{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
		},
		Box60Frozen: 800,
		Box80Rates: entity.ShippingRates{
			{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
			{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
		},
		Box80Frozen: 800,
		Box100Rates: entity.ShippingRates{
			{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
			{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
		},
		Box100Frozen:      800,
		HasFreeShipping:   true,
		FreeShippingRates: 3000,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateShippingInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Update(ctx, "shipping-id", params).Return(nil)
			},
			input: &store.UpdateShippingInput{
				ShippingID:        "shipping-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateShippingInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid box 60 rates",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: &store.UpdateShippingInput{
				ShippingID:        "shipping-id",
				Box60Rates:        []*store.UpdateShippingRate{},
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid box 80 rates",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: &store.UpdateShippingInput{
				ShippingID:   "shipping-id",
				Box60Rates:   rates,
				Box60Frozen:  800,
				Box80Rates:   []*store.UpdateShippingRate{},
				Box80Frozen:  800,
				Box100Rates:  rates,
				Box100Frozen: 800,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid box 100 rates",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: &store.UpdateShippingInput{
				ShippingID:   "shipping-id",
				Box60Rates:   rates,
				Box60Frozen:  800,
				Box80Rates:   rates,
				Box80Frozen:  800,
				Box100Rates:  []*store.UpdateShippingRate{},
				Box100Frozen: 800,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Update(ctx, "shipping-id", params).Return(assert.AnError)
			},
			input: &store.UpdateShippingInput{
				ShippingID:        "shipping-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateShipping(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateShippingInUse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateShippingInUseInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().UpdateInUse(ctx, "shop-id", "shipping-id").Return(nil)
			},
			input: &store.UpdateShippingInUseInput{
				ShopID:     "shop-id",
				ShippingID: "shipping-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateShippingInUseInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update in use",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().UpdateInUse(ctx, "shop-id", "shipping-id").Return(assert.AnError)
			},
			input: &store.UpdateShippingInUseInput{
				ShopID:     "shop-id",
				ShippingID: "shipping-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateShippingInUse(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateDefaultShipping(t *testing.T) {
	t.Parallel()

	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for val := range codes.PrefectureNames {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := []*store.UpdateDefaultShippingRate{
		{Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Name: "その他", Price: 500, PrefectureCodes: others},
	}
	params := &database.UpdateShippingParams{
		Box60Rates: entity.ShippingRates{
			{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
			{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
		},
		Box60Frozen: 800,
		Box80Rates: entity.ShippingRates{
			{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
			{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
		},
		Box80Frozen: 800,
		Box100Rates: entity.ShippingRates{
			{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
			{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
		},
		Box100Frozen:      800,
		HasFreeShipping:   true,
		FreeShippingRates: 3000,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateDefaultShippingInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Update(ctx, entity.DefaultShippingID, params).Return(nil)
			},
			input: &store.UpdateDefaultShippingInput{
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateDefaultShippingInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid box 60 rates",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: &store.UpdateDefaultShippingInput{
				Box60Rates:        []*store.UpdateDefaultShippingRate{},
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid box 80 rates",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: &store.UpdateDefaultShippingInput{
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        []*store.UpdateDefaultShippingRate{},
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "invalid box 100 rates",
			setup: func(ctx context.Context, mocks *mocks) {
			},
			input: &store.UpdateDefaultShippingInput{
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       []*store.UpdateDefaultShippingRate{},
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Update(ctx, entity.DefaultShippingID, params).Return(assert.AnError)
			},
			input: &store.UpdateDefaultShippingInput{
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateDefaultShipping(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpsertShipping(t *testing.T) {
	t.Parallel()

	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for val := range codes.PrefectureNames {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := []*store.UpsertShippingRate{
		{Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Name: "その他", Price: 500, PrefectureCodes: others},
	}
	params := &database.UpdateShippingParams{
		Box60Rates: entity.ShippingRates{
			{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
			{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
		},
		Box60Frozen: 800,
		Box80Rates: entity.ShippingRates{
			{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
			{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
		},
		Box80Frozen: 800,
		Box100Rates: entity.ShippingRates{
			{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
			{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
		},
		Box100Frozen:      800,
		HasFreeShipping:   true,
		FreeShippingRates: 3000,
	}
	shipping := func(revisionID int64) *entity.Shipping {
		rates := entity.ShippingRates{
			{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
			{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
		}
		return &entity.Shipping{
			ID:            "coordinator-id",
			ShopID:        "shop-id",
			CoordinatorID: "coordinator-id",
			InUse:         true,
			ShippingRevision: entity.ShippingRevision{
				ShippingID:        "coordinator-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
		}
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpsertShippingInput
		expectErr error
	}{
		{
			name: "success to create",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(nil, database.ErrNotFound)
				mocks.db.Shipping.EXPECT().Create(ctx, shipping(0)).Return(nil)
			},
			input: &store.UpsertShippingInput{
				ShopID:            "shop-id",
				CoordinatorID:     "coordinator-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
				InUse:             true,
			},
			expectErr: nil,
		},
		{
			name: "success to update",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(shipping(1), nil)
				mocks.db.Shipping.EXPECT().Update(ctx, "coordinator-id", params).Return(nil)
			},
			input: &store.UpsertShippingInput{
				ShopID:            "shop-id",
				CoordinatorID:     "coordinator-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
				InUse:             true,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpsertShippingInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(nil, assert.AnError)
			},
			input: &store.UpsertShippingInput{
				ShopID:            "shop-id",
				CoordinatorID:     "coordinator-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
				InUse:             true,
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(nil, database.ErrNotFound)
				mocks.db.Shipping.EXPECT().Create(ctx, shipping(0)).Return(assert.AnError)
			},
			input: &store.UpsertShippingInput{
				ShopID:            "shop-id",
				CoordinatorID:     "coordinator-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
				InUse:             true,
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to update",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(ctx, "coordinator-id").Return(shipping(1), nil)
				mocks.db.Shipping.EXPECT().Update(ctx, "coordinator-id", params).Return(assert.AnError)
			},
			input: &store.UpsertShippingInput{
				ShopID:            "shop-id",
				CoordinatorID:     "coordinator-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
				InUse:             true,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpsertShipping(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteShipping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteShippingInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Delete(ctx, "shipping-id").Return(nil)
			},
			input:     &store.DeleteShippingInput{ShippingID: "shipping-id"},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteShippingInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Delete(ctx, "shipping-id").Return(assert.AnError)
			},
			input:     &store.DeleteShippingInput{ShippingID: "shipping-id"},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteShipping(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
