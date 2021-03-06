package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/codes"
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
		Limit:  20,
		Offset: 0,
	}
	shikoku := []int64{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(len(shikoku))
	set.AddInt64s(shikoku...)
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
				Limit:  20,
				Offset: 0,
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
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list shippings",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().List(gomock.Any(), params).Return(nil, errmock)
				mocks.db.Shipping.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListShippingsInput{
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrUnknown,
		},
		{
			name: "failed to count shippings",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().List(gomock.Any(), params).Return(shippings, nil)
				mocks.db.Shipping.EXPECT().Count(gomock.Any(), params).Return(int64(0), errmock)
			},
			input: &store.ListShippingsInput{
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrUnknown,
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

func TestGetShipping(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 7, 16, 18, 30, 0, 0)
	shikoku := []int64{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(len(shikoku))
	set.AddInt64s(shikoku...)
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
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Get(ctx, "shipping-id").Return(nil, errmock)
			},
			input: &store.GetShippingInput{
				ShippingID: "shipping-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
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
	set := set.New(len(shikoku))
	set.AddInt64s(shikoku...)
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
						}
						assert.Equal(t, expect, shipping)
						return nil
					})
			},
			input: &store.CreateShippingInput{
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
			name: "failed to create shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Create(ctx, gomock.Any()).Return(errmock)
			},
			input: &store.CreateShippingInput{
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
			},
			expectErr: exception.ErrUnknown,
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
	set := set.New(len(shikoku))
	set.AddInt64s(shikoku...)
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
		Name: "デフォルト配送設定",
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
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Update(ctx, "shipping-id", params).Return(errmock)
			},
			input: &store.UpdateShippingInput{
				ShippingID:         "shipping-id",
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
			},
			expectErr: exception.ErrUnknown,
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
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shipping.EXPECT().Delete(ctx, "shipping-id").Return(errmock)
			},
			input: &store.DeleteShippingInput{
				ShippingID: "shipping-id",
			},
			expectErr: exception.ErrUnknown,
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
