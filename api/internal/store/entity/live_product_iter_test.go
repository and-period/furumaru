package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiveProducts_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products LiveProducts
	}{
		{
			name: "success",
			products: LiveProducts{
				{LiveID: "live-01", ProductID: "product-01"},
				{LiveID: "live-01", ProductID: "product-02"},
			},
		},
		{
			name:     "empty",
			products: LiveProducts{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			for i := range tt.products.All() {
				indices = append(indices, i)
			}
			assert.Len(t, indices, len(tt.products))
		})
	}
}

func TestLiveProducts_IterGroupByLiveID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		products   LiveProducts
		expectKeys int
	}{
		{
			name: "success",
			products: LiveProducts{
				{LiveID: "live-01", ProductID: "product-01"},
				{LiveID: "live-01", ProductID: "product-02"},
				{LiveID: "live-02", ProductID: "product-03"},
			},
			expectKeys: 2,
		},
		{
			name:       "empty",
			products:   LiveProducts{},
			expectKeys: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]LiveProducts)
			for k, v := range tt.products.IterGroupByLiveID() {
				result[k] = v
			}
			assert.Len(t, result, tt.expectKeys)
		})
	}
}
