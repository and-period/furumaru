package entity

import "iter"

// All はインデックスとお知らせのペアを返すイテレーターを返す。
func (ns Notifications) All() iter.Seq2[int, *Notification] {
	return func(yield func(int, *Notification) bool) {
		for i, n := range ns {
			if !yield(i, n) {
				return
			}
		}
	}
}

// IterPromotionIDs はプロモーション種別のお知らせからプロモーションIDを返すイテレーターを返す。
// NotificationTypePromotion でないお知らせはスキップされる。
func (ns Notifications) IterPromotionIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, n := range ns {
			if n.Type != NotificationTypePromotion {
				continue
			}
			if !yield(n.PromotionID) {
				return
			}
		}
	}
}

// IterAdminIDs はお知らせの作成者IDと更新者IDを返すイテレーターを返す。
func (ns Notifications) IterAdminIDs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, n := range ns {
			if !yield(n.CreatedBy) {
				return
			}
			if !yield(n.UpdatedBy) {
				return
			}
		}
	}
}
