package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProducers_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		producers Producers
	}{
		{
			name: "success",
			producers: Producers{
				{AdminID: "admin-id01", Username: "生産者A"},
				{AdminID: "admin-id02", Username: "生産者B"},
			},
		},
		{
			name:      "empty",
			producers: Producers{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, p := range tt.producers.All() {
				indices = append(indices, i)
				ids = append(ids, p.AdminID)
			}
			for i, p := range tt.producers {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, p.AdminID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.producers))
		})
	}
}

func TestProducers_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	producers := Producers{
		{AdminID: "admin-id01"},
		{AdminID: "admin-id02"},
		{AdminID: "admin-id03"},
	}
	var count int
	for range producers.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestProducers_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		producers Producers
		expectIDs []string
	}{
		{
			name: "success",
			producers: Producers{
				{AdminID: "admin-id01"},
				{AdminID: "admin-id02"},
			},
			expectIDs: []string{"admin-id01", "admin-id02"},
		},
		{
			name:      "empty",
			producers: Producers{},
			expectIDs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Producer)
			for k, v := range tt.producers.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.producers))
			for _, id := range tt.expectIDs {
				assert.Contains(t, result, id)
				assert.Equal(t, id, result[id].AdminID)
			}
		})
	}
}
