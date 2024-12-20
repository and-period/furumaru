package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductReview(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewProductReviewParams
		expect *ProductReview
	}{
		{
			name: "success",
			params: &NewProductReviewParams{
				ProductID: "product-id",
				UserID:    "user-id",
				Rate:      5,
				Title:     "最高の商品",
				Comment:   "おすすめできる商品です",
			},
			expect: &ProductReview{
				ProductID: "product-id",
				UserID:    "user-id",
				Rate:      5,
				Title:     "最高の商品",
				Comment:   "おすすめできる商品です",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewProductReview(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}
