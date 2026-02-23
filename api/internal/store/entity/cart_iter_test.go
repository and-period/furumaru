package entity

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCartBaskets_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		baskets CartBaskets
	}{
		{
			name: "success",
			baskets: CartBaskets{
				{BoxNumber: 1, CoordinatorID: "coord-01"},
				{BoxNumber: 2, CoordinatorID: "coord-02"},
			},
		},
		{
			name:    "empty",
			baskets: CartBaskets{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var count int
			for range tt.baskets.All() {
				count++
			}
			assert.Equal(t, len(tt.baskets), count)
		})
	}
}

func TestCartBaskets_IterFilterByCoordinatorID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		baskets        CartBaskets
		coordinatorIDs []string
		expectCount    int
	}{
		{
			name: "success",
			baskets: CartBaskets{
				{BoxNumber: 1, CoordinatorID: "coord-01"},
				{BoxNumber: 2, CoordinatorID: "coord-02"},
				{BoxNumber: 3, CoordinatorID: "coord-01"},
			},
			coordinatorIDs: []string{"coord-01"},
			expectCount:    2,
		},
		{
			name: "no match",
			baskets: CartBaskets{
				{BoxNumber: 1, CoordinatorID: "coord-01"},
			},
			coordinatorIDs: []string{"coord-99"},
			expectCount:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			collected := slices.Collect(tt.baskets.IterFilterByCoordinatorID(tt.coordinatorIDs...))
			assert.Len(t, collected, tt.expectCount)
		})
	}
}

func TestCartBaskets_IterFilterByBoxNumber(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		baskets     CartBaskets
		targets     []int64
		expectCount int
	}{
		{
			name: "success",
			baskets: CartBaskets{
				{BoxNumber: 1, CoordinatorID: "coord-01"},
				{BoxNumber: 2, CoordinatorID: "coord-02"},
				{BoxNumber: 3, CoordinatorID: "coord-01"},
			},
			targets:     []int64{1, 3},
			expectCount: 2,
		},
		{
			name: "zero includes all",
			baskets: CartBaskets{
				{BoxNumber: 1, CoordinatorID: "coord-01"},
				{BoxNumber: 2, CoordinatorID: "coord-02"},
			},
			targets:     []int64{0},
			expectCount: 2,
		},
		{
			name: "no match",
			baskets: CartBaskets{
				{BoxNumber: 1, CoordinatorID: "coord-01"},
			},
			targets:     []int64{99},
			expectCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			collected := slices.Collect(tt.baskets.IterFilterByBoxNumber(tt.targets...))
			assert.Len(t, collected, tt.expectCount)
		})
	}
}

func TestCartItems_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		items CartItems
	}{
		{
			name: "success",
			items: CartItems{
				{ProductID: "product-01", Quantity: 2},
				{ProductID: "product-02", Quantity: 1},
			},
		},
		{
			name:  "empty",
			items: CartItems{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var count int
			for range tt.items.All() {
				count++
			}
			assert.Equal(t, len(tt.items), count)
		})
	}
}

func TestCartItems_IterMapByProductID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		items CartItems
	}{
		{
			name: "success",
			items: CartItems{
				{ProductID: "product-01", Quantity: 2},
				{ProductID: "product-02", Quantity: 1},
			},
		},
		{
			name:  "empty",
			items: CartItems{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*CartItem)
			for k, v := range tt.items.IterMapByProductID() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.items))
			for _, item := range tt.items {
				assert.Contains(t, result, item.ProductID)
			}
		})
	}
}
