package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpots_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		spots Spots
	}{
		{
			name: "success",
			spots: Spots{
				{ID: "spot-01", Name: "農園A", UserType: SpotUserTypeUser},
				{ID: "spot-02", Name: "農園B", UserType: SpotUserTypeCoordinator},
			},
		},
		{
			name:  "empty",
			spots: Spots{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, s := range tt.spots.All() {
				indices = append(indices, i)
				ids = append(ids, s.ID)
			}
			for i, s := range tt.spots {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, s.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.spots))
		})
	}
}

func TestSpots_IterGroupByUserType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		spots      Spots
		expectKeys int
	}{
		{
			name: "success",
			spots: Spots{
				{ID: "spot-01", UserType: SpotUserTypeUser},
				{ID: "spot-02", UserType: SpotUserTypeUser},
				{ID: "spot-03", UserType: SpotUserTypeCoordinator},
			},
			expectKeys: 2,
		},
		{
			name:       "empty",
			spots:      Spots{},
			expectKeys: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[SpotUserType]Spots)
			for k, v := range tt.spots.IterGroupByUserType() {
				result[k] = v
			}
			assert.Len(t, result, tt.expectKeys)
		})
	}
}
