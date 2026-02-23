package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductReviews_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews ProductReviews
	}{
		{
			name: "success",
			reviews: ProductReviews{
				{ID: "review-01", ProductID: "product-01"},
				{ID: "review-02", ProductID: "product-02"},
			},
		},
		{
			name:    "empty",
			reviews: ProductReviews{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, r := range tt.reviews.All() {
				indices = append(indices, i)
				ids = append(ids, r.ID)
			}
			for i, r := range tt.reviews {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, r.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.reviews))
		})
	}
}

func TestProductReviews_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews ProductReviews
	}{
		{
			name: "success",
			reviews: ProductReviews{
				{ID: "review-01", ProductID: "product-01"},
				{ID: "review-02", ProductID: "product-02"},
			},
		},
		{
			name:    "empty",
			reviews: ProductReviews{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*ProductReview)
			for k, v := range tt.reviews.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.reviews))
			for _, r := range tt.reviews {
				assert.Contains(t, result, r.ID)
			}
		})
	}
}

func TestAggregatedProductReviews_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews AggregatedProductReviews
	}{
		{
			name: "success",
			reviews: AggregatedProductReviews{
				{ProductID: "product-01", Count: 10, Average: 4.5},
				{ProductID: "product-02", Count: 5, Average: 3.8},
			},
		},
		{
			name:    "empty",
			reviews: AggregatedProductReviews{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var count int
			for range tt.reviews.All() {
				count++
			}
			assert.Equal(t, len(tt.reviews), count)
		})
	}
}

func TestAggregatedProductReviews_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews AggregatedProductReviews
	}{
		{
			name: "success",
			reviews: AggregatedProductReviews{
				{ProductID: "product-01", Count: 10, Average: 4.5},
				{ProductID: "product-02", Count: 5, Average: 3.8},
			},
		},
		{
			name:    "empty",
			reviews: AggregatedProductReviews{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*AggregatedProductReview)
			for k, v := range tt.reviews.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.reviews))
			for _, r := range tt.reviews {
				assert.Contains(t, result, r.ProductID)
			}
		})
	}
}
