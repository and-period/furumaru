package entity

import "iter"

// All はインデックスとユーザーのペアを返すイテレーターを返す。
func (us Users) All() iter.Seq2[int, *User] {
	return func(yield func(int, *User) bool) {
		for i, u := range us {
			if !yield(i, u) {
				return
			}
		}
	}
}

// IterMap はユーザーIDをキー、ユーザーを値とするイテレーターを返す。
func (us Users) IterMap() iter.Seq2[string, *User] {
	return MapIter(us, func(u *User) (string, *User) {
		return u.ID, u
	})
}

// IterGroupByRegistered は会員登録フラグをキー、ユーザーを値とするイテレーターを返す。
func (us Users) IterGroupByRegistered() iter.Seq2[bool, *User] {
	return MapIter(us, func(u *User) (bool, *User) {
		return u.Registered, u
	})
}

// IterGroupByUserType はユーザー種別をキー、ユーザーを値とするイテレーターを返す。
func (us Users) IterGroupByUserType() iter.Seq2[UserType, *User] {
	return MapIter(us, func(u *User) (UserType, *User) {
		return u.Type, u
	})
}
