package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSchedules_All(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name      string
		schedules Schedules
	}{
		{
			name: "success",
			schedules: Schedules{
				{
					MessageType: ScheduleTypeNotification,
					MessageID:   "message-id01",
					Status:      ScheduleStatusDone,
					Count:       1,
					SentAt:      now,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					MessageType: ScheduleTypeStartLive,
					MessageID:   "message-id02",
					Status:      ScheduleStatusWaiting,
					Count:       0,
					SentAt:      now,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
		{
			name:      "empty",
			schedules: Schedules{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, s := range tt.schedules.All() {
				indices = append(indices, i)
				ids = append(ids, s.MessageID)
			}
			for i, s := range tt.schedules {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, s.MessageID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.schedules))
		})
	}
}

func TestSchedules_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	now := time.Now()
	schedules := Schedules{
		{MessageID: "message-id01", SentAt: now},
		{MessageID: "message-id02", SentAt: now},
		{MessageID: "message-id03", SentAt: now},
	}
	var count int
	for range schedules.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestSchedules_IterMap(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name      string
		schedules Schedules
		expectIDs []string
	}{
		{
			name: "success",
			schedules: Schedules{
				{
					MessageType: ScheduleTypeNotification,
					MessageID:   "message-id01",
					Status:      ScheduleStatusDone,
					Count:       1,
					SentAt:      now,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					MessageType: ScheduleTypeStartLive,
					MessageID:   "message-id02",
					Status:      ScheduleStatusWaiting,
					Count:       0,
					SentAt:      now,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			expectIDs: []string{"message-id01", "message-id02"},
		},
		{
			name:      "empty",
			schedules: Schedules{},
			expectIDs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Schedule)
			for k, v := range tt.schedules.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.schedules))
			for _, id := range tt.expectIDs {
				assert.Contains(t, result, id)
				assert.Equal(t, id, result[id].MessageID)
			}
		})
	}
}
