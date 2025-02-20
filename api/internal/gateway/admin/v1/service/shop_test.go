package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/stretchr/testify/assert"
)

func TestShops(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		shops    entity.Shops
		expect   Shops
		response []*response.Shop
	}{
		{
			name: "success",
			shops: entity.Shops{
				{
					ID:             "shop-id",
					CoordinatorID:  "coordinator-id",
					ProducerIDs:    []string{"producer-id"},
					ProductTypeIDs: []string{"product-type-id"},
					BusinessDays:   []time.Weekday{time.Monday},
					Name:           "テスト店舗",
					Activated:      true,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
			expect: Shops{
				{
					Shop: response.Shop{
						ID:             "shop-id",
						Name:           "テスト店舗",
						CoordinatorID:  "coordinator-id",
						ProducerIDs:    []string{"producer-id"},
						ProductTypeIDs: []string{"product-type-id"},
						BusinessDays:   []time.Weekday{time.Monday},
						CreatedAt:      now.Unix(),
						UpdatedAt:      now.Unix(),
					},
				},
			},
			response: []*response.Shop{
				{
					ID:             "shop-id",
					Name:           "テスト店舗",
					CoordinatorID:  "coordinator-id",
					ProducerIDs:    []string{"producer-id"},
					ProductTypeIDs: []string{"product-type-id"},
					BusinessDays:   []time.Weekday{time.Monday},
					CreatedAt:      now.Unix(),
					UpdatedAt:      now.Unix(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewShops(tt.shops)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
		})
	}
}

func TestShop_GetID(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name   string
		shop   *Shop
		expect string
	}{
		{
			name: "success",
			shop: &Shop{
				Shop: response.Shop{
					ID:             "shop-id",
					Name:           "テスト店舗",
					CoordinatorID:  "coordinator-id",
					ProducerIDs:    []string{"producer-id"},
					ProductTypeIDs: []string{"product-type-id"},
					BusinessDays:   []time.Weekday{time.Monday},
					CreatedAt:      now.Unix(),
					UpdatedAt:      now.Unix(),
				},
			},
			expect: "shop-id",
		},
		{
			name:   "empty",
			shop:   nil,
			expect: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.shop.GetID())
		})
	}
}
