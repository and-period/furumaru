package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLiveProduct(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		liveID    string
		productID string
		expect    *LiveProduct
	}{
		{
			name:      "success",
			liveID:    "live-id",
			productID: "product-id",
			expect: &LiveProduct{
				LiveID:    "live-id",
				ProductID: "product-id",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewLiveProduct(tt.liveID, tt.productID))
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
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
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
					ProductID: "product-id",
				},
			},
			expect: []string{"product-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.products.ProductIDs())
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.products.GroupByLiveID())
		})
	}
}

func TestLiveProducts_SortByCreatedAt(t *testing.T) {
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
				{ProductID: "product-id01", CreatedAt: now},
				{ProductID: "product-id02", CreatedAt: now.Add(-time.Minute)},
				{ProductID: "product-id03", CreatedAt: now.Add(time.Minute)},
				{ProductID: "product-id04", CreatedAt: now.Add(time.Hour)},
				{ProductID: "product-id05", CreatedAt: now.Add(-time.Hour)},
			},
			expect: LiveProducts{
				{ProductID: "product-id05", CreatedAt: now.Add(-time.Hour)},
				{ProductID: "product-id02", CreatedAt: now.Add(-time.Minute)},
				{ProductID: "product-id01", CreatedAt: now},
				{ProductID: "product-id03", CreatedAt: now.Add(time.Minute)},
				{ProductID: "product-id04", CreatedAt: now.Add(time.Hour)},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.products.SortByCreatedAt()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
