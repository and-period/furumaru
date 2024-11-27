package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListPromotions(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 8, 13, 18, 30, 0, 0)
	params := &database.ListPromotionsParams{
		Limit:  30,
		Offset: 0,
		Orders: []*database.ListPromotionsOrder{
			{Key: database.ListPromotionsOrderByPublic, OrderByASC: true},
		},
	}
	promotions := entity.Promotions{
		{
			ID:           "promotion-id",
			Title:        "夏の採れたて野菜マルシェを開催!!",
			Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
			Public:       true,
			PublishedAt:  now,
			DiscountType: entity.DiscountTypeFreeShipping,
			DiscountRate: 0,
			Code:         "code0001",
			CodeType:     entity.PromotionCodeTypeOnce,
			StartAt:      now,
			EndAt:        now.AddDate(0, 1, 0),
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListPromotionsInput
		expect      entity.Promotions
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().List(gomock.Any(), params).Return(promotions, nil)
				mocks.db.Promotion.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListPromotionsInput{
				Limit:  30,
				Offset: 0,
				Orders: []*store.ListPromotionsOrder{
					{Key: store.ListPromotionsOrderByPublic, OrderByASC: true},
				},
			},
			expect:      promotions,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &store.ListPromotionsInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list promotions",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Promotion.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListPromotionsInput{
				Limit:  30,
				Offset: 0,
				Orders: []*store.ListPromotionsOrder{
					{Key: store.ListPromotionsOrderByPublic, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count promotions",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().List(gomock.Any(), params).Return(promotions, nil)
				mocks.db.Promotion.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListPromotionsInput{
				Limit:  30,
				Offset: 0,
				Orders: []*store.ListPromotionsOrder{
					{Key: store.ListPromotionsOrderByPublic, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListPromotions(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestMultiGetPromotions(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 8, 13, 18, 30, 0, 0)
	promotions := entity.Promotions{
		{
			ID:           "promotion-id",
			Title:        "夏の採れたて野菜マルシェを開催!!",
			Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
			Public:       true,
			PublishedAt:  now,
			DiscountType: entity.DiscountTypeFreeShipping,
			DiscountRate: 0,
			Code:         "code0001",
			CodeType:     entity.PromotionCodeTypeOnce,
			StartAt:      now,
			EndAt:        now.AddDate(0, 1, 0),
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetPromotionsInput
		expect    entity.Promotions
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().MultiGet(ctx, []string{"promotion-id"}).Return(promotions, nil)
			},
			input: &store.MultiGetPromotionsInput{
				PromotionIDs: []string{"promotion-id"},
			},
			expect:    promotions,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetPromotionsInput{
				PromotionIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get promotions",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().MultiGet(ctx, []string{"promotion-id"}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetPromotionsInput{
				PromotionIDs: []string{"promotion-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetPromotions(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetPromotion(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 8, 13, 18, 30, 0, 0)
	promotion := &entity.Promotion{
		ID:           "promotion-id",
		Title:        "夏の採れたて野菜マルシェを開催!!",
		Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
		Public:       true,
		PublishedAt:  now,
		DiscountType: entity.DiscountTypeFreeShipping,
		DiscountRate: 0,
		Code:         "code0001",
		CodeType:     entity.PromotionCodeTypeOnce,
		StartAt:      now.AddDate(0, -1, 0),
		EndAt:        now.AddDate(0, 1, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetPromotionInput
		expect    *entity.Promotion
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().Get(ctx, "promotion-id").Return(promotion, nil)
			},
			input: &store.GetPromotionInput{
				PromotionID: "promotion-id",
			},
			expect:    promotion,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetPromotionInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get promotion",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().Get(ctx, "promotion-id").Return(nil, assert.AnError)
			},
			input: &store.GetPromotionInput{
				PromotionID: "promotion-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get enabled promotion",
			setup: func(ctx context.Context, mocks *mocks) {
				promotion := &entity.Promotion{}
				mocks.db.Promotion.EXPECT().Get(ctx, "promotion-id").Return(promotion, nil)
			},
			input: &store.GetPromotionInput{
				PromotionID: "promotion-id",
				OnlyEnabled: true,
			},
			expect:    nil,
			expectErr: exception.ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetPromotion(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}, withNow(now)))
	}
}

func TestGetPromotionByCode(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 8, 13, 18, 30, 0, 0)
	promotion := &entity.Promotion{
		ID:           "promotion-id",
		Title:        "夏の採れたて野菜マルシェを開催!!",
		Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
		Public:       true,
		PublishedAt:  now,
		DiscountType: entity.DiscountTypeFreeShipping,
		DiscountRate: 0,
		Code:         "code0001",
		CodeType:     entity.PromotionCodeTypeOnce,
		StartAt:      now.AddDate(0, -1, 0),
		EndAt:        now.AddDate(0, 1, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetPromotionByCodeInput
		expect    *entity.Promotion
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().GetByCode(ctx, "code0001").Return(promotion, nil)
			},
			input: &store.GetPromotionByCodeInput{
				PromotionCode: "code0001",
			},
			expect:    promotion,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetPromotionByCodeInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get promotion",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().GetByCode(ctx, "code0001").Return(nil, assert.AnError)
			},
			input: &store.GetPromotionByCodeInput{
				PromotionCode: "code0001",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get enabled promotion",
			setup: func(ctx context.Context, mocks *mocks) {
				promotion := &entity.Promotion{}
				mocks.db.Promotion.EXPECT().GetByCode(ctx, "code0001").Return(promotion, nil)
			},
			input: &store.GetPromotionByCodeInput{
				PromotionCode: "code0001",
				OnlyEnabled:   true,
			},
			expect:    nil,
			expectErr: exception.ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetPromotionByCode(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}, withNow(now)))
	}
}

func TestCreatePromotion(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 8, 1, 0, 0, 0, 0)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreatePromotionInput
		expect    *entity.Promotion
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, promotion *entity.Promotion) error {
						expect := &entity.Promotion{
							ID:           promotion.ID, // ignore
							Title:        "プロモーションタイトル",
							Description:  "プロモーションの詳細です。",
							Public:       true,
							PublishedAt:  now,
							DiscountType: entity.DiscountTypeRate,
							DiscountRate: 10,
							Code:         "excode01",
							CodeType:     entity.PromotionCodeTypeAlways,
							StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
							EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
						}
						assert.Equal(t, expect, promotion)
						return nil
					})
			},
			input: &store.CreatePromotionInput{
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				DiscountType: entity.DiscountTypeRate,
				DiscountRate: 10,
				Code:         "excode01",
				CodeType:     entity.PromotionCodeTypeAlways,
				StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreatePromotionInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create promotion",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreatePromotionInput{
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				DiscountType: entity.DiscountTypeRate,
				DiscountRate: 10,
				Code:         "excode01",
				CodeType:     entity.PromotionCodeTypeAlways,
				StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreatePromotion(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}, withNow(now)))
	}
}

