package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchedules_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		schedules Schedules
	}{
		{
			name: "success",
			schedules: Schedules{
				{ID: "schedule-01", Title: "春の収穫祭"},
				{ID: "schedule-02", Title: "秋の収穫祭"},
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
				ids = append(ids, s.ID)
			}
			for i, s := range tt.schedules {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, s.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.schedules))
		})
	}
}

func TestSchedules_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	schedules := Schedules{
		{ID: "schedule-01"},
		{ID: "schedule-02"},
		{ID: "schedule-03"},
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
	tests := []struct {
		name      string
		schedules Schedules
	}{
		{
			name: "success",
			schedules: Schedules{
				{ID: "schedule-01", Title: "春の収穫祭"},
				{ID: "schedule-02", Title: "秋の収穫祭"},
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
			result := make(map[string]*Schedule)
			for k, v := range tt.schedules.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.schedules))
			for _, s := range tt.schedules {
				assert.Contains(t, result, s.ID)
			}
		})
	}
}
