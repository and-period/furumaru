package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShippings_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		shippings Shippings
	}{
		{
			name: "success",
			shippings: Shippings{
				{ID: "shipping-01", Name: "通常配送"},
				{ID: "shipping-02", Name: "冷凍配送"},
			},
		},
		{
			name:      "empty",
			shippings: Shippings{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, s := range tt.shippings.All() {
				indices = append(indices, i)
				ids = append(ids, s.ID)
			}
			for i, s := range tt.shippings {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, s.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.shippings))
		})
	}
}

func TestShippings_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	shippings := Shippings{
		{ID: "shipping-01"},
		{ID: "shipping-02"},
		{ID: "shipping-03"},
	}
	var count int
	for range shippings.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestShippings_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		shippings Shippings
	}{
		{
			name: "success",
			shippings: Shippings{
				{ID: "shipping-01", Name: "通常配送"},
				{ID: "shipping-02", Name: "冷凍配送"},
			},
		},
		{
			name:      "empty",
			shippings: Shippings{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Shipping)
			for k, v := range tt.shippings.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.shippings))
			for _, s := range tt.shippings {
				assert.Contains(t, result, s.ID)
			}
		})
	}
}
