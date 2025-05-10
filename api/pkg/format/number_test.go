package format

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRound(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input  float64
		places float64
		expect float64
	}{
		{
			input:  10.123,
			places: 1,
			expect: 10.1,
		},
		{
			input:  10.567,
			places: 1,
			expect: 10.6,
		},
		{
			input:  10.987,
			places: 1,
			expect: 11,
		},
		{
			input:  10.123,
			places: 2,
			expect: 10.12,
		},
		{
			input:  10.987,
			places: 2,
			expect: 10.99,
		},
		{
			input:  10.123,
			places: 0,
			expect: 10,
		},
	}
	for _, tt := range tests {

		t.Run(fmt.Sprintf("input=%f, places=%f", tt.input, tt.places), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, Round(tt.input, tt.places))
		})
	}
}
