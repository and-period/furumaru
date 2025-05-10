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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductRevision(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProductRevisions_ProductIDs(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name      string
		revisions ProductRevisions
		expect    []string
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
			expect: []string{"product-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.revisions.ProductIDs())
		})
	}
}

func TestProductRevisions_MapByProductID(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.revisions.MapByProductID())
		})
	}
}

func TestProductRevisions_Merge(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name      string
		revisions ProductRevisions
		products  map[string]*Product
		expect    Products
		hasErr    bool
	}{
		{
			name: "success",
			revisions: ProductRevisions{
				{
					ID:        1,
					ProductID: "product-id01",
					Price:     3980,
					Cost:      880,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        2,
					ProductID: "product-id02",
					Price:     1500,
					Cost:      200,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        3,
					ProductID: "product-id01",
					Price:     2000,
					Cost:      880,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			products: map[string]*Product{
				"product-id01": {
					ID:   "product-id01",
					Name: "芽が出たじゃがいも",
				},
			},
			expect: Products{
				{
					ID:   "product-id01",
					Name: "芽が出たじゃがいも",
					ProductRevision: ProductRevision{
						ID:        1,
						ProductID: "product-id01",
						Price:     3980,
						Cost:      880,
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				{
					ID: "product-id02",
					ProductRevision: ProductRevision{
						ID:        2,
						ProductID: "product-id02",
						Price:     1500,
						Cost:      200,
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				{
					ID:   "product-id01",
					Name: "芽が出たじゃがいも",
					ProductRevision: ProductRevision{
						ID:        3,
						ProductID: "product-id01",
						Price:     2000,
						Cost:      880,
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.revisions.Merge(tt.products)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
