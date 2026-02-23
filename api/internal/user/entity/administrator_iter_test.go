package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdministrators_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		administrators Administrators
	}{
		{
			name: "success",
			administrators: Administrators{
				{AdminID: "admin-id01"},
				{AdminID: "admin-id02"},
			},
		},
		{
			name:           "empty",
			administrators: Administrators{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, a := range tt.administrators.All() {
				indices = append(indices, i)
				ids = append(ids, a.AdminID)
			}
			for i, a := range tt.administrators {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, a.AdminID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.administrators))
		})
	}
}

func TestAdministrators_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	administrators := Administrators{
		{AdminID: "admin-id01"},
		{AdminID: "admin-id02"},
		{AdminID: "admin-id03"},
	}
	var count int
	for range administrators.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestAdministrators_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		administrators Administrators
		expectIDs      []string
	}{
		{
			name: "success",
			administrators: Administrators{
				{AdminID: "admin-id01"},
				{AdminID: "admin-id02"},
			},
			expectIDs: []string{"admin-id01", "admin-id02"},
		},
		{
			name:           "empty",
			administrators: Administrators{},
			expectIDs:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Administrator)
			for k, v := range tt.administrators.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.administrators))
			for _, id := range tt.expectIDs {
				assert.Contains(t, result, id)
				assert.Equal(t, id, result[id].AdminID)
			}
		})
	}
}
