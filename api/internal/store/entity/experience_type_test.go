package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExperienceType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewExperienceTypeParams
		expect *ExperienceType
	}{
		{
			name: "success",
			params: &NewExperienceTypeParams{
				Name: "体験種別",
			},
			expect: &ExperienceType{
				Name: "体験種別",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewExperienceType(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}
