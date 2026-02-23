package entity

import (
	"iter"
	"time"
)

// All はインデックスと注文履歴のペアを返すイテレーターを返す。
func (os Orders) All() iter.Seq2[int, *Order] {
	return func(yield func(int, *Order) bool) {
		for i, o := range os {
			if !yield(i, o) {
				return
			}
		}
	}
}

// IterMap は注文履歴IDをキー、注文履歴を値とするイテレーターを返す。
func (os Orders) IterMap() iter.Seq2[string, *Order] {
	return MapIter(os, func(o *Order) (string, *Order) {
		return o.ID, o
	})
}

// All はインデックスと注文履歴集計情報のペアを返すイテレーターを返す。
func (os AggregatedUserOrders) All() iter.Seq2[int, *AggregatedUserOrder] {
	return func(yield func(int, *AggregatedUserOrder) bool) {
		for i, o := range os {
			if !yield(i, o) {
				return
			}
		}
	}
}

// IterMap はユーザーIDをキー、注文履歴集計情報を値とするイテレーターを返す。
func (os AggregatedUserOrders) IterMap() iter.Seq2[string, *AggregatedUserOrder] {
	return MapIter(os, func(o *AggregatedUserOrder) (string, *AggregatedUserOrder) {
		return o.UserID, o
	})
}

// All はインデックスと支払い情報別集計情報のペアを返すイテレーターを返す。
func (ps AggregatedOrderPayments) All() iter.Seq2[int, *AggregatedOrderPayment] {
	return func(yield func(int, *AggregatedOrderPayment) bool) {
		for i, p := range ps {
			if !yield(i, p) {
				return
			}
		}
	}
}

// IterMap は支払い種別をキー、支払い情報別集計情報を値とするイテレーターを返す。
func (ps AggregatedOrderPayments) IterMap() iter.Seq2[PaymentMethodType, *AggregatedOrderPayment] {
	return MapIter(ps, func(p *AggregatedOrderPayment) (PaymentMethodType, *AggregatedOrderPayment) {
		return p.PaymentMethodType, p
	})
}

// All はインデックスとプロモーション利用履歴集計情報のペアを返すイテレーターを返す。
func (os AggregatedOrderPromotions) All() iter.Seq2[int, *AggregatedOrderPromotion] {
	return func(yield func(int, *AggregatedOrderPromotion) bool) {
		for i, o := range os {
			if !yield(i, o) {
				return
			}
		}
	}
}

// IterMap はプロモーションIDをキー、プロモーション利用履歴集計情報を値とするイテレーターを返す。
func (os AggregatedOrderPromotions) IterMap() iter.Seq2[string, *AggregatedOrderPromotion] {
	return MapIter(os, func(o *AggregatedOrderPromotion) (string, *AggregatedOrderPromotion) {
		return o.PromotionID, o
	})
}

// All はインデックスと期間別注文履歴集計情報のペアを返すイテレーターを返す。
func (os AggregatedPeriodOrders) All() iter.Seq2[int, *AggregatedPeriodOrder] {
	return func(yield func(int, *AggregatedPeriodOrder) bool) {
		for i, o := range os {
			if !yield(i, o) {
				return
			}
		}
	}
}

// IterMapByPeriod は期間をキー、期間別注文履歴集計情報を値とするイテレーターを返す。
func (os AggregatedPeriodOrders) IterMapByPeriod() iter.Seq2[time.Time, *AggregatedPeriodOrder] {
	return MapIter(os, func(o *AggregatedPeriodOrder) (time.Time, *AggregatedPeriodOrder) {
		return o.Period, o
	})
}
