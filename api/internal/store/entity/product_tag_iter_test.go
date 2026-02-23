package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductTags_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		tags ProductTags
	}{
		{
			name: "success",
			tags: ProductTags{
				{ID: "tag-01", Name: "有機"},
				{ID: "tag-02", Name: "無農薬"},
			},
		},
		{
			name: "empty",
			tags: ProductTags{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, tag := range tt.tags.All() {
				indices = append(indices, i)
				ids = append(ids, tag.ID)
			}
			for i, tag := range tt.tags {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, tag.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.tags))
		})
	}
}

func TestProductTags_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		tags ProductTags
	}{
		{
			name: "success",
			tags: ProductTags{
				{ID: "tag-01", Name: "有機"},
				{ID: "tag-02", Name: "無農薬"},
			},
		},
		{
			name: "empty",
			tags: ProductTags{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*ProductTag)
			for k, v := range tt.tags.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.tags))
			for _, tag := range tt.tags {
				assert.Contains(t, result, tag.ID)
			}
		})
	}
}
