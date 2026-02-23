package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddresses_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		addresses Addresses
	}{
		{
			name: "success",
			addresses: Addresses{
				{ID: "address-id01", UserID: "user-id01"},
				{ID: "address-id02", UserID: "user-id02"},
			},
		},
		{
			name:      "empty",
			addresses: Addresses{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, a := range tt.addresses.All() {
				indices = append(indices, i)
				ids = append(ids, a.ID)
			}
			for i, a := range tt.addresses {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, a.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.addresses))
		})
	}
}

func TestAddresses_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	addresses := Addresses{
		{ID: "address-id01"},
		{ID: "address-id02"},
		{ID: "address-id03"},
	}
	var count int
	for range addresses.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestAddresses_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		addresses Addresses
		expectIDs []string
	}{
		{
			name: "success",
			addresses: Addresses{
				{ID: "address-id01"},
				{ID: "address-id02"},
			},
			expectIDs: []string{"address-id01", "address-id02"},
		},
		{
			name:      "empty",
			addresses: Addresses{},
			expectIDs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Address)
			for k, v := range tt.addresses.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.addresses))
			for _, id := range tt.expectIDs {
				assert.Contains(t, result, id)
				assert.Equal(t, id, result[id].ID)
			}
		})
	}
}

func TestAddresses_IterMapByRevision(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		addresses Addresses
		expect    map[int64]string
	}{
		{
			name: "success",
			addresses: Addresses{
				{
					ID:              "address-id01",
					AddressRevision: AddressRevision{ID: 1, AddressID: "address-id01"},
				},
				{
					ID:              "address-id02",
					AddressRevision: AddressRevision{ID: 2, AddressID: "address-id02"},
				},
			},
			expect: map[int64]string{
				1: "address-id01",
				2: "address-id02",
			},
		},
		{
			name:      "empty",
			addresses: Addresses{},
			expect:    map[int64]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[int64]*Address)
			for k, v := range tt.addresses.IterMapByRevision() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.expect))
			for revID, addressID := range tt.expect {
				assert.Contains(t, result, revID)
				assert.Equal(t, addressID, result[revID].ID)
			}
		})
	}
}

func TestAddresses_IterMapByUserID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		addresses Addresses
		expect    map[string]string
	}{
		{
			name: "success",
			addresses: Addresses{
				{ID: "address-id01", UserID: "user-id01"},
				{ID: "address-id02", UserID: "user-id02"},
			},
			expect: map[string]string{
				"user-id01": "address-id01",
				"user-id02": "address-id02",
			},
		},
		{
			name:      "empty",
			addresses: Addresses{},
			expect:    map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Address)
			for k, v := range tt.addresses.IterMapByUserID() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.expect))
			for userID, addressID := range tt.expect {
				assert.Contains(t, result, userID)
				assert.Equal(t, addressID, result[userID].ID)
			}
		})
	}
}
