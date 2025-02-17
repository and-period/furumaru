package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestShop(t *testing.T) {
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
				Name:           "test",
				ProductTypeIDs: []string{"product-type-id"},
				BusinessDays:   []time.Weekday{time.Monday, time.Tuesday},
			},
			expect: &Shop{
				CoordinatorID:  "coordinator-id",
				Name:           "test",
				ProductTypeIDs: []string{"product-type-id"},
				BusinessDays:   []time.Weekday{time.Monday, time.Tuesday},
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
			name: "enabled",
			shop: &Shop{
				Activated: true,
			},
			expect: true,
		},
		{
			name: "disabled",
			shop: &Shop{
				Activated: false,
			},
			expect: false,
		},
		{
			name: "deleted",
			shop: &Shop{
				Activated: true,
				DeletedAt: gorm.DeletedAt{Time: time.Now()},
			},
			expect: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.shop.Enabled()
			assert.Equal(t, tt.expect, actual)
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
				{ID: "1"},
				{ID: "2"},
				{ID: "3"},
			},
			expect: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.shops.IDs()
			assert.ElementsMatch(t, tt.expect, actual)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.shops.MapByCoordinatorID()
			assert.Equal(t, tt.expect, actual)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.shops.GroupByProducerID()
			assert.Equal(t, tt.expect, actual)
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
				{ID: "1"},
				{ID: "2"},
				{ID: "3"},
			},
			producers: map[string]ShopProducers{
				"1": {
					{ProducerID: "1"},
					{ProducerID: "2"},
				},
				"2": {
					{ProducerID: "3"},
				},
			},
			expect: Shops{
				{ID: "1", ProducerIDs: []string{"1", "2"}},
				{ID: "2", ProducerIDs: []string{"3"}},
				{ID: "3", ProducerIDs: []string{}},
			},
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
