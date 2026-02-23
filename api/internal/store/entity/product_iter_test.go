package entity

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProducts_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products Products
	}{
		{
			name: "success",
			products: Products{
				{ID: "product-id01", Name: "じゃがいも"},
				{ID: "product-id02", Name: "にんじん"},
				{ID: "product-id03", Name: "たまねぎ"},
			},
		},
		{
			name:     "empty",
			products: Products{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, p := range tt.products.All() {
				indices = append(indices, i)
				ids = append(ids, p.ID)
			}
			for i, p := range tt.products {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, p.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.products))
		})
	}
}

func TestProducts_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	products := Products{
		{ID: "product-id01"},
		{ID: "product-id02"},
		{ID: "product-id03"},
	}
	var count int
	for range products.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func collectProductIDs(seq func(func(*Product) bool)) []string {
	products := slices.Collect(seq)
	if len(products) == 0 {
		return nil
	}
	ids := make([]string, len(products))
	for i, p := range products {
		ids[i] = p.ID
	}
	return ids
}

func TestProducts_IterFilter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		products   Products
		productIDs []string
		expectIDs  []string
	}{
		{
			name: "success",
			products: Products{
				{ID: "product-id01"},
				{ID: "product-id02"},
				{ID: "product-id03"},
			},
			productIDs: []string{
				"product-id01",
				"product-id03",
			},
			expectIDs: []string{
				"product-id01",
				"product-id03",
			},
		},
		{
			name: "empty result",
			products: Products{
				{ID: "product-id01"},
				{ID: "product-id02"},
				{ID: "product-id03"},
			},
			productIDs: []string{},
			expectIDs:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := collectProductIDs(tt.products.IterFilter(tt.productIDs...))
			assert.Equal(t, tt.expectIDs, actual)
		})
	}
}

func TestProducts_IterFilterByProducerID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		products   Products
		producerID string
		expectIDs  []string
	}{
		{
			name: "success",
			products: Products{
				{ID: "product-id01", ProducerID: "producer-01"},
				{ID: "product-id02", ProducerID: "producer-02"},
				{ID: "product-id03", ProducerID: "producer-01"},
			},
			producerID: "producer-01",
			expectIDs:  []string{"product-id01", "product-id03"},
		},
		{
			name: "no match",
			products: Products{
				{ID: "product-id01", ProducerID: "producer-01"},
			},
			producerID: "producer-99",
			expectIDs:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := collectProductIDs(tt.products.IterFilterByProducerID(tt.producerID))
			assert.Equal(t, tt.expectIDs, actual)
		})
	}
}

func TestProducts_IterFilterBySales(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		products  Products
		expectIDs []string
	}{
		{
			name: "success",
			products: Products{
				{ID: "product-id01", Status: ProductStatusForSale},
				{ID: "product-id02", Status: ProductStatusPrivate},
				{ID: "product-id03", Status: ProductStatusForSale},
			},
			expectIDs: []string{"product-id01", "product-id03"},
		},
		{
			name: "no sales",
			products: Products{
				{ID: "product-id01", Status: ProductStatusPrivate},
			},
			expectIDs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := collectProductIDs(tt.products.IterFilterBySales())
			assert.Equal(t, tt.expectIDs, actual)
		})
	}
}

func TestProducts_IterFilterByPublished(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		products  Products
		expectIDs []string
	}{
		{
			name: "success",
			products: Products{
				{ID: "product-id01", Status: ProductStatusForSale},
				{ID: "product-id02", Status: ProductStatusPrivate},
				{ID: "product-id03", Status: ProductStatusPresale},
				{ID: "product-id04", Status: ProductStatusArchived},
			},
			expectIDs: []string{"product-id01", "product-id03"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := collectProductIDs(tt.products.IterFilterByPublished())
			assert.Equal(t, tt.expectIDs, actual)
		})
	}
}

func TestProducts_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		products  Products
		expectIDs []string
	}{
		{
			name: "success",
			products: Products{
				{ID: "product-id01", Name: "じゃがいも"},
				{ID: "product-id02", Name: "にんじん"},
			},
			expectIDs: []string{"product-id01", "product-id02"},
		},
		{
			name:      "empty",
			products:  Products{},
			expectIDs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Product)
			for k, v := range tt.products.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.products))
			for _, id := range tt.expectIDs {
				assert.Contains(t, result, id)
				assert.Equal(t, id, result[id].ID)
			}
		})
	}
}

func TestProducts_IterMapByRevision(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products Products
		expect   map[int64]string // revision ID -> product ID
	}{
		{
			name: "success",
			products: Products{
				{
					ID:              "product-id01",
					ProductRevision: ProductRevision{ID: 1, ProductID: "product-id01"},
				},
				{
					ID:              "product-id02",
					ProductRevision: ProductRevision{ID: 2, ProductID: "product-id02"},
				},
			},
			expect: map[int64]string{
				1: "product-id01",
				2: "product-id02",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[int64]*Product)
			for k, v := range tt.products.IterMapByRevision() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.expect))
			for revID, productID := range tt.expect {
				assert.Contains(t, result, revID)
				assert.Equal(t, productID, result[revID].ID)
			}
		})
	}
}
