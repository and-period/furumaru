package entity

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestPromotion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewPromotionParams
		expect *Promotion
	}{
		{
			name: "success for all",
			params: &NewPromotionParams{
				ShopID:       "",
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				DiscountType: DiscountTypeRate,
				DiscountRate: 10,
				Code:         "excode01",
				CodeType:     PromotionCodeTypeAlways,
				StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
			expect: &Promotion{
				ShopID:       "",
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				TargetType:   PromotionTargetTypeAllShop,
				DiscountType: DiscountTypeRate,
				DiscountRate: 10,
				Code:         "excode01",
				CodeType:     PromotionCodeTypeAlways,
				StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
		},
		{
			name: "success for shop",
			params: &NewPromotionParams{
				ShopID:       "shop-id",
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				DiscountType: DiscountTypeRate,
				DiscountRate: 10,
				Code:         "excode01",
				CodeType:     PromotionCodeTypeAlways,
				StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
			expect: &Promotion{
				ShopID:       "shop-id",
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				TargetType:   PromotionTargetTypeSpecificShop,
				DiscountType: DiscountTypeRate,
				DiscountRate: 10,
				Code:         "excode01",
				CodeType:     PromotionCodeTypeAlways,
				StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewPromotion(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestPromotion_CalcDiscount(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		promotion   *Promotion
		total       int64
		shippingFee int64
		expect      int64
	}{
		{
			name:      "empty",
			promotion: nil,
			expect:    0,
		},
		{
			name: "金額固定割引 割引金額未満の支払い",
			promotion: &Promotion{
				DiscountType: DiscountTypeAmount,
				DiscountRate: 500,
			},
			total:       300,
			shippingFee: 500,
			expect:      300,
		},
		{
			name: "金額固定割引 割引金額以上の支払い",
			promotion: &Promotion{
				DiscountType: DiscountTypeAmount,
				DiscountRate: 500,
			},
			total:       1980,
			shippingFee: 500,
			expect:      500,
		},
		{
			name: "料率指定での割引",
			promotion: &Promotion{
				DiscountType: DiscountTypeRate,
				DiscountRate: 10,
			},
			total:       1980,
			shippingFee: 500,
			expect:      198,
		},
		{
			name: "料率指定での割引 割引率が0%",
			promotion: &Promotion{
				DiscountType: DiscountTypeRate,
				DiscountRate: 0,
			},
			total:       1980,
			shippingFee: 500,
			expect:      0,
		},
		{
			name: "送料無料",
			promotion: &Promotion{
				DiscountType: DiscountTypeFreeShipping,
				DiscountRate: 0,
			},
			total:       1980,
			shippingFee: 500,
			expect:      500,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.promotion.CalcDiscount(tt.total, tt.shippingFee)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestPromotion_IsEnabled(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		promotion *Promotion
		shopID    string
		expect    bool
	}{
		{
			name: "enabled for all",
			promotion: &Promotion{
				Status:     PromotionStatusEnabled,
				TargetType: PromotionTargetTypeAllShop,
			},
			shopID: "shop-id",
			expect: true,
		},
		{
			name: "enabled specific shop",
			promotion: &Promotion{
				Status:     PromotionStatusEnabled,
				TargetType: PromotionTargetTypeSpecificShop,
				ShopID:     "shop-id",
			},
			shopID: "shop-id",
			expect: true,
		},
		{
			name: "disabled",
			promotion: &Promotion{
				Status: PromotionStatusPrivate,
			},
			expect: false,
		},
		{
			name: "enabled specific shop, but wrong shop id",
			promotion: &Promotion{
				Status:     PromotionStatusEnabled,
				TargetType: PromotionTargetTypeSpecificShop,
				ShopID:     "other-id",
			},
			shopID: "shop-id",
			expect: false,
		},
		{
			name:      "empty",
			promotion: nil,
			shopID:    "shop-id",
			expect:    false,
		},
		{
			name: "unknown target type",
			promotion: &Promotion{
				Status:     PromotionStatusEnabled,
				TargetType: -1,
			},
			shopID: "shop-id",
			expect: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.promotion.IsEnabled(tt.shopID))
		})
	}
}

func TestPromotion_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		promotion *Promotion
		expect    error
	}{
		{
			name: "success",
			promotion: &Promotion{
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				DiscountType: DiscountTypeRate,
				DiscountRate: 10,
				Code:         "excode01",
				CodeType:     PromotionCodeTypeAlways,
				StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
			expect: nil,
		},
		{
			name: "amount error",
			promotion: &Promotion{
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				DiscountType: DiscountTypeAmount,
				DiscountRate: 0,
				Code:         "excode01",
				CodeType:     PromotionCodeTypeAlways,
				StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
			expect: errInvalidDiscount,
		},
		{
			name: "rate error",
			promotion: &Promotion{
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				DiscountType: DiscountTypeRate,
				DiscountRate: 0,
				Code:         "excode01",
				CodeType:     PromotionCodeTypeAlways,
				StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
			expect: errInvalidDiscount,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ErrorIs(t, tt.promotion.Validate(), tt.expect)
		})
	}
}

func TestPromotions_IDs(t *testing.T) {
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
					ID:           "promotion-id",
					Status:       PromotionStatusEnabled,
					Title:        "プロモーションタイトル",
					Description:  "プロモーションの詳細です。",
					Public:       true,
					DiscountType: DiscountTypeRate,
					DiscountRate: 0,
					Code:         "excode01",
					CodeType:     PromotionCodeTypeAlways,
					StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
				},
			},
			expect: []string{"promotion-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.promotions.IDs())
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
					ID:           "promotion-id",
					ShopID:       "shop-id",
					Status:       PromotionStatusEnabled,
					Title:        "プロモーションタイトル",
					Description:  "プロモーションの詳細です。",
					Public:       true,
					DiscountType: DiscountTypeRate,
					DiscountRate: 0,
					Code:         "excode01",
					CodeType:     PromotionCodeTypeAlways,
					StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
				},
			},
			expect: []string{"shop-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.promotions.ShopIDs())
		})
	}
}

func TestPromotions_Fill(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name       string
		promotions Promotions
		expect     PromotionStatus
	}{
		{
			name: "private",
			promotions: Promotions{
				{
					Public:  false,
					StartAt: now.Add(-time.Hour),
					EndAt:   now.Add(time.Hour),
				},
			},
			expect: PromotionStatusPrivate,
		},
		{
			name: "waiting",
			promotions: Promotions{
				{
					Public:  true,
					StartAt: now.Add(time.Hour),
					EndAt:   now.Add(2 * time.Hour),
				},
			},
			expect: PromotionStatusWaiting,
		},
		{
			name: "enabled",
			promotions: Promotions{
				{
					Public:  true,
					StartAt: now.Add(-time.Hour),
					EndAt:   now.Add(time.Hour),
				},
			},
			expect: PromotionStatusEnabled,
		},
		{
			name: "finished",
			promotions: Promotions{
				{
					Public:  true,
					StartAt: now.Add(-2 * time.Hour),
					EndAt:   now.Add(-time.Hour),
				},
			},
			expect: PromotionStatusFinished,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.promotions.Fill(now)
			assert.Equal(t, tt.expect, tt.promotions[0].Status)
		})
	}
}
