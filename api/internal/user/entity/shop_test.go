package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewShop(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name   string
		params *ShopParams
		expect *Shop
	}{
		{
			name: "success",
			params: &ShopParams{
				CoordinatorID:  "coordinator-id",
				Name:           "テスト店舗",
				ProductTypeIDs: []string{"product-type-id01", "product-type-id02"},
				BusinessDays:   []time.Weekday{time.Monday, time.Tuesday, time.Wednesday},
			},
			expect: &Shop{
				CoordinatorID:  "coordinator-id",
				Name:           "テスト店舗",
				ProductTypeIDs: []string{"product-type-id01", "product-type-id02"},
				BusinessDays:   []time.Weekday{time.Monday, time.Tuesday, time.Wednesday},
				Activated:      true,
			},
		},
		{
			name: "success with empty product types",
			params: &ShopParams{
				CoordinatorID:  "coordinator-id",
				Name:           "テスト店舗",
				ProductTypeIDs: []string{},
				BusinessDays:   []time.Weekday{},
			},
			expect: &Shop{
				CoordinatorID:  "coordinator-id",
				Name:           "テスト店舗",
				ProductTypeIDs: []string{},
				BusinessDays:   []time.Weekday{},
				Activated:      true,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewShop(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestShop_Enabled(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name   string
		shop   *Shop
		expect bool
	}{
		{
			name: "enabled when activated and not deleted",
			shop: &Shop{
				Activated: true,
				DeletedAt: gorm.DeletedAt{},
			},
			expect: true,
		},
		{
			name: "disabled when not activated",
			shop: &Shop{
				Activated: false,
				DeletedAt: gorm.DeletedAt{},
			},
			expect: false,
		},
		{
			name: "disabled when deleted",
			shop: &Shop{
				Activated: true,
				DeletedAt: gorm.DeletedAt{
					Time:  time.Now(),
					Valid: true,
				},
			},
			expect: false,
		},
		{
			name: "disabled when not activated and deleted",
			shop: &Shop{
				Activated: false,
				DeletedAt: gorm.DeletedAt{
					Time:  time.Now(),
					Valid: true,
				},
			},
			expect: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.shop.Enabled())
		})
	}
}

func TestShop_Fill(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name      string
		shop      *Shop
		producers ShopProducers
		expect    []string
	}{
		{
			name: "success",
			shop: &Shop{
				ID: "shop-id",
			},
			producers: ShopProducers{
				{ProducerID: "producer-id01"},
				{ProducerID: "producer-id02"},
			},
			expect: []string{"producer-id01", "producer-id02"},
		},
		{
			name: "empty producers",
			shop: &Shop{
				ID: "shop-id",
			},
			producers: ShopProducers{},
			expect:    []string{},
		},
		{
			name:      "nil producers",
			shop:      &Shop{ID: "shop-id"},
			producers: nil,
			expect:    []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.shop.Fill(tt.producers)
			assert.Equal(t, tt.expect, tt.shop.ProducerIDs)
		})
	}
}

func TestShops_IDs(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name   string
		shops  Shops
		expect []string
	}{
		{
			name: "success",
			shops: Shops{
				{ID: "shop-id01"},
				{ID: "shop-id02"},
				{ID: "shop-id03"},
			},
			expect: []string{"shop-id01", "shop-id02", "shop-id03"},
		},
		{
			name: "duplicated",
			shops: Shops{
				{ID: "shop-id01"},
				{ID: "shop-id02"},
				{ID: "shop-id01"},
			},
			expect: []string{"shop-id01", "shop-id02"},
		},
		{
			name:   "empty",
			shops:  Shops{},
			expect: []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.shops.IDs())
		})
	}
}

func TestShops_ProductTypeIDs(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name   string
		shops  Shops
		expect []string
	}{
		{
			name: "success",
			shops: Shops{
				{ID: "shop-id01", ProductTypeIDs: []string{"type-id01", "type-id02"}},
				{ID: "shop-id02", ProductTypeIDs: []string{"type-id02", "type-id03"}},
				{ID: "shop-id03", ProductTypeIDs: []string{"type-id03", "type-id04"}},
			},
			expect: []string{"type-id01", "type-id02", "type-id03", "type-id04"},
		},
		{
			name: "duplicated product types",
			shops: Shops{
				{ID: "shop-id01", ProductTypeIDs: []string{"type-id01", "type-id01"}},
				{ID: "shop-id02", ProductTypeIDs: []string{"type-id01"}},
			},
			expect: []string{"type-id01"},
		},
		{
			name: "empty product types",
			shops: Shops{
				{ID: "shop-id01", ProductTypeIDs: []string{}},
				{ID: "shop-id02"},
			},
			expect: []string{},
		},
		{
			name:   "empty shops",
			shops:  Shops{},
			expect: []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.shops.ProductTypeIDs())
		})
	}
}

