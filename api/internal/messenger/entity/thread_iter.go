package entity

import "iter"

// All はインデックスとスレッドのペアを返すイテレーターを返す。
func (ts Threads) All() iter.Seq2[int, *Thread] {
	return func(yield func(int, *Thread) bool) {
		for i, t := range ts {
			if !yield(i, t) {
				return
			}
		}
	}
}

// IterIDs はスレッドIDを返すイテレーターを返す。
func (ts Threads) IterIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, t := range ts {
			if !yield(t.ID) {
				return
			}
		}
	}
}

// IterUserIDs はユーザー種別がユーザーのスレッドからユーザーIDを返すイテレーターを返す。
// ThreadUserTypeUser でないスレッドはスキップされる。
func (ts Threads) IterUserIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, t := range ts {
			if t.UserType != ThreadUserTypeUser {
				continue
			}
			if !yield(t.UserID) {
				return
			}
		}
	}
}

// IterAdminIDs はユーザー種別が管理者のスレッドからユーザーIDを返すイテレーターを返す。
// ThreadUserTypeAdmin でないスレッドはスキップされる。
func (ts Threads) IterAdminIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, t := range ts {
			if t.UserType != ThreadUserTypeAdmin {
				continue
			}
			if !yield(t.UserID) {
				return
			}
		}
	}
}
