package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductRevisions_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		revisions ProductRevisions
	}{
		{
			name: "success",
			revisions: ProductRevisions{
				{ID: 1, ProductID: "product-01"},
				{ID: 2, ProductID: "product-02"},
			},
		},
		{
			name:      "empty",
			revisions: ProductRevisions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var count int
			for range tt.revisions.All() {
				count++
			}
			assert.Equal(t, len(tt.revisions), count)
		})
	}
}

func TestProductRevisions_IterMapByProductID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		revisions ProductRevisions
	}{
		{
			name: "success",
			revisions: ProductRevisions{
				{ID: 1, ProductID: "product-01"},
				{ID: 2, ProductID: "product-02"},
			},
		},
		{
			name:      "empty",
			revisions: ProductRevisions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*ProductRevision)
			for k, v := range tt.revisions.IterMapByProductID() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.revisions))
			for _, r := range tt.revisions {
				assert.Contains(t, result, r.ProductID)
			}
		})
	}
}
