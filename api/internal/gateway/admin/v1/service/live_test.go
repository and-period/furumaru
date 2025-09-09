package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestLive(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name   string
		live   *entity.Live
		expect *Live
	}{
		{
			name: "success",
			live: &entity.Live{
				ID:         "live-id",
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    now.AddDate(0, -1, 0),
				EndAt:      now.AddDate(0, 1, 0),
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			expect: &Live{
				Live: types.Live{
					ID:         "live-id",
					ScheduleID: "schedule-id",
					ProducerID: "producer-id",
					ProductIDs: []string{"product-id"},
					Comment:    "よろしくお願いします。",
					StartAt:    1638284400,
					EndAt:      1643641200,
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewLive(tt.live)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestLive_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		live   *Live
		expect *types.Live
	}{
		{
			name: "success",
			live: &Live{
				Live: types.Live{
					ID:         "live-id",
					ScheduleID: "schedule-id",
					ProducerID: "producer-id",
					ProductIDs: []string{"product-id"},
					Comment:    "よろしくお願いします。",
					StartAt:    1638284400,
					EndAt:      1643641200,
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
				},
			},
			expect: &types.Live{
				ID:         "live-id",
				ScheduleID: "schedule-id",
				ProducerID: "producer-id",
				ProductIDs: []string{"product-id"},
				Comment:    "よろしくお願いします。",
				StartAt:    1638284400,
				EndAt:      1643641200,
				CreatedAt:  1640962800,
				UpdatedAt:  1640962800,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.live.Response())
		})
	}
}

func TestLives(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name   string
		lives  entity.Lives
		expect Lives
	}{
		{
			name: "success",
			lives: entity.Lives{
				{
					ID:         "live-id",
					ScheduleID: "schedule-id",
					ProducerID: "producer-id",
					ProductIDs: []string{"product-id"},
					Comment:    "よろしくお願いします。",
					StartAt:    now.AddDate(0, -1, 0),
					EndAt:      now.AddDate(0, 1, 0),
					CreatedAt:  now,
					UpdatedAt:  now,
				},
			},
			expect: Lives{
				{
					Live: types.Live{
						ID:         "live-id",
						ScheduleID: "schedule-id",
						ProducerID: "producer-id",
						ProductIDs: []string{"product-id"},
						Comment:    "よろしくお願いします。",
						StartAt:    1638284400,
						EndAt:      1643641200,
						CreatedAt:  1640962800,
						UpdatedAt:  1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewLives(tt.lives)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestLives_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		lives  Lives
		expect []*types.Live
	}{
		{
			name: "success",
			lives: Lives{
				{
					Live: types.Live{
						ID:         "live-id",
						ScheduleID: "schedule-id",
						ProducerID: "producer-id",
						ProductIDs: []string{"product-id"},
						Comment:    "よろしくお願いします。",
						StartAt:    1638284400,
						EndAt:      1643641200,
						CreatedAt:  1640962800,
						UpdatedAt:  1640962800,
					},
				},
			},
			expect: []*types.Live{
				{
					ID:         "live-id",
					ScheduleID: "schedule-id",
					ProducerID: "producer-id",
					ProductIDs: []string{"product-id"},
					Comment:    "よろしくお願いします。",
					StartAt:    1638284400,
					EndAt:      1643641200,
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.lives.Response())
		})
	}
}
