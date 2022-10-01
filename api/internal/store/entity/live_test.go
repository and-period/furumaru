package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestLive(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewLiveParams
		expect *Live
	}{
		{
			name: "success",
			params: &NewLiveParams{
				ScheduleID:  "schedule-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				ProducerID:  "producer-id",
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				Recommends:  []string{"product-id1", "product-id2"},
			},
			expect: &Live{
				ScheduleID:  "schedule-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				ProducerID:  "producer-id",
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				Recommends:  []string{"product-id1", "product-id2"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewLive(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestLive_FillJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		live   *Live
		expect *Live
		hasErr bool
	}{
		{
			name: "success",
			live: &Live{
				ID:          "live-id",
				Title:       "ライブのタイトル",
				Description: "ライブの説明",
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				Recommends:  []string{"product-id1", "product-id2"},
			},
			expect: &Live{
				ID:             "live-id",
				Title:          "ライブのタイトル",
				Description:    "ライブの説明",
				StartAt:        jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:          jst.Date(2022, 9, 1, 0, 0, 0, 0),
				Recommends:     []string{"product-id1", "product-id2"},
				RecommendsJSON: datatypes.JSON([]byte(`["product-id1","product-id2"]`)),
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.live.FillJSON()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.live)
		})
	}
}

func TestLive_Marshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		recommends []string
		expect     []byte
		hasErr     bool
	}{
		{
			name:       "success",
			recommends: []string{"product-id1", "product-id2"},
			expect:     []byte(`["product-id1", "product-id2"]`),
			hasErr:     false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := Marshal(tt.recommends)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
