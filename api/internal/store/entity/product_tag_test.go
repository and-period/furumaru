package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductTag(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		input  string
		expect *ProductTag
	}{
		{
			name:  "success",
			input: "贈答品",
			expect: &ProductTag{
				Name: "贈答品",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductTag(tt.input)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}
