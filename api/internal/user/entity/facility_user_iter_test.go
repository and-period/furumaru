package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFacilityUsers_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		facilityUsers FacilityUsers
	}{
		{
			name: "success",
			facilityUsers: FacilityUsers{
				{UserID: "user-id01", Lastname: "山田"},
				{UserID: "user-id02", Lastname: "田中"},
			},
		},
		{
			name:          "empty",
			facilityUsers: FacilityUsers{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, f := range tt.facilityUsers.All() {
				indices = append(indices, i)
				ids = append(ids, f.UserID)
			}
			for i, f := range tt.facilityUsers {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, f.UserID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.facilityUsers))
		})
	}
}

func TestFacilityUsers_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	facilityUsers := FacilityUsers{
		{UserID: "user-id01"},
		{UserID: "user-id02"},
		{UserID: "user-id03"},
	}
	var count int
	for range facilityUsers.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestFacilityUsers_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		facilityUsers FacilityUsers
		expectIDs     []string
	}{
		{
			name: "success",
			facilityUsers: FacilityUsers{
				{UserID: "user-id01"},
				{UserID: "user-id02"},
			},
			expectIDs: []string{"user-id01", "user-id02"},
		},
		{
			name:          "empty",
			facilityUsers: FacilityUsers{},
			expectIDs:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*FacilityUser)
			for k, v := range tt.facilityUsers.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.facilityUsers))
			for _, id := range tt.expectIDs {
				assert.Contains(t, result, id)
				assert.Equal(t, id, result[id].UserID)
			}
		})
	}
}
