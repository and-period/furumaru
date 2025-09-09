package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestSpotTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		spotTypes entity.SpotTypes
		expect    SpotTypes
	}{
		{
			name: "success",
			spotTypes: entity.SpotTypes{
				{
					ID:        "spot-type-id",
					Name:      "じゃがいも収穫",
					CreatedAt: jst.Date(2024, 8, 24, 18, 30, 0, 0),
					UpdatedAt: jst.Date(2024, 8, 24, 18, 30, 0, 0),
				},
			},
			expect: SpotTypes{
				{
					SpotType: types.SpotType{
						ID:        "spot-type-id",
						Name:      "じゃがいも収穫",
						CreatedAt: 1724491800,
						UpdatedAt: 1724491800,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewSpotTypes(tt.spotTypes)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestSpotTypes_Response(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		spotTypes SpotTypes
		expect    []*types.SpotType
	}{
		{
			name: "success",
			spotTypes: SpotTypes{
				{
					SpotType: types.SpotType{
						ID:        "spot-type-id",
						Name:      "じゃがいも収穫",
						CreatedAt: 1724491800,
						UpdatedAt: 1724491800,
					},
				},
			},
			expect: []*types.SpotType{
				{
					ID:        "spot-type-id",
					Name:      "じゃがいも収穫",
					CreatedAt: 1724491800,
					UpdatedAt: 1724491800,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.spotTypes.Response())
		})
	}
}
