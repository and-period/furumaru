package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/stretchr/testify/assert"
)

func TestSpotTypes(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name      string
		spotTypes entity.SpotTypes
		expect    SpotTypes
	}{
		{
			name: "success",
			spotTypes: entity.SpotTypes{
				{
					ID:        "1",
					Name:      "じゃがいも収穫",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			expect: SpotTypes{
				{
					SpotType: response.SpotType{
						ID:   "1",
						Name: "じゃがいも収穫",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
		expect    []*response.SpotType
	}{
		{
			name: "success",
			spotTypes: SpotTypes{
				{
					SpotType: response.SpotType{
						ID:   "1",
						Name: "じゃがいも収穫",
					},
				},
			},
			expect: []*response.SpotType{
				{
					ID:   "1",
					Name: "じゃがいも収穫",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.spotTypes.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
