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
				ScheduleID: "schedule-id",
			},
			expect: &Broadcast{
				ID:         "",
				ScheduleID: "schedule-id",
				Status:     BroadcastStatusDisabled,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewBroadcast(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}
