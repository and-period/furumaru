package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsers_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		users Users
	}{
		{
			name: "success",
			users: Users{
				{ID: "user-id01", Registered: true},
				{ID: "user-id02", Registered: false},
				{ID: "user-id03", Registered: true},
			},
		},
		{
			name:  "empty",
			users: Users{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, u := range tt.users.All() {
				indices = append(indices, i)
				ids = append(ids, u.ID)
			}
			for i, u := range tt.users {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, u.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.users))
		})
	}
}

func TestUsers_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	users := Users{
		{ID: "user-id01"},
		{ID: "user-id02"},
		{ID: "user-id03"},
	}
	var count int
	for range users.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestUsers_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		users     Users
		expectIDs []string
	}{
		{
			name: "success",
			users: Users{
				{ID: "user-id01"},
				{ID: "user-id02"},
			},
			expectIDs: []string{"user-id01", "user-id02"},
		},
		{
			name:      "empty",
			users:     Users{},
			expectIDs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*User)
			for k, v := range tt.users.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.users))
			for _, id := range tt.expectIDs {
				assert.Contains(t, result, id)
				assert.Equal(t, id, result[id].ID)
			}
		})
	}
}

func TestUsers_IterGroupByRegistered(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		users  Users
		expect map[bool][]string
	}{
		{
			name: "success",
			users: Users{
				{ID: "user-id01", Registered: true},
				{ID: "user-id02", Registered: false},
				{ID: "user-id03", Registered: true},
			},
			expect: map[bool][]string{
				true:  {"user-id01", "user-id03"},
				false: {"user-id02"},
			},
		},
		{
			name:   "empty",
			users:  Users{},
			expect: map[bool][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[bool][]string)
			for k, v := range tt.users.IterGroupByRegistered() {
				result[k] = append(result[k], v.ID)
			}
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestUsers_IterGroupByUserType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		users  Users
		expect map[UserType][]string
	}{
		{
			name: "success",
			users: Users{
				{ID: "user-id01", Type: UserTypeMember},
				{ID: "user-id02", Type: UserTypeGuest},
				{ID: "user-id03", Type: UserTypeMember},
			},
			expect: map[UserType][]string{
				UserTypeMember: {"user-id01", "user-id03"},
				UserTypeGuest:  {"user-id02"},
			},
		},
		{
			name:   "empty",
			users:  Users{},
			expect: map[UserType][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[UserType][]string)
			for k, v := range tt.users.IterGroupByUserType() {
				result[k] = append(result[k], v.ID)
			}
			assert.Equal(t, tt.expect, result)
		})
	}
}
