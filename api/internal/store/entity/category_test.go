package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewCategoryParams
		input  string
		expect *Category
	}{
		{
			name: "success",
			params: &NewCategoryParams{
				Name: "野菜",
			},
			expect: &Category{
				Name: "野菜",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewCategory(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}
