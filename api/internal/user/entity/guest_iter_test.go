package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuests_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		guests Guests
	}{
		{
			name: "success",
			guests: Guests{
				{UserID: "user-id01", Lastname: "山田"},
				{UserID: "user-id02", Lastname: "田中"},
			},
		},
		{
			name:   "empty",
			guests: Guests{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, g := range tt.guests.All() {
				indices = append(indices, i)
				ids = append(ids, g.UserID)
			}
			for i, g := range tt.guests {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, g.UserID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.guests))
		})
	}
}

func TestGuests_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	guests := Guests{
		{UserID: "user-id01"},
		{UserID: "user-id02"},
		{UserID: "user-id03"},
	}
	var count int
	for range guests.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestGuests_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		guests    Guests
		expectIDs []string
	}{
		{
			name: "success",
			guests: Guests{
				{UserID: "user-id01"},
				{UserID: "user-id02"},
			},
			expectIDs: []string{"user-id01", "user-id02"},
		},
		{
			name:      "empty",
			guests:    Guests{},
			expectIDs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Guest)
			for k, v := range tt.guests.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.guests))
			for _, id := range tt.expectIDs {
				assert.Contains(t, result, id)
				assert.Equal(t, id, result[id].UserID)
			}
		})
	}
}
