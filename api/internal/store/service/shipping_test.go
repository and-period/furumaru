package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListShippings(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 7, 15, 18, 30, 0, 0)
	params := &database.ListShippingsParams{
		CoordinatorID: "coordinator-id",
		Limit:         20,
		Offset:        0,
		Orders: []*database.ListShippingsOrder{
			{Key: entity.ShippingOrderByName, OrderByASC: true},
		},
	}
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
	shippings := entity.Shippings{
		{
			ID:                 "shipping-id",
			CoordinatorID:      "coordinator-id",
			Name:               "デフォルト配送設定",
			IsDefault:          true,
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
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListShippingsInput
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
			input: &store.ListShippingsInput{
				CoordinatorID: "coordinator-id",
				Limit:         20,
				Offset:        0,
				Orders: []*store.ListShippingsOrder{
					{Key: entity.ShippingOrderByName, OrderByASC: true},
				},
			},
			expect:      shippings,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &store.ListShippingsInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   store.ErrInvalidArgument,
		},
		{
			name: "failed to list shippings",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Shipping.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListShippingsInput{
				CoordinatorID: "coordinator-id",
				Limit:         20,
				Offset:        0,
				Orders: []*store.ListShippingsOrder{
					{Key: entity.ShippingOrderByName, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   store.ErrInternal,
		},
		{
			name: "failed to count shippings",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().List(gomock.Any(), params).Return(shippings, nil)
				mocks.db.Shipping.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListShippingsInput{
				CoordinatorID: "coordinator-id",
				Limit:         20,
				Offset:        0,
				Orders: []*store.ListShippingsOrder{
					{Key: entity.ShippingOrderByName, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   store.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListShippings(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestMutiGetShippings(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 7, 16, 18, 30, 0, 0)
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
	shippings := entity.Shippings{
		{
			ID:                 "shipping-id",
			CoordinatorID:      "coordinator-id",
			Name:               "デフォルト配送設定",
			IsDefault:          true,
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
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetShippingsInput
		expect    entity.Shippings
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().MultiGet(ctx, []string{"shipping-id"}).Return(shippings, nil)
			},
			input: &store.MultiGetShippingsInput{
				ShippingIDs: []string{"shipping-id"},
			},
			expect:    shippings,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetShippingsInput{
				ShippingIDs: []string{""},
			},
			expect:    nil,
			expectErr: store.ErrInvalidArgument,
		},
		{
			name: "failed to multiGet",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().MultiGet(ctx, []string{"shipping-id"}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetShippingsInput{
				ShippingIDs: []string{"shipping-id"},
			},
			expect:    nil,
			expectErr: store.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetShippings(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetShipping(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 7, 16, 18, 30, 0, 0)
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
		ID:                 "shipping-id",
		CoordinatorID:      "coordinator-id",
		Name:               "デフォルト配送設定",
		IsDefault:          true,
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
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetShippingInput{},
			expect:    nil,
			expectErr: store.ErrInvalidArgument,
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
			expectErr: store.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetShipping(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateShipping(t *testing.T) {
	t.Parallel()

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
	rates := []*store.CreateShippingRate{
		{Name: "四国", Price: 250, Prefectures: shikoku},
		{Name: "その他", Price: 500, Prefectures: others},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateShippingInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, shipping *entity.Shipping) error {
						rates := entity.ShippingRates{
							{Number: 1, Name: "四国", Price: 250, Prefectures: shikoku},
							{Number: 2, Name: "その他", Price: 500, Prefectures: others},
						}
						expect := &entity.Shipping{
							ID:                 shipping.ID, // ignore
							CoordinatorID:      "coordinator-id",
							Name:               "デフォルト配送設定",
							IsDefault:          true,
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
						}
						assert.Equal(t, expect, shipping)
						return nil
					})
			},
			input: &store.CreateShippingInput{
				CoordinatorID:      "coordinator-id",
				Name:               "デフォルト配送設定",
				IsDefault:          true,
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
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateShippingInput{},
			expectErr: store.ErrInvalidArgument,
		},
		{
			name: "failed to create shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateShippingInput{
				CoordinatorID:      "coordinator-id",
				Name:               "デフォルト配送設定",
				IsDefault:          true,
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
			},
			expectErr: store.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateShipping(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateShipping(t *testing.T) {
	t.Parallel()

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
	rates := []*store.UpdateShippingRate{
		{Name: "四国", Price: 250, Prefectures: shikoku},
		{Name: "その他", Price: 500, Prefectures: others},
	}
	params := &database.UpdateShippingParams{
		Name:      "デフォルト配送設定",
		IsDefault: true,
		Box60Rates: entity.ShippingRates{
			{Number: 1, Name: "四国", Price: 250, Prefectures: shikoku},
			{Number: 2, Name: "その他", Price: 500, Prefectures: others},
		},
		Box60Refrigerated: 500,
		Box60Frozen:       800,
		Box80Rates: entity.ShippingRates{
			{Number: 1, Name: "四国", Price: 250, Prefectures: shikoku},
			{Number: 2, Name: "その他", Price: 500, Prefectures: others},
		},
		Box80Refrigerated: 500,
		Box80Frozen:       800,
		Box100Rates: entity.ShippingRates{
			{Number: 1, Name: "四国", Price: 250, Prefectures: shikoku},
			{Number: 2, Name: "その他", Price: 500, Prefectures: others},
		},
		Box100Refrigerated: 500,
		Box100Frozen:       800,
		HasFreeShipping:    true,
		FreeShippingRates:  3000,
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
				ShippingID:         "shipping-id",
				Name:               "デフォルト配送設定",
				IsDefault:          true,
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
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateShippingInput{},
			expectErr: store.ErrInvalidArgument,
		},
		{
			name: "failed to update shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Update(ctx, "shipping-id", params).Return(assert.AnError)
			},
			input: &store.UpdateShippingInput{
				ShippingID:         "shipping-id",
				Name:               "デフォルト配送設定",
				IsDefault:          true,
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
			},
			expectErr: store.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateShipping(ctx, tt.input)
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
			input: &store.DeleteShippingInput{
				ShippingID: "shipping-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteShippingInput{},
			expectErr: store.ErrInvalidArgument,
		},
		{
			name: "failed to delete shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Delete(ctx, "shipping-id").Return(assert.AnError)
			},
			input: &store.DeleteShippingInput{
				ShippingID: "shipping-id",
			},
			expectErr: store.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteShipping(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
