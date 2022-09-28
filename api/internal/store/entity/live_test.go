package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
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
