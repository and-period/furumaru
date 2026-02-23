package collection

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterIter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		input  []int
		fn     func(int) bool
		expect []int
	}{
		{
			name:  "filter even numbers",
			input: []int{1, 2, 3, 4, 5, 6},
			fn: func(v int) bool {
				return v%2 == 0
			},
			expect: []int{2, 4, 6},
		},
		{
			name:  "filter none",
			input: []int{1, 3, 5},
			fn: func(v int) bool {
				return v%2 == 0
			},
			expect: nil,
		},
		{
			name:  "empty slice",
			input: []int{},
			fn: func(v int) bool {
				return true
			},
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := slices.Collect(FilterIter(tt.input, tt.fn))
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestFilterIter_EarlyBreak(t *testing.T) {
	t.Parallel()
	input := []int{1, 2, 3, 4, 5, 6}
	seq := FilterIter(input, func(v int) bool {
		return v%2 == 0
	})
	// Only take the first matching element
	var first int
	for v := range seq {
		first = v
		break
	}
	assert.Equal(t, 2, first)
}

func TestMapIter(t *testing.T) {
	t.Parallel()
	type item struct {
		ID   string
		Name string
	}
	tests := []struct {
		name   string
		input  []item
		expect map[string]string
	}{
		{
			name: "map items by ID",
			input: []item{
				{ID: "a", Name: "Alice"},
				{ID: "b", Name: "Bob"},
			},
			expect: map[string]string{
				"a": "Alice",
				"b": "Bob",
			},
		},
		{
			name:   "empty slice",
			input:  []item{},
			expect: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			seq := MapIter(tt.input, func(i item) (string, string) {
				return i.ID, i.Name
			})
			result := make(map[string]string)
			for k, v := range seq {
				result[k] = v
			}
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestMapIter_EarlyBreak(t *testing.T) {
	t.Parallel()
	items := []int{1, 2, 3}
	result := make(map[int]int)
	for k, v := range MapIter(items, func(i int) (int, int) { return i, i * 10 }) {
		result[k] = v
		if len(result) >= 1 {
			break
		}
	}
	assert.Len(t, result, 1)
}
