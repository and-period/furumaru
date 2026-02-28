package mysql

import (
	"database/sql/driver"
	"encoding/base64"
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
// string型で返すことで、ドライバが[]byteをバイナリ(base64)として扱うのを防ぐ
func (j JSONColumn[T]) Value() (driver.Value, error) {
	b, err := json.Marshal(j.Val)
	if err != nil {
		return nil, err
	}
	return string(b), nil
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
	if err := json.Unmarshal(bytes, &j.Val); err != nil {
		// TiDBのJSONカラムがbase64エンコードされた文字列として返される場合のフォールバック
		// 例: "WyJ0YWcxIiwidGFnMiJd" (base64) → ["tag1","tag2"] (JSON)
		var s string
		if jsonErr := json.Unmarshal(bytes, &s); jsonErr == nil {
			if decoded, decErr := base64.StdEncoding.DecodeString(s); decErr == nil {
				return json.Unmarshal(decoded, &j.Val)
			}
		}
		return err
	}
	return nil
}
