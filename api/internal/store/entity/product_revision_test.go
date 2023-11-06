package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProductRevision(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewProductRevisionParams
		expect *ProductRevision
	}{
		{
			name: "success",
			params: &NewProductRevisionParams{
				ProductID: "product-id",
				Price:     3980,
				Cost:      880,
			},
			expect: &ProductRevision{
				ProductID: "product-id",
				Price:     3980,
				Cost:      880,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductRevision(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProductRevisions_Map(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name      string
		revisions ProductRevisions
		expect    map[string]*ProductRevision
	}{
		{
			name: "success",
			revisions: ProductRevisions{
				{
					ID:        1,
					ProductID: "product-id",
					Price:     3980,
					Cost:      880,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			expect: map[string]*ProductRevision{
				"product-id": {
					ID:        1,
					ProductID: "product-id",
					Price:     3980,
					Cost:      880,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.revisions.Map())
		})
	}
}
