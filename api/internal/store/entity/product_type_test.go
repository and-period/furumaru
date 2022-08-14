package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		iconURL    string
		categoryID string
		expect     *ProductType
	}{
		{
			name:       "success",
			input:      "じゃがいも",
			iconURL:    "https://and-period.jp/icon.png",
			categoryID: "category-id",
			expect: &ProductType{
				Name:       "じゃがいも",
				IconURL:    "https://and-period.jp/icon.png",
				CategoryID: "category-id",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductType(tt.input, tt.iconURL, tt.categoryID)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}
