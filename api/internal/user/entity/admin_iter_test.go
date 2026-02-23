package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdmins_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admins Admins
	}{
		{
			name: "success",
			admins: Admins{
				{ID: "admin-id01", Type: AdminTypeAdministrator},
				{ID: "admin-id02", Type: AdminTypeCoordinator},
				{ID: "admin-id03", Type: AdminTypeProducer},
			},
		},
		{
			name:   "empty",
			admins: Admins{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, a := range tt.admins.All() {
				indices = append(indices, i)
				ids = append(ids, a.ID)
			}
			for i, a := range tt.admins {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, a.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.admins))
		})
	}
}

func TestAdmins_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	admins := Admins{
		{ID: "admin-id01"},
		{ID: "admin-id02"},
		{ID: "admin-id03"},
	}
	var count int
	for range admins.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestAdmins_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		admins    Admins
		expectIDs []string
	}{
		{
			name: "success",
			admins: Admins{
				{ID: "admin-id01"},
				{ID: "admin-id02"},
			},
			expectIDs: []string{"admin-id01", "admin-id02"},
		},
		{
			name:      "empty",
			admins:    Admins{},
			expectIDs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Admin)
			for k, v := range tt.admins.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.admins))
			for _, id := range tt.expectIDs {
				assert.Contains(t, result, id)
				assert.Equal(t, id, result[id].ID)
			}
		})
	}
}

func TestAdmins_IterGroupByType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admins Admins
		expect map[AdminType][]string
	}{
		{
			name: "success",
			admins: Admins{
				{ID: "admin-id01", Type: AdminTypeAdministrator},
				{ID: "admin-id02", Type: AdminTypeCoordinator},
				{ID: "admin-id03", Type: AdminTypeAdministrator},
			},
			expect: map[AdminType][]string{
				AdminTypeAdministrator: {"admin-id01", "admin-id03"},
				AdminTypeCoordinator:   {"admin-id02"},
			},
		},
		{
			name:   "empty",
			admins: Admins{},
			expect: map[AdminType][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[AdminType][]string)
			for k, v := range tt.admins.IterGroupByType() {
				result[k] = append(result[k], v.ID)
			}
			assert.Equal(t, tt.expect, result)
		})
	}
}
