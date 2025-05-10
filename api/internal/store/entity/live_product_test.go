package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLiveProduct(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewLiveProductParams
		expect *LiveProduct
	}{
		{
			name: "success",
			params: &NewLiveProductParams{
				LiveID:    "live-id",
				ProductID: "product-id",
				Priority:  1,
			},
			expect: &LiveProduct{
				LiveID:    "live-id",
				ProductID: "product-id",
				Priority:  1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewLiveProduct(tt.params))
		})
	}
}

func TestLiveProducts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		liveID     string
		productIDs []string
		expect     LiveProducts
	}{
		{
			name:       "success",
			liveID:     "live-id",
			productIDs: []string{"product-id"},
			expect: LiveProducts{
				{
					LiveID:    "live-id",
					ProductID: "product-id",
					Priority:  1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewLiveProducts(tt.liveID, tt.productIDs))
		})
	}
}

func TestLiveProducts_ProductIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products LiveProducts
		expect   []string
	}{
		{
			name: "success",
			products: LiveProducts{
				{
					LiveID:    "live-id",
					ProductID: "product-id-0",
					Priority:  2,
				},
				{
					LiveID:    "live-id",
					ProductID: "product-id-1",
					Priority:  1,
				},
				{
					LiveID:    "live-id",
					ProductID: "product-id-2",
					Priority:  0,
				},
			},
			expect: []string{"product-id-2", "product-id-1", "product-id-0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Exactly(t, tt.expect, tt.products.ProductIDs())
		})
	}
}

func TestLiveProducts_GroupByLiveID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products LiveProducts
		expect   map[string]LiveProducts
	}{
		{
			name: "success",
			products: LiveProducts{
				{
					LiveID:    "live-id",
					ProductID: "product-id",
				},
			},
			expect: map[string]LiveProducts{
				"live-id": {
					{
						LiveID:    "live-id",
						ProductID: "product-id",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.products.GroupByLiveID())
		})
	}
}

func TestLiveProducts_SortByPrimary(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		products LiveProducts
		expect   LiveProducts
	}{
		{
			name: "success",
			products: LiveProducts{
				{ProductID: "product-id01", Priority: 3, CreatedAt: now},
				{ProductID: "product-id02", Priority: 2, CreatedAt: now.Add(-time.Minute)},
				{ProductID: "product-id03", Priority: 4, CreatedAt: now.Add(time.Minute)},
				{ProductID: "product-id04", Priority: 5, CreatedAt: now.Add(time.Hour)},
				{ProductID: "product-id05", Priority: 1, CreatedAt: now.Add(-time.Hour)},
			},
			expect: LiveProducts{
				{ProductID: "product-id05", Priority: 1, CreatedAt: now.Add(-time.Hour)},
				{ProductID: "product-id02", Priority: 2, CreatedAt: now.Add(-time.Minute)},
				{ProductID: "product-id01", Priority: 3, CreatedAt: now},
				{ProductID: "product-id03", Priority: 4, CreatedAt: now.Add(time.Minute)},
				{ProductID: "product-id04", Priority: 5, CreatedAt: now.Add(time.Hour)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.products.SortByPrimary()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