func TestUpdatePromotion(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 8, 1, 0, 0, 0, 0)
	params := &database.UpdatePromotionParams{
		Title:        "プロモーションタイトル",
		Description:  "プロモーションの詳細です。",
		Public:       true,
		PublishedAt:  now,
		DiscountType: entity.DiscountTypeRate,
		DiscountRate: 10,
		Code:         "excode01",
		CodeType:     entity.PromotionCodeTypeAlways,
		StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
		EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdatePromotionInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().Update(ctx, "promotion-id", params).Return(nil)
			},
			input: &store.UpdatePromotionInput{
				PromotionID:  "promotion-id",
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				DiscountType: entity.DiscountTypeRate,
				DiscountRate: 10,
				Code:         "excode01",
				CodeType:     entity.PromotionCodeTypeAlways,
				StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdatePromotionInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update promotion",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().Update(ctx, "promotion-id", params).Return(assert.AnError)
			},
			input: &store.UpdatePromotionInput{
				PromotionID:  "promotion-id",
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				DiscountType: entity.DiscountTypeRate,
				DiscountRate: 10,
				Code:         "excode01",
				CodeType:     entity.PromotionCodeTypeAlways,
				StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdatePromotion(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}, withNow(now)))
	}
}

func TestDeletePromotion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeletePromotionInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().Delete(ctx, "promotion-id").Return(nil)
			},
			input: &store.DeletePromotionInput{
				PromotionID: "promotion-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeletePromotionInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete promotion",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Promotion.EXPECT().Delete(ctx, "promotion-id").Return(assert.AnError)
			},
			input: &store.DeletePromotionInput{
				PromotionID: "promotion-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeletePromotion(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
