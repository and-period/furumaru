package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		expect *ContactCategory
	}{
		{
			name:  "success",
			input: "野菜",
			expect: &ContactCategory{
				Title: "野菜",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewContactCategory(tt.input)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}
