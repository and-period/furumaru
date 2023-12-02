package entity

import (
	"testing"
	"time"

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
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				Comment:    "よろしくお願いします。",
				ProductIDs: []string{},
				StartAt:    jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:      jst.Date(2022, 9, 1, 0, 0, 0, 0),
			},
			expect: &Live{
				ScheduleID:   "schedule-id",
				ProducerID:   "producer-id",
				ProductIDs:   []string{},
				Comment:      "よろしくお願いします。",
				StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
				LiveProducts: LiveProducts{},
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

func TestLive_Fill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		live     *Live
		products LiveProducts
		now      time.Time
		expect   *Live
	}{
		{
			name: "success",
			live: &Live{
				ID:         "live-id",
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:      jst.Date(2022, 9, 1, 0, 0, 0, 0),
				CreatedAt:  jst.Date(2022, 7, 1, 0, 0, 0, 0),
				UpdatedAt:  jst.Date(2022, 7, 1, 0, 0, 0, 0),
			},
			products: LiveProducts{
				{
					LiveID:    "live-id",
					ProductID: "product-id",
					CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				},
			},
			now: jst.Date(2022, 7, 1, 0, 0, 0, 0),
			expect: &Live{
				ID:         "live-id",
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    jst.Date(2022, 8, 1, 0, 0, 0, 0),
				EndAt:      jst.Date(2022, 9, 1, 0, 0, 0, 0),
				LiveProducts: LiveProducts{
					{
						LiveID:    "live-id",
						ProductID: "product-id",
						CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
						UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					},
				},
				CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.live.Fill(tt.products)
			assert.Equal(t, tt.expect, tt.live)
		})
	}
}

