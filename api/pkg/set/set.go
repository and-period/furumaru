package set

import "sort"

// Set - 重複排除処理用の構造体
type Set struct {
	values map[interface{}]struct{}
}

// New - 構造体の生成(cap指定)
func New(cap int) *Set {
	return &Set{
		values: make(map[interface{}]struct{}, cap),
	}
}

// Set - 構造体の中身を初期化
func (s *Set) Reset(cap int) {
	s.values = make(map[interface{}]struct{}, cap)
}

// Len - 長さを取得
func (s *Set) Len() int {
	return len(s.values)
}

// Contains - 指定された値がすべて含まれているかの判定
func (s *Set) Contains(values ...interface{}) bool {
	for i := range values {
		if _, ok := s.values[values[i]]; !ok {
			return false
		}
	}
	return true
}

// FindOrAdd - 指定された値のいずれかが存在するかの判定 && 存在しないものは代入
func (s *Set) FindOrAdd(values ...interface{}) bool {
	isFind := false
	for i := range values {
		if _, ok := s.values[values[i]]; ok {
			continue
		}
		s.Add(values[i])
		isFind = true
	}
	return isFind
}

// Add - 指定された値を代入
func (s *Set) Add(values ...interface{}) {
	for _, v := range values {
		s.values[v] = struct{}{}
	}
}

// AddStrings - 指定された文字列を代入
func (s *Set) AddStrings(values ...string) {
	for _, v := range values {
		s.Add(v)
	}
}

// AddStrings - 指定された数字(int64)を代入
func (s *Set) AddInt64s(values ...int64) {
	for _, v := range values {
		s.Add(v)
	}
}

// Remove - 指定された値を削除
func (s *Set) Remove(v interface{}) {
	delete(s.values, v)
}

// Do - 構造体内の内のデータに対し、指定された処理を実行
func (s *Set) Do(f func(interface{})) {
	for v := range s.values {
		f(v)
	}
}

// Strings - 文字列型の配列として返す
func (s *Set) Strings() []string {
	res := make([]string, 0, s.Len())
	for v := range s.values {
		res = append(res, v.(string))
	}
	return res
}

// SortStrings - 文字列型の配列 (昇順でソートされた状態) として返す
func (s *Set) SortStrings() []string {
	res := s.Strings()
	sort.Strings(res)
	return res
}

// Int64s - int64型の配列として返す
func (s *Set) Int64s() []int64 {
	res := make([]int64, 0, s.Len())
	for v := range s.values {
		res = append(res, v.(int64))
	}
	return res
}

// SortInt64s - int64型の配列 (昇順でソートされた状態) として返す
func (s *Set) SortInt64s() []int64 {
	res := s.Int64s()
	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})
	return res
}
