package entity

import "iter"

// FilterIter は述語関数に一致する要素を返すイテレーターを返す。
func FilterIter[T any](s []T, fn func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s {
			if fn(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// MapIter はスライスの各要素をキーと値のペアに変換するイテレーターを返す。
func MapIter[T any, K comparable, V any](s []T, fn func(T) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, v := range s {
			k, val := fn(v)
			if !yield(k, val) {
				return
			}
		}
	}
}
