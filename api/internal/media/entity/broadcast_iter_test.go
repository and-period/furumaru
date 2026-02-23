package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBroadcasts_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		broadcasts Broadcasts
	}{
		{
			name: "success",
			broadcasts: Broadcasts{
				{ID: "broadcast-id01", ScheduleID: "schedule-id01"},
				{ID: "broadcast-id02", ScheduleID: "schedule-id02"},
			},
		},
		{
			name:       "empty",
			broadcasts: Broadcasts{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, b := range tt.broadcasts.All() {
				indices = append(indices, i)
				ids = append(ids, b.ID)
			}
			for i, b := range tt.broadcasts {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, b.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.broadcasts))
		})
	}
}

func TestBroadcasts_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	broadcasts := Broadcasts{
		{ID: "broadcast-id01"},
		{ID: "broadcast-id02"},
		{ID: "broadcast-id03"},
	}
	var count int
	for range broadcasts.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestBroadcasts_IterScheduleIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		broadcasts Broadcasts
		expect     []string
	}{
		{
			name: "success",
			broadcasts: Broadcasts{
				{ID: "broadcast-id01", ScheduleID: "schedule-id01"},
				{ID: "broadcast-id02", ScheduleID: "schedule-id02"},
			},
			expect: []string{"schedule-id01", "schedule-id02"},
		},
		{
			name:       "empty",
			broadcasts: Broadcasts{},
			expect:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var actual []string
			for id := range tt.broadcasts.IterScheduleIDs() {
				actual = append(actual, id)
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestBroadcasts_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		broadcasts Broadcasts
		expectIDs  []string
	}{
		{
			name: "success",
			broadcasts: Broadcasts{
				{ID: "broadcast-id01", ScheduleID: "schedule-id01"},
				{ID: "broadcast-id02", ScheduleID: "schedule-id02"},
			},
			expectIDs: []string{"broadcast-id01", "broadcast-id02"},
		},
		{
			name:       "empty",
			broadcasts: Broadcasts{},
			expectIDs:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Broadcast)
			for k, v := range tt.broadcasts.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.broadcasts))
			for _, id := range tt.expectIDs {
				assert.Contains(t, result, id)
				assert.Equal(t, id, result[id].ID)
			}
		})
	}
}
