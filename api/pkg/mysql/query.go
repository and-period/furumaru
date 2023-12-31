package mysql

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
