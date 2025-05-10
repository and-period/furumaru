package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
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
			expect: PromotionStatusPrivate,
		},
		{
			name:   "waiting",
			status: entity.PromotionStatusWaiting,
			expect: PromotionStatusWaiting,
		},
		{
			name:   "enabled",
			status: entity.PromotionStatusEnabled,
			expect: PromotionStatusEnabled,
		},
		{
			name:   "finisihed",
			status: entity.PromotionStatusFinished,
			expect: PromotionStatusFinished,
		},
		{
			name:   "unknown",
			status: entity.PromotionStatusUnknown,
			expect: PromotionStatusUnknown,
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
		expect int32
	}{
		{
			name:   "success",
			status: PromotionStatusEnabled,
			expect: 3,
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
			expect:       DiscountTypeAmount,
		},
		{
			name:         "success to rate",
			discountType: entity.DiscountTypeRate,
			expect:       DiscountTypeRate,
		},
		{
			name:         "success to free shipping",
			discountType: entity.DiscountTypeFreeShipping,
			expect:       DiscountTypeFreeShipping,
		},
		{
			name:         "success to unknown",
			discountType: entity.DiscountTypeUnknown,
			expect:       DiscountTypeUnknown,
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
			discountType: DiscountTypeAmount,
			expect:       entity.DiscountTypeAmount,
		},
		{
			name:         "success to rate",
			discountType: DiscountTypeRate,
			expect:       entity.DiscountTypeRate,
		},
		{
			name:         "success to free shipping",
			discountType: DiscountTypeFreeShipping,
			expect:       entity.DiscountTypeFreeShipping,
		},
		{
			name:         "success to unknown",
			discountType: DiscountTypeUnknown,
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
		expect       int32
	}{
		{
			name:         "success",
			discountType: DiscountTypeAmount,
			expect:       1,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.discountType.Response())
		})
	}
}

func TestPromotionTargetType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		targetType entity.PromotionTargetType
		expect     PromotionTargetType
		response   int32
	}{
		{
			name:       "all shop",
			targetType: entity.PromotionTargetTypeAllShop,
			expect:     PromotionTargetTypeAllShop,
			response:   1,
		},
		{
			name:       "only shop",
			targetType: entity.PromotionTargetTypeSpecificShop,
			expect:     PromotionTargetTypeSpecificShop,
			response:   2,
		},
		{
			name:       "unknown",
			targetType: -1,
			expect:     PromotionTargetTypeUnknown,
			response:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewPromotionTargetType(tt.targetType)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
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
				Promotion: response.Promotion{
					ID:           "promotion-id",
					ShopID:       "shop-id",
					Status:       int32(PromotionStatusEnabled),
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					Public:       true,
					TargetType:   PromotionTargetTypeSpecificShop.Response(),
					DiscountType: DiscountTypeFreeShipping.Response(),
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
		expect    *response.Promotion
	}{
		{
			name: "success",
			promotion: &Promotion{
				Promotion: response.Promotion{
					ID:           "promotion-id",
					Status:       int32(PromotionStatusEnabled),
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					Public:       true,
					DiscountType: DiscountTypeFreeShipping.Response(),
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
			expect: &response.Promotion{
				ID:           "promotion-id",
				Status:       int32(PromotionStatusEnabled),
				Title:        "夏の採れたて野菜マルシェを開催!!",
				Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
				Public:       true,
				DiscountType: DiscountTypeFreeShipping.Response(),
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
					Promotion: response.Promotion{
						ID:           "promotion-id",
						ShopID:       "shop-id",
						Status:       int32(PromotionStatusEnabled),
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Public:       true,
						TargetType:   PromotionTargetTypeSpecificShop.Response(),
						DiscountType: DiscountTypeFreeShipping.Response(),
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
					Promotion: response.Promotion{
						ID:           "promotion-id",
						ShopID:       "shop-id",
						Status:       int32(PromotionStatusEnabled),
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Public:       true,
						DiscountType: DiscountTypeFreeShipping.Response(),
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
					Promotion: response.Promotion{
						ID:           "promotion-id",
						Status:       int32(PromotionStatusEnabled),
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Public:       true,
						DiscountType: DiscountTypeFreeShipping.Response(),
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
					Promotion: response.Promotion{
						ID:           "promotion-id",
						Status:       int32(PromotionStatusEnabled),
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Public:       true,
						DiscountType: DiscountTypeFreeShipping.Response(),
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
		expect     []*response.Promotion
	}{
		{
			name: "success",
			promotions: Promotions{
				{
					Promotion: response.Promotion{
						ID:           "promotion-id",
						Status:       int32(PromotionStatusEnabled),
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Public:       true,
						DiscountType: DiscountTypeFreeShipping.Response(),
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
			expect: []*response.Promotion{
				{
					ID:           "promotion-id",
					Status:       int32(PromotionStatusEnabled),
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					Public:       true,
					DiscountType: DiscountTypeFreeShipping.Response(),
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
