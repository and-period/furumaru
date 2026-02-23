package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLives_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		lives Lives
	}{
		{
			name: "success",
			lives: Lives{
				{ID: "live-01", ScheduleID: "schedule-01"},
				{ID: "live-02", ScheduleID: "schedule-01"},
			},
		},
		{
			name:  "empty",
			lives: Lives{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, l := range tt.lives.All() {
				indices = append(indices, i)
				ids = append(ids, l.ID)
			}
			for i, l := range tt.lives {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, l.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.lives))
		})
	}
}

func TestLives_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	lives := Lives{
		{ID: "live-01"},
		{ID: "live-02"},
		{ID: "live-03"},
	}
	var count int
	for range lives.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestLives_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		lives Lives
	}{
		{
			name: "success",
			lives: Lives{
				{ID: "live-01", ScheduleID: "schedule-01"},
				{ID: "live-02", ScheduleID: "schedule-02"},
			},
		},
		{
			name:  "empty",
			lives: Lives{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*Live)
			for k, v := range tt.lives.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.lives))
			for _, l := range tt.lives {
				assert.Contains(t, result, l.ID)
			}
		})
	}
}

func TestLives_IterGroupByScheduleID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lives      Lives
		expectKeys int
	}{
		{
			name: "success",
			lives: Lives{
				{ID: "live-01", ScheduleID: "schedule-01"},
				{ID: "live-02", ScheduleID: "schedule-01"},
				{ID: "live-03", ScheduleID: "schedule-02"},
			},
			expectKeys: 2,
		},
		{
			name:       "empty",
			lives:      Lives{},
			expectKeys: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]Lives)
			for k, v := range tt.lives.IterGroupByScheduleID() {
				result[k] = v
			}
			assert.Len(t, result, tt.expectKeys)
		})
	}
}