func TestLives_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		live     *Live
		schedule *Schedule
		lives    Lives
		expect   error
	}{
		{
			name: "success empty",
			live: &Live{
				ID:         "live-id",
				ScheduleID: "schedule-id",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			schedule: &Schedule{
				ID:      "schedule-id",
				StartAt: jst.Date(2023, 7, 15, 17, 0, 0, 0),
				EndAt:   jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			lives:  Lives{},
			expect: nil,
		},
		{
			name: "success same live",
			live: &Live{
				ID:         "live-id",
				ScheduleID: "schedule-id",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			schedule: &Schedule{
				ID:      "schedule-id",
				StartAt: jst.Date(2023, 7, 15, 17, 0, 0, 0),
				EndAt:   jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			lives: Lives{
				{
					ID:         "live-id",
					ScheduleID: "schedule-id",
					StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
					EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
				},
			},
			expect: nil,
		},
		{
			name: "success before",
			live: &Live{
				ID:         "live-id",
				ScheduleID: "schedule-id",
				StartAt:    jst.Date(2023, 7, 15, 17, 0, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 18, 30, 0, 0),
			},
			schedule: &Schedule{
				ID:      "schedule-id",
				StartAt: jst.Date(2023, 7, 15, 17, 0, 0, 0),
				EndAt:   jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			lives: Lives{
				{
					ID:         "other-id",
					ScheduleID: "schedule-id",
					StartAt:    jst.Date(2023, 7, 15, 18, 30, 0, 0),
					EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
				},
			},
			expect: nil,
		},
		{
			name: "success after",
			live: &Live{
				ID:         "live-id",
				ScheduleID: "schedule-id",
				StartAt:    jst.Date(2023, 7, 15, 18, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			schedule: &Schedule{
				ID:      "schedule-id",
				StartAt: jst.Date(2023, 7, 15, 17, 0, 0, 0),
				EndAt:   jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			lives: Lives{
				{
					ID:         "other-id",
					ScheduleID: "schedule-id",
					StartAt:    jst.Date(2023, 7, 15, 17, 0, 0, 0),
					EndAt:      jst.Date(2023, 7, 15, 18, 30, 0, 0),
				},
			},
			expect: nil,
		},
		{
			name: "not found schedule",
			live: &Live{
				ID:         "live-id",
				ScheduleID: "schedule-id",
				StartAt:    jst.Date(2023, 7, 15, 19, 30, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			schedule: nil,
			lives:    Lives{},
			expect:   errNotFoundSchedule,
		},
		{
			name: "invalid schedule",
			live: &Live{
				ID:         "live-id",
				ScheduleID: "schedule-id",
				StartAt:    jst.Date(2023, 7, 15, 12, 0, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 14, 30, 0, 0),
			},
			schedule: &Schedule{
				ID:      "schedule-id",
				StartAt: jst.Date(2023, 7, 15, 17, 0, 0, 0),
				EndAt:   jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			lives:  Lives{},
			expect: errInvalidLiveSchedule,
		},
		{
			name: "invalid live schedule",
			live: &Live{
				ID:         "live-id",
				ScheduleID: "schedule-id",
				StartAt:    jst.Date(2023, 7, 15, 18, 00, 0, 0),
				EndAt:      jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			schedule: &Schedule{
				ID:      "schedule-id",
				StartAt: jst.Date(2023, 7, 15, 17, 0, 0, 0),
				EndAt:   jst.Date(2023, 7, 15, 21, 0, 0, 0),
			},
			lives: Lives{
				{
					ID:         "other-id",
					ScheduleID: "schedule-id",
					StartAt:    jst.Date(2023, 7, 15, 17, 0, 0, 0),
					EndAt:      jst.Date(2023, 7, 15, 19, 0, 0, 0),
				},
			},
			expect: errInvalidLiveSchedule,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.live.Validate(tt.schedule, tt.lives)
			assert.ErrorIs(t, err, tt.expect)
		})
	}
}

func TestLives_IDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		lives  Lives
		expect []string
	}{
		{
			name: "success",
			lives: Lives{
				{
					ID:         "live-id",
					ScheduleID: "schedule-id",
					ProducerID: "producer-id",
					Comment:    "よろしくお願いします。",
					ProductIDs: []string{"product-id"},
					StartAt:    jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:      jst.Date(2022, 9, 1, 0, 0, 0, 0),
					LiveProducts: LiveProducts{
						{
							LiveID:    "live-id",
							ProductID: "product-id",
							CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
							UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
						},
					},
					CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				},
			},
			expect: []string{"live-id"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.lives.IDs())
		})
	}
}

func TestLives_ProducerIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		lives  Lives
		expect []string
	}{
		{
			name: "success",
			lives: Lives{
				{
					ID:         "live-id",
					ScheduleID: "schedule-id",
					ProducerID: "producer-id",
					Comment:    "よろしくお願いします。",
					ProductIDs: []string{"product-id"},
					StartAt:    jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:      jst.Date(2022, 9, 1, 0, 0, 0, 0),
					LiveProducts: LiveProducts{
						{
							LiveID:    "live-id",
							ProductID: "product-id",
							CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
							UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
						},
					},
					CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				},
			},
			expect: []string{"producer-id"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.lives.ProducerIDs())
		})
	}
}

func TestLives_ProductIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		lives  Lives
		expect []string
	}{
		{
			name: "success",
			lives: Lives{
				{
					ID:         "live-id",
					ScheduleID: "schedule-id",
					ProducerID: "producer-id",
					Comment:    "よろしくお願いします。",
					ProductIDs: []string{"product-id"},
					StartAt:    jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:      jst.Date(2022, 9, 1, 0, 0, 0, 0),
					LiveProducts: LiveProducts{
						{
							LiveID:    "live-id",
							ProductID: "product-id",
							CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
							UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
						},
					},
					CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				},
			},
			expect: []string{"product-id"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.lives.ProductIDs())
		})
	}
}

func TestLives_Fill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		lives    Lives
		products map[string]LiveProducts
		now      time.Time
		expect   Lives
	}{
		{
			name: "success",
			lives: Lives{
				{
					ID:         "live-id",
					ScheduleID: "schedule-id",
					ProducerID: "producer-id",
					ProductIDs: []string{"product-id"},
					Comment:    "よろしくお願いします。",
					StartAt:    jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:      jst.Date(2022, 9, 1, 0, 0, 0, 0),
					CreatedAt:  jst.Date(2022, 7, 1, 0, 0, 0, 0),
					UpdatedAt:  jst.Date(2022, 7, 1, 0, 0, 0, 0),
				},
			},
			products: map[string]LiveProducts{
				"live-id": {
					{
						LiveID:    "live-id",
						ProductID: "product-id",
						CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
						UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					},
				},
			},
			now: jst.Date(2022, 7, 1, 0, 0, 0, 0),
			expect: Lives{
				{
					ID:         "live-id",
					ScheduleID: "schedule-id",
					ProducerID: "producer-id",
					ProductIDs: []string{"product-id"},
					Comment:    "よろしくお願いします。",
					StartAt:    jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:      jst.Date(2022, 9, 1, 0, 0, 0, 0),
					LiveProducts: LiveProducts{
						{
							LiveID:    "live-id",
							ProductID: "product-id",
							CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
							UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
						},
					},
					CreatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 1, 0, 0, 0, 0),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.lives.Fill(tt.products)
			assert.Equal(t, tt.expect, tt.lives)
		})
	}
}

func TestLives_GroupByScheduleID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		lives  Lives
		expect map[string]Lives
	}{
		{
			name: "success",
			lives: Lives{
				{
					ID:         "live-id",
					ScheduleID: "schedule-id",
				},
			},
			expect: map[string]Lives{
				"schedule-id": {
					{
						ID:         "live-id",
						ScheduleID: "schedule-id",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.lives.GroupByScheduleID())
		})
	}
}
