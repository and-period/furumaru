package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinators_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		coordinators Coordinators
	}{
		{
			name: "success",
			coordinators: Coordinators{
				{AdminID: "admin-id01", Username: "コーディネータA"},
				{AdminID: "admin-id02", Username: "コーディネータB"},
			},
		},
		{
			name:         "empty",
			coordinators: Coordinators{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, c := range tt.coordinators.All() {
				indices = append(indices, i)
				ids = append(ids, c.AdminID)
			}
			for i, c := range tt.coordinators {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, c.AdminID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.coordinators))
		})
	}
}

func TestCoordinators_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	coordinators := Coordinators{
		{AdminID: "admin-id01"},
		{AdminID: "admin-id02"},
		{AdminID: "admin-id03"},
	}
	var count int
	for range coordinators.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestCoordinators_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		coordinators Coordinators
		expectIDs    []string
	}{
		{
			name: "success",
			coordinators: Coordinators{
				{AdminID: "admin-id01"},
				{AdminID: "admin-id02"},
			},
			expectIDs: []string{"admin-id01", "admin-id02"},
		},
		{
			name:         "empty",
			coordinators: Coordinators{},
			expectIDs:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Coordinator)
			for k, v := range tt.coordinators.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.coordinators))
			for _, id := range tt.expectIDs {
				assert.Contains(t, result, id)
				assert.Equal(t, id, result[id].AdminID)
			}
		})
	}
}
