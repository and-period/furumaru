package mysql

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/binary"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Integer interface {
	~int | ~int32 | ~int64
}

// NullInt ゼロ値の場合にnilを返す関数
func NullInt[T Integer](val T) any {
	if val == 0 {
		return nil
	}
	return val
}

// NullString ゼロ値の場合にnilを返す関数
func NullString(val string) any {
	if val == "" {
		return nil
	}
	return val
}

// Geometry 座標情報
//
//nolint:recvcheck
type Geometry struct {
	X float64
	Y float64
}

func (g Geometry) Value() (driver.Value, error) {
	if g.X == 0 && g.Y == 0 {
		return "POINT(0.000000 0.000000)", nil
	}
	return fmt.Sprintf("POINT(%f %f)", g.X, g.Y), nil
}

func (g *Geometry) Scan(value interface{}) error {
	if value == nil {
		*g = Geometry{}
		return nil
	}

	var longitude, latitude float64
	// @see: https://dev.mysql.com/doc/refman/8.0/en/gis-data-formats.html
	switch v := value.(type) {
	case []byte:
		if len(v) != 25 {
			return fmt.Errorf("query: expected []bytes with length 25, got %d", len(v))
		}
		buf := bytes.NewReader(v[9:17])
		if err := binary.Read(buf, binary.LittleEndian, &longitude); err != nil {
			return err
		}
		buf = bytes.NewReader(v[17:])
		if err := binary.Read(buf, binary.LittleEndian, &latitude); err != nil {
			return err
		}
	case string:
		if _, err := fmt.Sscanf(v, "POINT(%f %f)", &longitude, &latitude); err != nil {
			return err
		}
	default:
		return fmt.Errorf("query: unsupported Geometry type: %T", value)
	}

	*g = Geometry{X: longitude, Y: latitude}
	return nil
}

func (g Geometry) GormDataType() string {
	return "geometry"
}

func (g Geometry) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Name() {
	case "mysql":
		return "GEOMETRY"
	default:
		return ""
	}
}

func (g Geometry) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := g.Value()

	switch db.Name() {
	case "mysql":
		return gorm.Expr("ST_GeomFromText(?)", data)
	default:
		return gorm.Expr("?", data)
	}
}
