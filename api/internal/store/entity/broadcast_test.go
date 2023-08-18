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

func TestBroadcast_MediaLiveChannelID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		broadcast *Broadcast
		expect    string
	}{
		{
			name: "success",
			broadcast: &Broadcast{
				MediaLiveChannelArn: "arn:aws:medialive:ap-northeast-1:123456789012:channel:123456",
			},
			expect: "123456",
		},
		{
			name: "empty",
			broadcast: &Broadcast{
				MediaLiveChannelArn: "",
			},
			expect: "",
		},
		{
			name:      "nil",
			broadcast: nil,
			expect:    "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.broadcast.MediaLiveChannelID())
		})
	}
}

func TestBroadcast_MediaLiveRTMPInputID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		broadcast *Broadcast
		expect    string
	}{
		{
			name: "success",
			broadcast: &Broadcast{
				MediaLiveRTMPInputArn: "arn:aws:medialive:ap-northeast-1:123456789012:input:123456",
			},
			expect: "123456",
		},
		{
			name: "empty",
			broadcast: &Broadcast{
				MediaLiveRTMPInputArn: "",
			},
			expect: "",
		},
		{
			name:      "nil",
			broadcast: nil,
			expect:    "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.broadcast.MediaLiveRTMPInputID())
		})
	}
}

func TestBroadcast_MediaLiveMP4InputID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		broadcast *Broadcast
		expect    string
	}{
		{
			name: "success",
			broadcast: &Broadcast{
				MediaLiveMP4InputArn: "arn:aws:medialive:ap-northeast-1:123456789012:input:123456",
			},
			expect: "123456",
		},
		{
			name: "empty",
			broadcast: &Broadcast{
				MediaLiveMP4InputArn: "",
			},
			expect: "",
		},
		{
			name:      "nil",
			broadcast: nil,
			expect:    "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.broadcast.MediaLiveMP4InputID())
		})
	}
}
