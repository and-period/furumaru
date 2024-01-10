package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

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
		tt := tt
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
		tt := tt
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
		tt := tt
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
			expect: &Promotion{
				Promotion: response.Promotion{
					ID:           "promotion-id",
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					DiscountType: DiscountTypeFreeShipping.Response(),
					DiscountRate: 0,
					Code:         "code0001",
					StartAt:      1640962800,
					EndAt:        1643641200,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
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
		expect    *response.Promotion
	}{
		{
			name: "success",
			promotion: &Promotion{
				Promotion: response.Promotion{
					ID:           "promotion-id",
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					DiscountType: DiscountTypeFreeShipping.Response(),
					DiscountRate: 0,
					Code:         "code0001",
					StartAt:      1640962800,
					EndAt:        1643641200,
				},
			},
			expect: &response.Promotion{
				ID:           "promotion-id",
				Title:        "夏の採れたて野菜マルシェを開催!!",
				Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
				DiscountType: DiscountTypeFreeShipping.Response(),
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
		tt := tt
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
			},
			expect: Promotions{
				{
					Promotion: response.Promotion{
						ID:           "promotion-id",
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						DiscountType: DiscountTypeFreeShipping.Response(),
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
		tt := tt
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
		expect     []*response.Promotion
	}{
		{
			name: "success",
			promotions: Promotions{
				{
					Promotion: response.Promotion{
						ID:           "promotion-id",
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						DiscountType: DiscountTypeFreeShipping.Response(),
						DiscountRate: 0,
						Code:         "code0001",
						StartAt:      1640962800,
						EndAt:        1643641200,
					},
				},
			},
			expect: []*response.Promotion{
				{
					ID:           "promotion-id",
					Title:        "夏の採れたて野菜マルシェを開催!!",
					Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
					DiscountType: DiscountTypeFreeShipping.Response(),
					DiscountRate: 0,
					Code:         "code0001",
					StartAt:      1640962800,
					EndAt:        1643641200,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.promotions.Response())
		})
	}
}
