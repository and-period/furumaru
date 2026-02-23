package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategories_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		categories Categories
	}{
		{
			name: "success",
			categories: Categories{
				{ID: "cat-01", Name: "野菜"},
				{ID: "cat-02", Name: "果物"},
			},
		},
		{
			name:       "empty",
			categories: Categories{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var categories Categories
			for i, c := range tt.categories.All() {
				indices = append(indices, i)
				categories = append(categories, c)
			}
			for i, c := range tt.categories {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, c, categories[i])
				}
			}
			assert.Len(t, indices, len(tt.categories))
		})
	}
}

func TestCategories_IterMapByName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		categories Categories
		expect     map[string]*Category
	}{
		{
			name: "success",
			categories: Categories{
				{ID: "cat-01", Name: "野菜"},
				{ID: "cat-02", Name: "果物"},
			},
			expect: map[string]*Category{
				"野菜": {ID: "cat-01", Name: "野菜"},
				"果物": {ID: "cat-02", Name: "果物"},
			},
		},
		{
			name:       "empty",
			categories: Categories{},
			expect:     map[string]*Category{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Category)
			for k, v := range tt.categories.IterMapByName() {
				result[k] = v
			}
			assert.Equal(t, tt.expect, result)
		})
	}
}
