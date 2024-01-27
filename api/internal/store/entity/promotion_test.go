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
			name: "success",
			params: &NewPromotionParams{
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				PublishedAt:  jst.Date(2022, 8, 9, 18, 30, 0, 0),
				DiscountType: DiscountTypeRate,
				DiscountRate: 10,
				Code:         "excode01",
				CodeType:     PromotionCodeTypeAlways,
				StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
			expect: &Promotion{
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				PublishedAt:  jst.Date(2022, 8, 9, 18, 30, 0, 0),
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
		expect    bool
	}{
		{
			name: "enabled",
			promotion: &Promotion{
				Status: PromotionStatusEnabled,
			},
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
			name:      "empty",
			promotion: nil,
			expect:    false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.promotion.IsEnabled())
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
				PublishedAt:  jst.Date(2022, 8, 9, 18, 30, 0, 0),
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
				PublishedAt:  jst.Date(2022, 8, 9, 18, 30, 0, 0),
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
				PublishedAt:  jst.Date(2022, 8, 9, 18, 30, 0, 0),
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
