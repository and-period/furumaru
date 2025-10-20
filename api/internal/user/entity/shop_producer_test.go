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
				{ShopID: "shop-id01", ProducerID: "producer-id01"},
				{ShopID: "shop-id01", ProducerID: "producer-id02"},
				{ShopID: "shop-id02", ProducerID: "producer-id03"},
			},
			expect: []string{"producer-id01", "producer-id02", "producer-id03"},
		},
		{
			name: "empty",
			producers: ShopProducers{},
			expect:    []string{},
		},
		{
			name:      "nil",
			producers: nil,
			expect:    []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.producers.ProducerIDs())
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
				{ShopID: "shop-id01", ProducerID: "producer-id01"},
				{ShopID: "shop-id01", ProducerID: "producer-id02"},
				{ShopID: "shop-id02", ProducerID: "producer-id03"},
				{ShopID: "shop-id02", ProducerID: "producer-id04"},
				{ShopID: "shop-id03", ProducerID: "producer-id05"},
			},
			expect: map[string]ShopProducers{
				"shop-id01": {
					{ShopID: "shop-id01", ProducerID: "producer-id01"},
					{ShopID: "shop-id01", ProducerID: "producer-id02"},
				},
				"shop-id02": {
					{ShopID: "shop-id02", ProducerID: "producer-id03"},
					{ShopID: "shop-id02", ProducerID: "producer-id04"},
				},
				"shop-id03": {
					{ShopID: "shop-id03", ProducerID: "producer-id05"},
				},
			},
		},
		{
			name: "single shop multiple producers",
			producers: ShopProducers{
				{ShopID: "shop-id01", ProducerID: "producer-id01"},
				{ShopID: "shop-id01", ProducerID: "producer-id02"},
				{ShopID: "shop-id01", ProducerID: "producer-id03"},
			},
			expect: map[string]ShopProducers{
				"shop-id01": {
					{ShopID: "shop-id01", ProducerID: "producer-id01"},
					{ShopID: "shop-id01", ProducerID: "producer-id02"},
					{ShopID: "shop-id01", ProducerID: "producer-id03"},
				},
			},
		},
		{
			name:      "empty",
			producers: ShopProducers{},
			expect:    map[string]ShopProducers{},
		},
		{
			name:      "nil",
			producers: nil,
			expect:    map[string]ShopProducers{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.producers.GroupByShopID())
		})
	}
}
