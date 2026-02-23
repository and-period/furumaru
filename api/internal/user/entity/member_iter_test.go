package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMembers_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		members Members
	}{
		{
			name: "success",
			members: Members{
				{UserID: "user-id01", Username: "会員A"},
				{UserID: "user-id02", Username: "会員B"},
			},
		},
		{
			name:    "empty",
			members: Members{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, m := range tt.members.All() {
				indices = append(indices, i)
				ids = append(ids, m.UserID)
			}
			for i, m := range tt.members {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, m.UserID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.members))
		})
	}
}

func TestMembers_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	members := Members{
		{UserID: "user-id01"},
		{UserID: "user-id02"},
		{UserID: "user-id03"},
	}
	var count int
	for range members.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestMembers_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		members   Members
		expectIDs []string
	}{
		{
			name: "success",
			members: Members{
				{UserID: "user-id01"},
				{UserID: "user-id02"},
			},
			expectIDs: []string{"user-id01", "user-id02"},
		},
		{
			name:      "empty",
			members:   Members{},
			expectIDs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Member)
			for k, v := range tt.members.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.members))
			for _, id := range tt.expectIDs {
				assert.Contains(t, result, id)
				assert.Equal(t, id, result[id].UserID)
			}
		})
	}
}
