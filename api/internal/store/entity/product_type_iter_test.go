package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductTypes_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		types ProductTypes
	}{
		{
			name: "success",
			types: ProductTypes{
				{ID: "type-01", CategoryID: "cat-01", Name: "じゃがいも"},
				{ID: "type-02", CategoryID: "cat-01", Name: "にんじん"},
			},
		},
		{
			name:  "empty",
			types: ProductTypes{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, pt := range tt.types.All() {
				indices = append(indices, i)
				ids = append(ids, pt.ID)
			}
			for i, pt := range tt.types {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, pt.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.types))
		})
	}
}

func TestProductTypes_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		types ProductTypes
	}{
		{
			name: "success",
			types: ProductTypes{
				{ID: "type-01", Name: "じゃがいも"},
				{ID: "type-02", Name: "にんじん"},
			},
		},
		{
			name:  "empty",
			types: ProductTypes{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*ProductType)
			for k, v := range tt.types.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.types))
			for _, pt := range tt.types {
				assert.Contains(t, result, pt.ID)
			}
		})
	}
}
