package entity

import "iter"

// All はインデックスとお問い合わせのペアを返すイテレーターを返す。
func (cs Contacts) All() iter.Seq2[int, *Contact] {
	return func(yield func(int, *Contact) bool) {
		for i, c := range cs {
			if !yield(i, c) {
				return
			}
		}
	}
}

// IterIDs はお問い合わせIDを返すイテレーターを返す。
func (cs Contacts) IterIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, c := range cs {
			if !yield(c.ID) {
				return
			}
		}
	}
}

// IterCategoryIDs はお問い合わせ種別IDを返すイテレーターを返す。
func (cs Contacts) IterCategoryIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, c := range cs {
			if !yield(c.CategoryID) {
				return
			}
		}
	}
}

// IterUserIDs はユーザーIDを返すイテレーターを返す。
func (cs Contacts) IterUserIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, c := range cs {
			if !yield(c.UserID) {
				return
			}
		}
	}
}

// IterResponderIDs は対応者IDを返すイテレーターを返す。
func (cs Contacts) IterResponderIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, c := range cs {
			if !yield(c.ResponderID) {
				return
			}
		}
	}
}
