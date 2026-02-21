package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONColumn - GORM用のJSON型カラムのジェネリック実装
// database/sql/driver.Valuer と sql.Scanner を実装し、
// 任意の型をJSONとしてDBに読み書きできる
type JSONColumn[T any] struct {
	Val T
}

// NewJSONColumn - JSONColumn のコンストラクタ
func NewJSONColumn[T any](val T) JSONColumn[T] {
	return JSONColumn[T]{Val: val}
}

// Value - driver.Valuer の実装。DBへの書き込み時にJSON文字列に変換する
func (j JSONColumn[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Val)
}

// Scan - sql.Scanner の実装。DBからの読み込み時にJSONをデシリアライズする
func (j *JSONColumn[T]) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("mysql: unsupported type for JSONColumn: %T", value)
	}
	return json.Unmarshal(bytes, &j.Val)
}
