package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestPromotionStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.PromotionStatus
		expect types.PromotionStatus
	}{
		{
			name:   "private",
			status: entity.PromotionStatusPrivate,
			expect: types.PromotionStatusPrivate,
		},
		{
			name:   "waiting",
			status: entity.PromotionStatusWaiting,
			expect: types.PromotionStatusWaiting,
		},
		{
			name:   "enabled",
			status: entity.PromotionStatusEnabled,
			expect: types.PromotionStatusEnabled,
		},
		{
			name:   "finisihed",
			status: entity.PromotionStatusFinished,
			expect: types.PromotionStatusFinished,
		},
		{
			name:   "unknown",
			status: entity.PromotionStatusUnknown,
			expect: types.PromotionStatusUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewPromotionStatus(tt.status)
			assert.Equal(t, tt.expect, actual.Response())
		})
	}
}

func TestDiscountType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		discountType entity.DiscountType
		expect       types.DiscountType
	}{
		{
			name:         "success to amount",
			discountType: entity.DiscountTypeAmount,
			expect:       types.DiscountTypeAmount,
		},
		{
			name:         "success to rate",
			discountType: entity.DiscountTypeRate,
			expect:       types.DiscountTypeRate,
		},
		{
			name:         "success to free shipping",
			discountType: entity.DiscountTypeFreeShipping,
			expect:       types.DiscountTypeFreeShipping,
		},
		{
			name:         "success to unknown",
			discountType: entity.DiscountTypeUnknown,
			expect:       types.DiscountTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewDiscountType(tt.discountType)
			assert.Equal(t, tt.expect, actual.Response())
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

func TestPromotionTargetType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		targetType entity.PromotionTargetType
		expect     types.PromotionTargetType
	}{
		{
			name:       "all shop",
			targetType: entity.PromotionTargetTypeAllShop,
			expect:     types.PromotionTargetTypeAllShop,
		},
		{
			name:       "only shop",
			targetType: entity.PromotionTargetTypeSpecificShop,
			expect:     types.PromotionTargetTypeSpecificShop,
		},
		{
			name:       "unknown",
			targetType: -1,
			expect:     types.PromotionTargetTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewPromotionTargetType(tt.targetType)
			assert.Equal(t, tt.expect, actual.Response())
		})
	}
}

