package mysql

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONColumn_Value(t *testing.T) {
	t.Parallel()

	t.Run("string slice", func(t *testing.T) {
		t.Parallel()
		col := NewJSONColumn([]string{"a", "b", "c"})
		val, err := col.Value()
		require.NoError(t, err)
		assert.JSONEq(t, `["a","b","c"]`, val.(string))
	})

	t.Run("int slice", func(t *testing.T) {
		t.Parallel()
		col := NewJSONColumn([]int{1, 2, 3})
		val, err := col.Value()
		require.NoError(t, err)
		assert.JSONEq(t, `[1,2,3]`, val.(string))
	})

	t.Run("struct slice", func(t *testing.T) {
		t.Parallel()
		type item struct {
			Name  string `json:"name"`
			Price int    `json:"price"`
		}
		col := NewJSONColumn([]item{
			{Name: "apple", Price: 100},
			{Name: "banana", Price: 200},
		})
		val, err := col.Value()
		require.NoError(t, err)
		assert.JSONEq(t, `[{"name":"apple","price":100},{"name":"banana","price":200}]`, val.(string))
	})

	t.Run("map", func(t *testing.T) {
		t.Parallel()
		col := NewJSONColumn(map[string]int{"x": 1, "y": 2})
		val, err := col.Value()
		require.NoError(t, err)
		assert.JSONEq(t, `{"x":1,"y":2}`, val.(string))
	})

	t.Run("nil slice", func(t *testing.T) {
		t.Parallel()
		col := NewJSONColumn[[]string](nil)
		val, err := col.Value()
		require.NoError(t, err)
		assert.Equal(t, "null", val.(string))
	})
}

func TestJSONColumn_Scan(t *testing.T) {
	t.Parallel()

	t.Run("from bytes", func(t *testing.T) {
		t.Parallel()
		col := &JSONColumn[[]string]{}
		err := col.Scan([]byte(`["a","b","c"]`))
		require.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, col.Val)
	})

	t.Run("from string", func(t *testing.T) {
		t.Parallel()
		col := &JSONColumn[[]string]{}
		err := col.Scan(`["x","y"]`)
		require.NoError(t, err)
		assert.Equal(t, []string{"x", "y"}, col.Val)
	})

	t.Run("from nil", func(t *testing.T) {
		t.Parallel()
		col := &JSONColumn[[]string]{}
		err := col.Scan(nil)
		require.NoError(t, err)
		assert.Nil(t, col.Val)
	})

	t.Run("struct slice from bytes", func(t *testing.T) {
		t.Parallel()
		type item struct {
			Name  string `json:"name"`
			Price int    `json:"price"`
		}
		col := &JSONColumn[[]item]{}
		err := col.Scan([]byte(`[{"name":"apple","price":100}]`))
		require.NoError(t, err)
		require.Len(t, col.Val, 1)
		assert.Equal(t, "apple", col.Val[0].Name)
		assert.Equal(t, 100, col.Val[0].Price)
	})

	t.Run("map from bytes", func(t *testing.T) {
		t.Parallel()
		col := &JSONColumn[map[string]int]{}
		err := col.Scan([]byte(`{"x":1,"y":2}`))
		require.NoError(t, err)
		assert.Equal(t, map[string]int{"x": 1, "y": 2}, col.Val)
	})

	t.Run("unsupported type", func(t *testing.T) {
		t.Parallel()
		col := &JSONColumn[[]string]{}
		err := col.Scan(12345)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported type for JSONColumn")
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()
		col := &JSONColumn[[]string]{}
		err := col.Scan([]byte(`invalid-json`))
		assert.Error(t, err)
	})

	t.Run("base64 encoded string slice", func(t *testing.T) {
		t.Parallel()
		// base64("["a","b","c"]") → "WyJhIiwiYiIsImMiXQ=="
		encoded := base64.StdEncoding.EncodeToString([]byte(`["a","b","c"]`))
		col := &JSONColumn[[]string]{}
		err := col.Scan([]byte(`"` + encoded + `"`))
		require.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, col.Val)
	})

	t.Run("base64 encoded struct slice", func(t *testing.T) {
		t.Parallel()
		type item struct {
			Name  string `json:"name"`
			Price int    `json:"price"`
		}
		encoded := base64.StdEncoding.EncodeToString([]byte(`[{"name":"apple","price":100}]`))
		col := &JSONColumn[[]item]{}
		err := col.Scan([]byte(`"` + encoded + `"`))
		require.NoError(t, err)
		require.Len(t, col.Val, 1)
		assert.Equal(t, "apple", col.Val[0].Name)
		assert.Equal(t, 100, col.Val[0].Price)
	})

	t.Run("base64 encoded map", func(t *testing.T) {
		t.Parallel()
		encoded := base64.StdEncoding.EncodeToString([]byte(`{"x":1,"y":2}`))
		col := &JSONColumn[map[string]int]{}
		err := col.Scan([]byte(`"` + encoded + `"`))
		require.NoError(t, err)
		assert.Equal(t, map[string]int{"x": 1, "y": 2}, col.Val)
	})
}

func TestNewJSONColumn(t *testing.T) {
	t.Parallel()

	col := NewJSONColumn([]int32{1, 2, 3})
	assert.Equal(t, []int32{1, 2, 3}, col.Val)
}
