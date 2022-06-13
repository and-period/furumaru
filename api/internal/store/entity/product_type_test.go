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
		categoryID string
		expect     *ProductType
	}{
		{
			name:       "success",
			input:      "じゃがいも",
			categoryID: "category-id",
			expect: &ProductType{
				Name:       "じゃがいも",
				CategoryID: "category-id",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductType(tt.input, tt.categoryID)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}
