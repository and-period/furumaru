package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
