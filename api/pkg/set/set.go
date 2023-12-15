package set

// Set - 重複排除処理用の構造体
type Set[T comparable] struct {
	values map[T]struct{}
}

// New - 渡された値を基に構造体の生成
func New[T comparable](values ...T) *Set[T] {
	return NewEmpty[T](len(values)).Add(values...)
}

// New - 構造体の生成(cap指定)
func NewEmpty[T comparable](cap int) *Set[T] {
	return &Set[T]{
		values: make(map[T]struct{}, cap),
	}
}

// NewBy - 渡された操作を実行して構造体の生成
func NewBy[T comparable, V any](values []V, iteratee func(V) T) *Set[T] {
	set := NewEmpty[T](len(values))
	for i := range values {
		key := iteratee(values[i])
		set.Add(key)
	}
	return set
}

// Uniq - 渡された値の重複を排除して返す
func Uniq[T comparable](values ...T) []T {
	return New(values...).Slice()
}

// UniqBy - 渡された操作を実行して重複を排除した値を返す
func UniqBy[K comparable, V any](values []V, iteratee func(V) K) []K {
	set := NewEmpty[K](len(values))
	for i := range values {
		key := iteratee(values[i])
		set.Add(key)
	}
	return set.Slice()
}

// UniqWithErr - 渡された操作を実行して重複を排除した値を返す
func UniqWithErr[K comparable, V any](values []V, iteratee func(V) (K, error)) ([]K, error) {
	set := NewEmpty[K](len(values))
	for i := range values {
		key, err := iteratee(values[i])
		if err != nil {
			return nil, err
		}
		set.Add(key)
	}
	return set.Slice(), nil
}

// Len - 長さを取得
func (s *Set[T]) Len() int {
	return len(s.values)
}

// Reset - 構造体の中身を初期化
func (s *Set[T]) Reset(cap int) *Set[T] {
	s.values = make(map[T]struct{}, cap)
	return s
}

// Contains - 指定された値がすべて含まれているかの検証
func (s *Set[T]) Contains(values ...T) bool {
	for i := range values {
		if _, ok := s.values[values[i]]; ok {
			continue
		}
		return false
	}
	return true
}

// Add - 指定された値を代入
func (s *Set[T]) Add(values ...T) *Set[T] {
	for _, v := range values {
		s.values[v] = struct{}{}
	}
	return s
}

// FindOrAdd - 指定された値のいずれかが存在するかの判定 && 存在しないものは代入
func (s *Set[T]) FindOrAdd(values ...T) (*Set[T], bool) {
	var isExists bool
	for i := range values {
		if s.Contains(values[i]) {
			isExists = true
			continue
		}
		s.Add(values[i])
	}
	return s, isExists
}

// Remove - 指定された値を削除
func (s *Set[T]) Remove(v T) *Set[T] {
	delete(s.values, v)
	return s
}

// Slice - Sliceとして返す
func (s *Set[T]) Slice() []T {
	var zero T
	res := make([]T, 0, s.Len())
	for v := range s.values {
		if v == zero {
			continue
		}
		res = append(res, v)
	}
	return res
}
