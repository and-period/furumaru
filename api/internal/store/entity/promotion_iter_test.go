package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPromotions_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		promotions Promotions
	}{
		{
			name: "success",
			promotions: Promotions{
				{ID: "promo-01", Title: "夏セール"},
				{ID: "promo-02", Title: "冬セール"},
			},
		},
		{
			name:       "empty",
			promotions: Promotions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, p := range tt.promotions.All() {
				indices = append(indices, i)
				ids = append(ids, p.ID)
			}
			for i, p := range tt.promotions {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, p.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.promotions))
		})
	}
}

func TestPromotions_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	promotions := Promotions{
		{ID: "promo-01"},
		{ID: "promo-02"},
		{ID: "promo-03"},
	}
	var count int
	for range promotions.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestPromotions_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		promotions Promotions
	}{
		{
			name: "success",
			promotions: Promotions{
				{ID: "promo-01", Title: "夏セール"},
				{ID: "promo-02", Title: "冬セール"},
			},
		},
		{
			name:       "empty",
			promotions: Promotions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Promotion)
			for k, v := range tt.promotions.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.promotions))
			for _, p := range tt.promotions {
				assert.Contains(t, result, p.ID)
			}
		})
	}
}
