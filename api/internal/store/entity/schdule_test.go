package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestSchedule(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewScheduleParams
		expect *Schedule
	}{
		{
			name: "success",
			params: &NewScheduleParams{
				Title:       "スケジュールタイトル",
				Description: "スケジュールの詳細です。",
				ThumnailUrl: "サムネイルのURLです",
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				Canceled:    false,
			},
			expect: &Schedule{
				Title:       "スケジュールタイトル",
				Description: "スケジュールの詳細です。",
				ThumnailUrl: "サムネイルのURLです",
				StartAt:     jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 9, 1, 0, 0, 0, 0),
				Canceled:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewSchedule(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}
