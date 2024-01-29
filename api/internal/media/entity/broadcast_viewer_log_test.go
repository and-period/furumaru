package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBroadcastViewerLog(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *BroadcastViewerLogParams
		expect *BroadcastViewerLog
	}{
		{
			name: "success",
			params: &BroadcastViewerLogParams{
				BroadcastID: "broadcast-id",
				SessionID:   "session-id",
				UserID:      "user-id",
				UserAgent:   "user-agent",
				ClientIP:    "127.0.0.1",
			},
			expect: &BroadcastViewerLog{
				BroadcastID: "broadcast-id",
				SessionID:   "session-id",
				UserID:      "user-id",
				UserAgent:   "user-agent",
				ClientIP:    "127.0.0.1",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewBroadcastViewerLog(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
