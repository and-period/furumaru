package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		expect *Category
	}{
		{
			name:  "success",
			input: "野菜",
			expect: &Category{
				Name: "野菜",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewCategory(tt.input)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}
