package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet_New(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, New[int64](1))
}

func TestSet_NewEmpty(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, New[int64](1))
}

func TestSet_NewBy(t *testing.T) {
	t.Parallel()
	type input struct {
		key   string
		value int64
	}
	values := []*input{
		{key: "test01", value: 1},
		{key: "test02", value: 2},
		{key: "test03", value: 2},
		{key: "test04", value: 3},
	}
	iter := func(in *input) int64 {
		return in.value
	}
	assert.Len(t, UniqBy(values, iter), 3)
}

func TestSet_Uniq(t *testing.T) {
	t.Parallel()
	values := []int64{1, 2, 2, 3, 3, 3}
	assert.ElementsMatch(t, []int64{1, 2, 3}, Uniq(values...))
}

func TestSet_UniqBy(t *testing.T) {
	t.Parallel()
	type input struct {
		key   string
		value int64
	}
	values := []*input{
		{key: "test01", value: 1},
		{key: "test02", value: 2},
		{key: "test03", value: 2},
		{key: "test04", value: 3},
	}
	iter := func(in *input) int64 {
		return in.value
	}
	assert.ElementsMatch(t, []int64{1, 2, 3}, UniqBy(values, iter))
}

func TestSet_UniqWithErr(t *testing.T) {
	t.Parallel()
	type input struct {
		key   string
		value int64
	}
	values := []*input{
		{key: "test01", value: 1},
		{key: "test02", value: 2},
		{key: "test03", value: 2},
		{key: "test04", value: 3},
	}
	iter := func(in *input) (int64, error) {
		return in.value, nil
	}
	actual, err := UniqWithErr(values, iter)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int64{1, 2, 3}, actual)
	iter = func(in *input) (int64, error) {
		return 0, assert.AnError
	}
	actual, err = UniqWithErr(values, iter)
	assert.Error(t, err)
	assert.Equal(t, []int64(nil), actual)
}

func TestSet_Len(t *testing.T) {
	t.Parallel()
	set := New([]int64{1, 2, 2, 3, 3, 3}...)
	assert.Equal(t, 3, set.Len())
}

func TestSet_Reset(t *testing.T) {
	t.Parallel()
	set := New([]int64{1, 2, 2, 3, 3, 3}...)
	set.Reset(0)
	assert.ElementsMatch(t, []int64{}, set.Slice())
}

func TestSet_Contains(t *testing.T) {
	t.Parallel()
	set := New([]int64{1, 2, 2, 3, 3, 3}...)
	assert.True(t, set.Contains(1))
	assert.True(t, set.Contains(1, 2, 3))
	assert.False(t, set.Contains(4))
}

func TestSet_Add(t *testing.T) {
	t.Parallel()
	set := NewEmpty[int64](2)
	assert.ElementsMatch(t, []int64{}, set.Slice())
	set.Add(1, 2)
	assert.ElementsMatch(t, []int64{1, 2}, set.Slice())
}

func TestSet_FindOrAdd(t *testing.T) {
	t.Parallel()
	set := NewEmpty[int64](4)
	actual, exists := set.FindOrAdd(1)
	assert.False(t, exists)
	assert.ElementsMatch(t, []int64{1}, actual.Slice())
	actual, exists = set.FindOrAdd(1)
	assert.True(t, exists)
	assert.ElementsMatch(t, []int64{1}, actual.Slice())
	actual, exists = set.FindOrAdd(2, 3)
	assert.False(t, exists)
	assert.ElementsMatch(t, []int64{1, 2, 3}, actual.Slice())
	actual, exists = set.FindOrAdd(1, 4)
	assert.True(t, exists)
	assert.ElementsMatch(t, []int64{1, 2, 3, 4}, actual.Slice())
}

func TestSet_Remove(t *testing.T) {
	t.Parallel()
	set := New([]int64{1, 2, 3}...)
	set.Remove(1)
	assert.ElementsMatch(t, []int64{2, 3}, set.Slice())
}

func TestSet_Slice(t *testing.T) {
	t.Parallel()
	set := New([]string{"test01", "", "test02"}...)
	assert.ElementsMatch(t, []string{"test01", "test02"}, set.Slice())
}