func TestShops_MapByCoordinatorID(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name   string
		shops  Shops
		expect map[string]*Shop
	}{
		{
			name: "success",
			shops: Shops{
				{
					ID:            "shop-id01",
					CoordinatorID: "coordinator-id01",
					ProducerIDs:   []string{"producer-id01", "producer-id02"},
				},
				{
					ID:            "shop-id02",
					CoordinatorID: "coordinator-id02",
					ProducerIDs:   []string{"producer-id01", "producer-id03"},
				},
			},
			expect: map[string]*Shop{
				"coordinator-id01": {
					ID:            "shop-id01",
					CoordinatorID: "coordinator-id01",
					ProducerIDs:   []string{"producer-id01", "producer-id02"},
				},
				"coordinator-id02": {
					ID:            "shop-id02",
					CoordinatorID: "coordinator-id02",
					ProducerIDs:   []string{"producer-id01", "producer-id03"},
				},
			},
		},
		{
			name: "duplicated coordinator id overwrites",
			shops: Shops{
				{
					ID:            "shop-id01",
					CoordinatorID: "coordinator-id01",
				},
				{
					ID:            "shop-id02",
					CoordinatorID: "coordinator-id01",
				},
			},
			expect: map[string]*Shop{
				"coordinator-id01": {
					ID:            "shop-id02",
					CoordinatorID: "coordinator-id01",
				},
			},
		},
		{
			name:   "empty",
			shops:  Shops{},
			expect: map[string]*Shop{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.shops.MapByCoordinatorID())
		})
	}
}

func TestShops_GroupByProducerID(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name   string
		shops  Shops
		expect map[string]Shops
	}{
		{
			name: "success",
			shops: Shops{
				{
					ID:            "shop-id01",
					CoordinatorID: "coordinator-id01",
					ProducerIDs:   []string{"producer-id01", "producer-id02"},
				},
				{
					ID:            "shop-id02",
					CoordinatorID: "coordinator-id02",
					ProducerIDs:   []string{"producer-id01", "producer-id03"},
				},
			},
			expect: map[string]Shops{
				"producer-id01": {
					{
						ID:            "shop-id01",
						CoordinatorID: "coordinator-id01",
						ProducerIDs:   []string{"producer-id01", "producer-id02"},
					},
					{
						ID:            "shop-id02",
						CoordinatorID: "coordinator-id02",
						ProducerIDs:   []string{"producer-id01", "producer-id03"},
					},
				},
				"producer-id02": {
					{
						ID:            "shop-id01",
						CoordinatorID: "coordinator-id01",
						ProducerIDs:   []string{"producer-id01", "producer-id02"},
					},
				},
				"producer-id03": {
					{
						ID:            "shop-id02",
						CoordinatorID: "coordinator-id02",
						ProducerIDs:   []string{"producer-id01", "producer-id03"},
					},
				},
			},
		},
		{
			name: "shop with no producers",
			shops: Shops{
				{
					ID:            "shop-id01",
					CoordinatorID: "coordinator-id01",
					ProducerIDs:   []string{},
				},
				{
					ID:            "shop-id02",
					CoordinatorID: "coordinator-id02",
					ProducerIDs:   []string{"producer-id01"},
				},
			},
			expect: map[string]Shops{
				"producer-id01": {
					{
						ID:            "shop-id02",
						CoordinatorID: "coordinator-id02",
						ProducerIDs:   []string{"producer-id01"},
					},
				},
			},
		},
		{
			name:   "empty",
			shops:  Shops{},
			expect: map[string]Shops{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.shops.GroupByProducerID())
		})
	}
}

func TestShops_Fill(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name      string
		shops     Shops
		producers map[string]ShopProducers
		expect    Shops
	}{
		{
			name: "success",
			shops: Shops{
				{ID: "shop-id01"},
				{ID: "shop-id02"},
				{ID: "shop-id03"},
			},
			producers: map[string]ShopProducers{
				"shop-id01": {
					{ProducerID: "producer-id01"},
					{ProducerID: "producer-id02"},
				},
				"shop-id02": {
					{ProducerID: "producer-id03"},
				},
			},
			expect: Shops{
				{ID: "shop-id01", ProducerIDs: []string{"producer-id01", "producer-id02"}},
				{ID: "shop-id02", ProducerIDs: []string{"producer-id03"}},
				{ID: "shop-id03", ProducerIDs: []string{}},
			},
		},
		{
			name: "nil producers map",
			shops: Shops{
				{ID: "shop-id01"},
				{ID: "shop-id02"},
			},
			producers: nil,
			expect: Shops{
				{ID: "shop-id01", ProducerIDs: []string{}},
				{ID: "shop-id02", ProducerIDs: []string{}},
			},
		},
		{
			name:      "empty shops",
			shops:     Shops{},
			producers: map[string]ShopProducers{"shop-id01": {{ProducerID: "producer-id01"}}},
			expect:    Shops{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.shops.Fill(tt.producers)
			assert.Equal(t, tt.expect, tt.shops)
		})
	}
}