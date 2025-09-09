package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/stretchr/testify/assert"
)

func TestExperienceTypes(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name            string
		experienceTypes entity.ExperienceTypes
		expect          ExperienceTypes
	}{
		{
			name: "success",
			experienceTypes: entity.ExperienceTypes{
				{
					ID:        "1",
					Name:      "じゃがいも収穫",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			expect: ExperienceTypes{
				{
					ExperienceType: types.ExperienceType{
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
		expect          []*types.ExperienceType
	}{
		{
			name: "success",
			experienceTypes: ExperienceTypes{
				{
					ExperienceType: types.ExperienceType{
						ID:   "1",
						Name: "じゃがいも収穫",
					},
				},
			},
			expect: []*types.ExperienceType{
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
			actual := tt.experienceTypes.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
