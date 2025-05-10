package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrings(t *testing.T) {
	t.Parallel()
	assert.Len(t, NewStrings(10), 10)
}

func BenchmarkStrings(b *testing.B) {
	const size = 16
	for b.Loop() {
		NewStrings(size)
	}
}
