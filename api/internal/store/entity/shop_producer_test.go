package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShopProducers_ProducerIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		producers ShopProducers
		expect    []string
	}{
		{
			name: "success",
			producers: ShopProducers{
				{ProducerID: "1"},
				{ProducerID: "2"},
				{ProducerID: "3"},
			},
			expect: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.producers.ProducerIDs()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestShopProducers_GroupByShopID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		producers ShopProducers
		expect    map[string]ShopProducers
	}{
		{
			name: "success",
			producers: ShopProducers{
				{ShopID: "1", ProducerID: "1"},
				{ShopID: "1", ProducerID: "2"},
				{ShopID: "2", ProducerID: "3"},
			},
			expect: map[string]ShopProducers{
				"1": {
					{ShopID: "1", ProducerID: "1"},
					{ShopID: "1", ProducerID: "2"},
				},
				"2": {
					{ShopID: "2", ProducerID: "3"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.producers.GroupByShopID()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
