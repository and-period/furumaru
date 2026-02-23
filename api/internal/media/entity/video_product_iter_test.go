package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideoProducts_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products VideoProducts
	}{
		{
			name: "success",
			products: VideoProducts{
				{VideoID: "video-id01", ProductID: "product-id01", Priority: 1},
				{VideoID: "video-id01", ProductID: "product-id02", Priority: 2},
			},
		},
		{
			name:     "empty",
			products: VideoProducts{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var productIDs []string
			for i, p := range tt.products.All() {
				indices = append(indices, i)
				productIDs = append(productIDs, p.ProductID)
			}
			for i, p := range tt.products {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, p.ProductID, productIDs[i])
				}
			}
			assert.Len(t, indices, len(tt.products))
		})
	}
}

func TestVideoProducts_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	products := VideoProducts{
		{VideoID: "video-id01", ProductID: "product-id01"},
		{VideoID: "video-id01", ProductID: "product-id02"},
		{VideoID: "video-id01", ProductID: "product-id03"},
	}
	var count int
	for range products.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestVideoProducts_IterProductIDs(t *testing.T) {
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
				{VideoID: "video-id01", ProductID: "product-id02"},
			},
			expect: []string{"product-id01", "product-id02"},
		},
		{
			name:     "empty",
			products: VideoProducts{},
			expect:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var actual []string
			for id := range tt.products.IterProductIDs() {
				actual = append(actual, id)
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideoProducts_IterGroupByVideoID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products VideoProducts
		expect   map[string]int // videoID -> count
	}{
		{
			name: "success",
			products: VideoProducts{
				{VideoID: "video-id01", ProductID: "product-id01"},
				{VideoID: "video-id01", ProductID: "product-id02"},
				{VideoID: "video-id02", ProductID: "product-id03"},
			},
			expect: map[string]int{
				"video-id01": 2,
				"video-id02": 1,
			},
		},
		{
			name:     "empty",
			products: VideoProducts{},
			expect:   map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]int)
			for k, v := range tt.products.IterGroupByVideoID() {
				result[k] = len(v)
			}
			assert.Equal(t, tt.expect, result)
		})
	}
}
