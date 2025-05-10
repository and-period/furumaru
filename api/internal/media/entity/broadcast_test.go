package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBroadcast(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewBroadcastParams
		expect *Broadcast
	}{
		{
			name: "success",
			params: &NewBroadcastParams{
				ScheduleID:    "schedule-id",
				CoordinatorID: "coordinator-id",
			},
			expect: &Broadcast{
				ID:            "",
				ScheduleID:    "schedule-id",
				CoordinatorID: "coordinator-id",
				Type:          BroadcastTypeNormal,
				Status:        BroadcastStatusDisabled,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewBroadcast(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestBroadcasts_ScheduleID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		broadcasts Broadcasts
		expect     []string
	}{
		{
			name: "suceess",
			broadcasts: Broadcasts{
				{ID: "broadcast-id-01", ScheduleID: "schedule-id01"},
				{ID: "broadcast-id-02", ScheduleID: "schedule-id01"},
				{ID: "broadcast-id-03", ScheduleID: "schedule-id02"},
			},
			expect: []string{
				"schedule-id01",
				"schedule-id02",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.broadcasts.ScheduleIDs())
		})
	}
}