func TestPromotion(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name      string
		promotion *entity.Promotion
		aggregate *entity.AggregatedOrderPromotion
		expect    *Promotion
	}{
		{
			name: "success",
			promotion: &entity.Promotion{
				ID:           "promotion-id",
				ShopID:       "shop-id",
				Status:       entity.PromotionStatusEnabled,
				Title:        "夏の採れたて野菜マルシェを開催!!",
				Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
				Public:       true,
				TargetType:   entity.PromotionTargetTypeSpecificShop,
				DiscountType: entity.DiscountTypeFreeShipping,
				DiscountRate: 0,
				Code:         "code0001",
				CodeType:     entity.PromotionCodeTypeOnce,
				StartAt:      now,
				EndAt:        now.AddDate(0, 1, 0),
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			aggregate: &entity.AggregatedOrderPromotion{
				PromotionID:   "promotion-id",
				OrderCount:    2,
				DiscountTotal: 1000,
			},
			expect: &Promotion{
				Promotion: types.Promotion{
					ID:           "promotion-id",
					ShopID:       "shop-id",
					Status:       types.PromotionStatusEnabled,
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					Public:       true,
					TargetType:   types.PromotionTargetTypeSpecificShop,
					DiscountType: types.DiscountTypeFreeShipping,
					DiscountRate: 0,
					Code:         "code0001",
					UsedCount:    2,
					UsedAmount:   1000,
					StartAt:      1640962800,
					EndAt:        1643641200,
					CreatedAt:    1640962800,
					UpdatedAt:    1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewPromotion(tt.promotion, tt.aggregate))
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
					Status:       types.PromotionStatusEnabled,
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					Public:       true,
					DiscountType: types.DiscountTypeFreeShipping,
					DiscountRate: 0,
					Code:         "code0001",
					UsedCount:    2,
					UsedAmount:   1000,
					StartAt:      1640962800,
					EndAt:        1643641200,
					CreatedAt:    1640962800,
					UpdatedAt:    1640962800,
				},
			},
			expect: &types.Promotion{
				ID:           "promotion-id",
				Status:       types.PromotionStatusEnabled,
				Title:        "夏の採れたて野菜マルシェを開催!!",
				Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
				Public:       true,
				DiscountType: types.DiscountTypeFreeShipping,
				DiscountRate: 0,
				Code:         "code0001",
				UsedCount:    2,
				UsedAmount:   1000,
				StartAt:      1640962800,
				EndAt:        1643641200,
				CreatedAt:    1640962800,
				UpdatedAt:    1640962800,
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
		aggregates map[string]*entity.AggregatedOrderPromotion
		expect     Promotions
	}{
		{
			name: "success",
			promotions: entity.Promotions{
				{
					ID:           "promotion-id",
					ShopID:       "shop-id",
					Status:       entity.PromotionStatusEnabled,
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					Public:       true,
					TargetType:   entity.PromotionTargetTypeSpecificShop,
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
			aggregates: map[string]*entity.AggregatedOrderPromotion{
				"promotion-id": {
					PromotionID:   "promotion-id",
					OrderCount:    2,
					DiscountTotal: 1000,
				},
			},
			expect: Promotions{
				{
					Promotion: types.Promotion{
						ID:           "promotion-id",
						ShopID:       "shop-id",
						Status:       types.PromotionStatusEnabled,
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Public:       true,
						TargetType:   types.PromotionTargetTypeSpecificShop,
						DiscountType: types.DiscountTypeFreeShipping,
						DiscountRate: 0,
						Code:         "code0001",
						UsedCount:    2,
						UsedAmount:   1000,
						StartAt:      1640962800,
						EndAt:        1643641200,
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewPromotions(tt.promotions, tt.aggregates))
		})
	}
}

func TestPromotions_ShopIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		promotions Promotions
		expect     []string
	}{
		{
			name: "success",
			promotions: Promotions{
				{
					Promotion: types.Promotion{
						ID:           "promotion-id",
						ShopID:       "shop-id",
						Status:       types.PromotionStatusEnabled,
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Public:       true,
						DiscountType: types.DiscountTypeFreeShipping,
						DiscountRate: 0,
						Code:         "code0001",
						UsedCount:    2,
						UsedAmount:   1000,
						StartAt:      1640962800,
						EndAt:        1643641200,
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
			expect: []string{"shop-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.promotions.ShopIDs())
		})
	}
}

func TestPromotions_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		promotions Promotions
		expect     map[string]*Promotion
	}{
		{
			name: "success",
			promotions: Promotions{
				{
					Promotion: types.Promotion{
						ID:           "promotion-id",
						Status:       types.PromotionStatusEnabled,
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Public:       true,
						DiscountType: types.DiscountTypeFreeShipping,
						DiscountRate: 0,
						Code:         "code0001",
						UsedCount:    2,
						UsedAmount:   1000,
						StartAt:      1640962800,
						EndAt:        1643641200,
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
			expect: map[string]*Promotion{
				"promotion-id": {
					Promotion: types.Promotion{
						ID:           "promotion-id",
						Status:       types.PromotionStatusEnabled,
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Public:       true,
						DiscountType: types.DiscountTypeFreeShipping,
						DiscountRate: 0,
						Code:         "code0001",
						UsedCount:    2,
						UsedAmount:   1000,
						StartAt:      1640962800,
						EndAt:        1643641200,
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.promotions.Map())
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
						Status:       types.PromotionStatusEnabled,
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Public:       true,
						DiscountType: types.DiscountTypeFreeShipping,
						DiscountRate: 0,
						Code:         "code0001",
						UsedCount:    2,
						UsedAmount:   1000,
						StartAt:      1640962800,
						EndAt:        1643641200,
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
			expect: []*types.Promotion{
				{
					ID:           "promotion-id",
					Status:       types.PromotionStatusEnabled,
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					Public:       true,
					DiscountType: types.DiscountTypeFreeShipping,
					DiscountRate: 0,
					Code:         "code0001",
					UsedCount:    2,
					UsedAmount:   1000,
					StartAt:      1640962800,
					EndAt:        1643641200,
					CreatedAt:    1640962800,
					UpdatedAt:    1640962800,
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
