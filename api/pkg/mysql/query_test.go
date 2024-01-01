package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullInt(t *testing.T) {
	t.Parallel()
	t.Run("non null int", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, int(1), NullInt[int](1))
	})
	t.Run("non null int32", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, int32(1), NullInt[int32](1))
	})
	t.Run("non null int64", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, int64(1), NullInt[int64](1))
	})
	t.Run("null int", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, nil, NullInt[int](0))
	})
	t.Run("null int32", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, nil, NullInt[int32](0))
	})
	t.Run("null int64", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, nil, NullInt[int64](0))
	})
}

func TestNullString(t *testing.T) {
	t.Parallel()
	t.Run("non null string", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "hoge", NullString("hoge"))
	})
	t.Run("null string", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, nil, NullString(""))
	})
}
