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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewBroadcast(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestBroadcast_FillJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		broadcast *Broadcast
		expect    *Broadcast
	}{
		{
			name: "success",
			broadcast: &Broadcast{
				ID:            "broadcast-id",
				ScheduleID:    "schedule-id",
				CoordinatorID: "coordinator-id",
				Type:          BroadcastTypeNormal,
				Status:        BroadcastStatusDisabled,
				ArchiveMetadata: &BroadcastArchiveMetadata{
					Text: map[string]string{
						"jpn": "http://example.com/translate-jpn.vtt",
						"eng": "http://example.com/translate-eng.vtt",
					},
				},
			},
			expect: &Broadcast{
				ID:            "broadcast-id",
				ScheduleID:    "schedule-id",
				CoordinatorID: "coordinator-id",
				Type:          BroadcastTypeNormal,
				Status:        BroadcastStatusDisabled,
				ArchiveMetadata: &BroadcastArchiveMetadata{
					Text: map[string]string{
						"jpn": "http://example.com/translate-jpn.vtt",
						"eng": "http://example.com/translate-eng.vtt",
					},
				},
				ArchiveMetadataJSON: []byte(`{"text":{"eng":"http://example.com/translate-eng.vtt","jpn":"http://example.com/translate-jpn.vtt"}}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.broadcast.FillJSON()
			assert.Equal(t, tt.expect, tt.broadcast)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.broadcasts.ScheduleIDs())
		})
	}
}

func TestBroadcasts_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		broadcasts Broadcasts
		expect     Broadcasts
		hasErr     bool
	}{
		{
			name: "success",
			broadcasts: Broadcasts{
				{
					ID:                  "broadcast-id",
					ScheduleID:          "schedule-id",
					CoordinatorID:       "coordinator-id",
					Type:                BroadcastTypeNormal,
					Status:              BroadcastStatusDisabled,
					ArchiveMetadata:     nil,
					ArchiveMetadataJSON: []byte(`{"text":{"eng":"http://example.com/translate-eng.vtt","jpn":"http://example.com/translate-jpn.vtt"}}`),
				},
			},
			expect: Broadcasts{
				{
					ID:            "broadcast-id",
					ScheduleID:    "schedule-id",
					CoordinatorID: "coordinator-id",
					Type:          BroadcastTypeNormal,
					Status:        BroadcastStatusDisabled,
					ArchiveMetadata: &BroadcastArchiveMetadata{
						Text: map[string]string{
							"jpn": "http://example.com/translate-jpn.vtt",
							"eng": "http://example.com/translate-eng.vtt",
						},
					},
					ArchiveMetadataJSON: []byte(`{"text":{"eng":"http://example.com/translate-eng.vtt","jpn":"http://example.com/translate-jpn.vtt"}}`),
				},
			},
			hasErr: false,
		},
		{
			name: "empty",
			broadcasts: Broadcasts{
				{
					ID:                  "broadcast-id",
					ScheduleID:          "schedule-id",
					CoordinatorID:       "coordinator-id",
					Type:                BroadcastTypeNormal,
					Status:              BroadcastStatusDisabled,
					ArchiveMetadata:     nil,
					ArchiveMetadataJSON: nil,
				},
			},
			expect: Broadcasts{
				{
					ID:                  "broadcast-id",
					ScheduleID:          "schedule-id",
					CoordinatorID:       "coordinator-id",
					Type:                BroadcastTypeNormal,
					Status:              BroadcastStatusDisabled,
					ArchiveMetadata:     &BroadcastArchiveMetadata{},
					ArchiveMetadataJSON: nil,
				},
			},
			hasErr: false,
		},
		{
			name: "failed to unmarshal",
			broadcasts: Broadcasts{
				{
					ID:                  "broadcast-id",
					ScheduleID:          "schedule-id",
					CoordinatorID:       "coordinator-id",
					Type:                BroadcastTypeNormal,
					Status:              BroadcastStatusDisabled,
					ArchiveMetadata:     nil,
					ArchiveMetadataJSON: []byte(`{{`),
				},
			},
			expect: Broadcasts{
				{
					ID:                  "broadcast-id",
					ScheduleID:          "schedule-id",
					CoordinatorID:       "coordinator-id",
					Type:                BroadcastTypeNormal,
					Status:              BroadcastStatusDisabled,
					ArchiveMetadata:     &BroadcastArchiveMetadata{},
					ArchiveMetadataJSON: []byte(`{{`),
				},
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.broadcasts.Fill()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.broadcasts)
		})
	}
}
