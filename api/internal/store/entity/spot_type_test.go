package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpotType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewSpotTypeParams
		expect *SpotType
	}{
		{
			name: "success",
			params: &NewSpotTypeParams{
				Name: "スポット種別",
			},
			expect: &SpotType{
				Name: "スポット種別",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewSpotType(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}
