package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
		assert.Nil(t, NullInt[int](0))
	})
	t.Run("null int32", func(t *testing.T) {
		t.Parallel()
		assert.Nil(t, NullInt[int32](0))
	})
	t.Run("null int64", func(t *testing.T) {
		t.Parallel()
		assert.Nil(t, NullInt[int64](0))
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
		assert.Nil(t, NullString(""))
	})
}

func TestGeometry_Value(t *testing.T) {
	t.Parallel()
	t.Run("null geometry", func(t *testing.T) {
		t.Parallel()
		geometry := Geometry{}
		value, err := geometry.Value()
		assert.NoError(t, err)
		assert.Equal(t, "POINT(0.000000 0.000000)", value)
	})
	t.Run("non null geometry", func(t *testing.T) {
		t.Parallel()
		geometry := Geometry{X: 1, Y: 1}
		value, err := geometry.Value()
		assert.NoError(t, err)
		assert.Equal(t, "POINT(1.000000 1.000000)", value)
	})
}

func TestGeometry_Scan(t *testing.T) {
	t.Parallel()
	t.Run("null geometry", func(t *testing.T) {
		t.Parallel()
		geometry := Geometry{}
		err := geometry.Scan(nil)
		assert.NoError(t, err)
		assert.Equal(t, Geometry{}, geometry)
	})
	t.Run("valid byte type", func(t *testing.T) {
		t.Parallel()
		geometry := Geometry{}

		err := geometry.Scan(
			[]byte{
				0,
				0,
				0,
				0,
				1,
				1,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				240,
				63,
				0,
				0,
				0,
				0,
				0,
				0,
				240,
				191,
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, Geometry{X: 1, Y: -1}, geometry)
	})
	t.Run("invalid byte format", func(t *testing.T) {
		t.Parallel()
		geometry := Geometry{}

		err := geometry.Scan(
			[]byte{0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0},
		)
		assert.Error(t, err)
		assert.Equal(t, Geometry{}, geometry)
	})
	t.Run("valid string type", func(t *testing.T) {
		t.Parallel()
		geometry := Geometry{}

		err := geometry.Scan("POINT(1.000000 -1.000000)")
		assert.NoError(t, err)
		assert.Equal(t, Geometry{X: 1, Y: -1}, geometry)
	})
	t.Run("invalid string format", func(t *testing.T) {
		t.Parallel()
		geometry := Geometry{}

		err := geometry.Scan("POINT(1.000000 -1.000000")
		assert.Error(t, err)
		assert.Equal(t, Geometry{}, geometry)
	})
	t.Run("invalid type", func(t *testing.T) {
		t.Parallel()
		geometry := Geometry{}

		err := geometry.Scan(1)
		assert.Error(t, err)
		assert.Equal(t, Geometry{}, geometry)
	})
}

func TestGeometry_GormDataType(t *testing.T) {
	t.Parallel()
	t.Run("geometry", func(t *testing.T) {
		t.Parallel()
		geometry := Geometry{X: 1, Y: 1}
		dataType := geometry.GormDataType()
		assert.Equal(t, "geometry", dataType)
	})
}

func TestGeometry_GormDBDataType(t *testing.T) {
	t.Parallel()
	t.Run("mysql", func(t *testing.T) {
		t.Parallel()
		dialector := mysql.New(mysql.Config{})
		db := &gorm.DB{Config: &gorm.Config{Dialector: dialector}}
		field := &schema.Field{}

		geometry := Geometry{X: 1, Y: 1}
		dbDataType := geometry.GormDBDataType(db, field)
		assert.Equal(t, "GEOMETRY", dbDataType)
	})
}

func TestGeometry_GormValue(t *testing.T) {
	t.Parallel()
	t.Run("mysql", func(t *testing.T) {
		t.Parallel()
		dialector := mysql.New(mysql.Config{})
		db := &gorm.DB{Config: &gorm.Config{Dialector: dialector}}

		geometry := Geometry{X: 1, Y: 1}
		expr := geometry.GormValue(t.Context(), db)
		assert.Equal(t, "ST_GeomFromText(?)", expr.SQL)
		assert.Equal(t, []interface{}{"POINT(1.000000 1.000000)"}, expr.Vars)
	})
}
