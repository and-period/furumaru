package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideoProducts_ProductIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		products VideoProducts
		expect   []string
	}{
		{
			name: "success",
			products: VideoProducts{
				{VideoID: "video-id01", ProductID: "product-id01"},
				{VideoID: "video-id02", ProductID: "product-id02"},
				{VideoID: "video-id02", ProductID: "product-id03"},
			},
			expect: []string{"product-id01", "product-id02", "product-id03"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.products.ProductIDs()
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestVideoProducts_GroupByVideoID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		products VideoProducts
		expect   map[string]VideoProducts
	}{
		{
			name: "success",
			products: VideoProducts{
				{VideoID: "video-id01", ProductID: "product-id01"},
				{VideoID: "video-id02", ProductID: "product-id02"},
				{VideoID: "video-id02", ProductID: "product-id03"},
			},
			expect: map[string]VideoProducts{
				"video-id01": {
					{VideoID: "video-id01", ProductID: "product-id01"},
				},
				"video-id02": {
					{VideoID: "video-id02", ProductID: "product-id02"},
					{VideoID: "video-id02", ProductID: "product-id03"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.products.GroupByVideoID()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideoProducts_SortByPrority(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		products VideoProducts
		expect   VideoProducts
	}{
		{
			name: "success",
			products: VideoProducts{
				{VideoID: "video-id01", Priority: 1, ProductID: "product-id01"},
				{VideoID: "video-id02", Priority: 3, ProductID: "product-id03"},
				{VideoID: "video-id02", Priority: 2, ProductID: "product-id02"},
			},
			expect: VideoProducts{
				{VideoID: "video-id01", Priority: 1, ProductID: "product-id01"},
				{VideoID: "video-id02", Priority: 2, ProductID: "product-id02"},
				{VideoID: "video-id02", Priority: 3, ProductID: "product-id03"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.products.SortByPriority()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
