package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShops_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		shops Shops
	}{
		{
			name: "success",
			shops: Shops{
				{ID: "shop-id01", CoordinatorID: "coordinator-id01"},
				{ID: "shop-id02", CoordinatorID: "coordinator-id02"},
			},
		},
		{
			name:  "empty",
			shops: Shops{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, s := range tt.shops.All() {
				indices = append(indices, i)
				ids = append(ids, s.ID)
			}
			for i, s := range tt.shops {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, s.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.shops))
		})
	}
}

func TestShops_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	shops := Shops{
		{ID: "shop-id01"},
		{ID: "shop-id02"},
		{ID: "shop-id03"},
	}
	var count int
	for range shops.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestShops_IterMapByCoordinatorID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		shops  Shops
		expect map[string]string
	}{
		{
			name: "success",
			shops: Shops{
				{ID: "shop-id01", CoordinatorID: "coordinator-id01"},
				{ID: "shop-id02", CoordinatorID: "coordinator-id02"},
			},
			expect: map[string]string{
				"coordinator-id01": "shop-id01",
				"coordinator-id02": "shop-id02",
			},
		},
		{
			name:   "empty",
			shops:  Shops{},
			expect: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Shop)
			for k, v := range tt.shops.IterMapByCoordinatorID() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.expect))
			for coordID, shopID := range tt.expect {
				assert.Contains(t, result, coordID)
				assert.Equal(t, shopID, result[coordID].ID)
			}
		})
	}
}

func TestShops_IterGroupByProducerID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		shops  Shops
		expect map[string][]string
	}{
		{
			name: "success",
			shops: Shops{
				{ID: "shop-id01", ProducerIDs: []string{"producer-id01", "producer-id02"}},
				{ID: "shop-id02", ProducerIDs: []string{"producer-id02", "producer-id03"}},
			},
			expect: map[string][]string{
				"producer-id01": {"shop-id01"},
				"producer-id02": {"shop-id01", "shop-id02"},
				"producer-id03": {"shop-id02"},
			},
		},
		{
			name:   "empty",
			shops:  Shops{},
			expect: map[string][]string{},
		},
		{
			name: "no producer IDs",
			shops: Shops{
				{ID: "shop-id01", ProducerIDs: []string{}},
			},
			expect: map[string][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string][]string)
			for k, v := range tt.shops.IterGroupByProducerID() {
				result[k] = append(result[k], v.ID)
			}
			assert.Equal(t, tt.expect, result)
		})
	}
}
