package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestPromotionStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.PromotionStatus
		expect PromotionStatus
	}{
		{
			name:   "private",
			status: entity.PromotionStatusPrivate,
			expect: PromotionStatus(types.PromotionStatusPrivate),
		},
		{
			name:   "waiting",
			status: entity.PromotionStatusWaiting,
			expect: PromotionStatus(types.PromotionStatusWaiting),
		},
		{
			name:   "enabled",
			status: entity.PromotionStatusEnabled,
			expect: PromotionStatus(types.PromotionStatusEnabled),
		},
		{
			name:   "finisihed",
			status: entity.PromotionStatusFinished,
			expect: PromotionStatus(types.PromotionStatusFinished),
		},
		{
			name:   "unknown",
			status: entity.PromotionStatusUnknown,
			expect: PromotionStatus(types.PromotionStatusUnknown),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewPromotionStatus(tt.status)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestPromotionStatus_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status PromotionStatus
		expect types.PromotionStatus
	}{
		{
			name:   "success",
			status: PromotionStatus(types.PromotionStatusEnabled),
			expect: types.PromotionStatusEnabled,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.status.Response())
		})
	}
}

func TestDiscountType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		discountType entity.DiscountType
		expect       DiscountType
	}{
		{
			name:         "success to amount",
			discountType: entity.DiscountTypeAmount,
			expect:       DiscountType(types.DiscountTypeAmount),
		},
		{
			name:         "success to rate",
			discountType: entity.DiscountTypeRate,
			expect:       DiscountType(types.DiscountTypeRate),
		},
		{
			name:         "success to free shipping",
			discountType: entity.DiscountTypeFreeShipping,
			expect:       DiscountType(types.DiscountTypeFreeShipping),
		},
		{
			name:         "success to unknown",
			discountType: entity.DiscountTypeUnknown,
			expect:       DiscountType(types.DiscountTypeUnknown),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewDiscountType(tt.discountType))
		})
	}
}

func TestDiscountType_StoreEntity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		discountType DiscountType
		expect       entity.DiscountType
	}{
		{
			name:         "success to amount",
			discountType: DiscountType(types.DiscountTypeAmount),
			expect:       entity.DiscountTypeAmount,
		},
		{
			name:         "success to rate",
			discountType: DiscountType(types.DiscountTypeRate),
			expect:       entity.DiscountTypeRate,
		},
		{
			name:         "success to free shipping",
			discountType: DiscountType(types.DiscountTypeFreeShipping),
			expect:       entity.DiscountTypeFreeShipping,
		},
		{
			name:         "success to unknown",
			discountType: DiscountType(types.DiscountTypeUnknown),
			expect:       entity.DiscountTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.discountType.StoreEntity())
		})
	}
}

func TestDiscountType_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		discountType DiscountType
		expect       types.DiscountType
	}{
		{
			name:         "success",
			discountType: DiscountType(types.DiscountTypeAmount),
			expect:       types.DiscountTypeAmount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.discountType.Response())
		})
	}
}

func TestPromotion(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name      string
		promotion *entity.Promotion
		expect    *Promotion
	}{
		{
			name: "success",
			promotion: &entity.Promotion{
				ID:           "promotion-id",
				Title:        "夏の採れたて野菜マルシェを開催!!",
				Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
				Status:       entity.PromotionStatusEnabled,
				Public:       true,
				DiscountType: entity.DiscountTypeFreeShipping,
				DiscountRate: 0,
				Code:         "code0001",
				CodeType:     entity.PromotionCodeTypeOnce,
				StartAt:      now,
				EndAt:        now.AddDate(0, 1, 0),
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			expect: &Promotion{
				Promotion: types.Promotion{
					ID:           "promotion-id",
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					Status:       types.PromotionStatusEnabled,
					DiscountType: NewDiscountType(entity.DiscountTypeFreeShipping).Response(),
					DiscountRate: 0,
					Code:         "code0001",
					StartAt:      1640962800,
					EndAt:        1643641200,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewPromotion(tt.promotion))
		})
	}
}

func TestPromotion_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		promotion *Promotion
		expect    *types.Promotion
	}{
		{
			name: "success",
			promotion: &Promotion{
				Promotion: types.Promotion{
					ID:           "promotion-id",
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					Status:       types.PromotionStatusEnabled,
					DiscountType: NewDiscountType(entity.DiscountTypeFreeShipping).Response(),
					DiscountRate: 0,
					Code:         "code0001",
					StartAt:      1640962800,
					EndAt:        1643641200,
				},
			},
			expect: &types.Promotion{
				ID:           "promotion-id",
				Title:        "夏の採れたて野菜マルシェを開催!!",
				Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
				Status:       NewPromotionStatus(entity.PromotionStatusEnabled).Response(),
				DiscountType: NewDiscountType(entity.DiscountTypeFreeShipping).Response(),
				DiscountRate: 0,
				Code:         "code0001",
				StartAt:      1640962800,
				EndAt:        1643641200,
			},
		},
		{
			name:      "empty",
			promotion: nil,
			expect:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.promotion.Response())
		})
	}
}

func TestPromotions(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name       string
		promotions entity.Promotions
		expect     Promotions
	}{
		{
			name: "success",
			promotions: entity.Promotions{
				{
					ID:           "promotion-id",
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					Status:       entity.PromotionStatusEnabled,
					Public:       true,
					DiscountType: entity.DiscountTypeFreeShipping,
					DiscountRate: 0,
					Code:         "code0001",
					CodeType:     entity.PromotionCodeTypeOnce,
					StartAt:      now,
					EndAt:        now.AddDate(0, 1, 0),
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
			expect: Promotions{
				{
					Promotion: types.Promotion{
						ID:           "promotion-id",
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Status:       types.PromotionStatusEnabled,
						DiscountType: NewDiscountType(entity.DiscountTypeFreeShipping).Response(),
						DiscountRate: 0,
						Code:         "code0001",
						StartAt:      1640962800,
						EndAt:        1643641200,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewPromotions(tt.promotions))
		})
	}
}

func TestPromotions_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		promotions Promotions
		expect     []*types.Promotion
	}{
		{
			name: "success",
			promotions: Promotions{
				{
					Promotion: types.Promotion{
						ID:           "promotion-id",
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Status:       types.PromotionStatusEnabled,
						DiscountType: NewDiscountType(entity.DiscountTypeFreeShipping).Response(),
						DiscountRate: 0,
						Code:         "code0001",
						StartAt:      1640962800,
						EndAt:        1643641200,
					},
				},
			},
			expect: []*types.Promotion{
				{
					ID:           "promotion-id",
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					Status:       types.PromotionStatusEnabled,
					DiscountType: NewDiscountType(entity.DiscountTypeFreeShipping).Response(),
					DiscountRate: 0,
					Code:         "code0001",
					StartAt:      1640962800,
					EndAt:        1643641200,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.promotions.Response())
		})
	}
}
