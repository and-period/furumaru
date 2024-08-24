package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestExperienceTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		experienceTypes entity.ExperienceTypes
		expect          ExperienceTypes
	}{
		{
			name: "success",
			experienceTypes: entity.ExperienceTypes{
				{
					ID:        "experience-type-id",
					Name:      "じゃがいも収穫",
					CreatedAt: jst.Date(2024, 8, 24, 18, 30, 0, 0),
					UpdatedAt: jst.Date(2024, 8, 24, 18, 30, 0, 0),
				},
			},
			expect: ExperienceTypes{
				{
					ExperienceType: response.ExperienceType{
						ID:        "experience-type-id",
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
			actual := NewExperienceTypes(tt.experienceTypes)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperienceTypes_Response(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		experienceTypes ExperienceTypes
		expect          []*response.ExperienceType
	}{
		{
			name: "success",
			experienceTypes: ExperienceTypes{
				{
					ExperienceType: response.ExperienceType{
						ID:        "experience-type-id",
						Name:      "じゃがいも収穫",
						CreatedAt: 1724491800,
						UpdatedAt: 1724491800,
					},
				},
			},
			expect: []*response.ExperienceType{
				{
					ID:        "experience-type-id",
					Name:      "じゃがいも収穫",
					CreatedAt: 1724491800,
					UpdatedAt: 1724491800,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.experienceTypes.Response())
		})
	}
}
