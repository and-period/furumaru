package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExperienceTypes_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		types ExperienceTypes
	}{
		{
			name: "success",
			types: ExperienceTypes{
				{ID: "exp-type-01", Name: "収穫体験"},
				{ID: "exp-type-02", Name: "料理体験"},
			},
		},
		{
			name:  "empty",
			types: ExperienceTypes{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, et := range tt.types.All() {
				indices = append(indices, i)
				ids = append(ids, et.ID)
			}
			for i, et := range tt.types {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, et.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.types))
		})
	}
}

func TestExperienceTypes_IterMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		types ExperienceTypes
	}{
		{
			name: "success",
			types: ExperienceTypes{
				{ID: "exp-type-01", Name: "収穫体験"},
				{ID: "exp-type-02", Name: "料理体験"},
			},
		},
		{
			name:  "empty",
			types: ExperienceTypes{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*ExperienceType)
			for k, v := range tt.types.IterMap() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.types))
			for _, et := range tt.types {
				assert.Contains(t, result, et.ID)
			}
		})
	}
}
