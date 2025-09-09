package service

import (
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/shopspring/decimal"
)

// TopOrderPeriodType - 期間種別
type TopOrderPeriodType string

const (
	TopOrderPeriodTypeDay   TopOrderPeriodType = "day"   // 日
	TopOrderPeriodTypeWeek  TopOrderPeriodType = "week"  // 週
	TopOrderPeriodTypeMonth TopOrderPeriodType = "month" // 月
)

func NewTopOrderPeriodType(periodType entity.AggregateOrderPeriodType) TopOrderPeriodType {
	switch periodType {
	case entity.AggregateOrderPeriodTypeDay:
		return TopOrderPeriodTypeDay
	case entity.AggregateOrderPeriodTypeWeek:
		return TopOrderPeriodTypeWeek
	case entity.AggregateOrderPeriodTypeMonth:
		return TopOrderPeriodTypeMonth
	default:
		return ""
	}
}

func NewTopOrderPeriodTypeFromRequest(periodType string) TopOrderPeriodType {
	return TopOrderPeriodType(periodType)
}

func (t TopOrderPeriodType) String(period time.Time) string {
	switch t {
	case TopOrderPeriodTypeDay:
		return period.Format(time.DateOnly)
	case TopOrderPeriodTypeWeek:
		days := int(period.Weekday())
		week := jst.Date(period.Year(), period.Month(), period.Day()-days, 0, 0, 0, 0)
		return week.Format(time.DateOnly)
	case TopOrderPeriodTypeMonth:
		month := jst.Date(period.Year(), period.Month(), 1, 0, 0, 0, 0)
		return month.Format(time.DateOnly)
	default:
		return ""
	}
}

func (t TopOrderPeriodType) Add(ts time.Time) time.Time {
	switch t {
	case TopOrderPeriodTypeDay:
		return ts.AddDate(0, 0, 1)
	case TopOrderPeriodTypeWeek:
		return ts.AddDate(0, 0, 7)
	case TopOrderPeriodTypeMonth:
		return ts.AddDate(0, 1, 0)
	default:
		return ts
	}
}

func (t TopOrderPeriodType) Truncate(ts time.Time) time.Time {
	switch t {
	case TopOrderPeriodTypeDay:
		return jst.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0)
	case TopOrderPeriodTypeWeek:
		days := int(ts.Weekday())
		return jst.Date(ts.Year(), ts.Month(), ts.Day()-days, 0, 0, 0, 0)
	case TopOrderPeriodTypeMonth:
		return jst.Date(ts.Year(), ts.Month(), 1, 0, 0, 0, 0)
	default:
		return ts
	}
}

func (t TopOrderPeriodType) StoreEntity() entity.AggregateOrderPeriodType {
	switch t {
	case TopOrderPeriodTypeDay:
		return entity.AggregateOrderPeriodTypeDay
	case TopOrderPeriodTypeWeek:
		return entity.AggregateOrderPeriodTypeWeek
	case TopOrderPeriodTypeMonth:
		return entity.AggregateOrderPeriodTypeMonth
	default:
		return ""
	}
}

func (t TopOrderPeriodType) Response() string {
	switch t {
	case TopOrderPeriodTypeDay:
		return "day"
	case TopOrderPeriodTypeWeek:
		return "week"
	case TopOrderPeriodTypeMonth:
		return "month"
	default:
		return ""
	}
}

type TopOrderValue struct {
	types.TopOrderValue
}

func NewTopOrderValue(current, before int64) *TopOrderValue {
	if current == 0 && before == 0 {
		return newEmptyTopOrderValue()
	}

	var comparison float64
	if before == 0 {
		comparison = 100 // 0除算を避けるため
	} else {
		dcurrent := decimal.NewFromInt(current)
		dbefore := decimal.NewFromInt(before)
		diff := dcurrent.Sub(dbefore)
		dcomp := diff.Div(dbefore).Mul(decimal.NewFromInt(100))
		comparison, _ = dcomp.Float64()
	}

	return &TopOrderValue{
		TopOrderValue: types.TopOrderValue{
			Value:      current,
			Comparison: comparison,
		},
	}
}

func newEmptyTopOrderValue() *TopOrderValue {
	return &TopOrderValue{
		TopOrderValue: types.TopOrderValue{
			Value:      0,
			Comparison: 0,
		},
	}
}

func (v *TopOrderValue) Response() *types.TopOrderValue {
	return &v.TopOrderValue
}

type TopOrderSalesTrend struct {
	types.TopOrderSalesTrend
}

type TopOrderSalesTrends []*TopOrderSalesTrend

func NewTopOrderSalesTrend(aggregated *entity.AggregatedPeriodOrder, periodType TopOrderPeriodType) *TopOrderSalesTrend {
	return &TopOrderSalesTrend{
		TopOrderSalesTrend: types.TopOrderSalesTrend{
			Period:     periodType.String(aggregated.Period),
			SalesTotal: aggregated.SalesTotal,
		},
	}
}

func newEmptyTopOrderSalesTrend(ts time.Time, periodType TopOrderPeriodType) *TopOrderSalesTrend {
	return &TopOrderSalesTrend{
		TopOrderSalesTrend: types.TopOrderSalesTrend{
			Period:     periodType.String(ts),
			SalesTotal: 0,
		},
	}
}

func (t *TopOrderSalesTrend) Response() *types.TopOrderSalesTrend {
	return &t.TopOrderSalesTrend
}

func NewTopOrderSalesTrends(
	periodType TopOrderPeriodType,
	startAt, endAt time.Time,
	aggregated entity.AggregatedPeriodOrders,
) TopOrderSalesTrends {
	aggregatedMap := aggregated.MapByPeriod()
	res := make(TopOrderSalesTrends, 0, len(aggregated))
	for ts := periodType.Truncate(startAt); ts.Before(endAt); ts = periodType.Add(ts) {
		if aggregate, ok := aggregatedMap[ts]; ok {
			res = append(res, NewTopOrderSalesTrend(aggregate, periodType))
			continue
		}
		res = append(res, newEmptyTopOrderSalesTrend(ts, periodType))
	}
	return res
}

func (t TopOrderSalesTrends) Response() []*types.TopOrderSalesTrend {
	res := make([]*types.TopOrderSalesTrend, len(t))
	for i := range t {
		res[i] = t[i].Response()
	}
	return res
}

type TopOrderPayment struct {
	types.TopOrderPayment
}

type TopOrderPayments []*TopOrderPayment

func NewTopOrderPayment(payment *entity.AggregatedOrderPayment, total int64) *TopOrderPayment {
	var rate float64
	if total > 0 {
		dcount := decimal.NewFromInt(payment.OrderCount)
		dtotal := decimal.NewFromInt(total)
		rate, _ = dcount.Div(dtotal).Mul(decimal.NewFromInt(100)).Float64()
	}

	return &TopOrderPayment{
		TopOrderPayment: types.TopOrderPayment{
			PaymentMethodType: NewPaymentMethodType(payment.PaymentMethodType).Response(),
			OrderCount:        payment.OrderCount,
			UserCount:         payment.UserCount,
			SalesTotal:        payment.SalesTotal,
			Rate:              rate,
		},
	}
}

func (p *TopOrderPayment) Response() *types.TopOrderPayment {
	return &p.TopOrderPayment
}

func NewTopOrderPayments(payments entity.AggregatedOrderPayments) TopOrderPayments {
	total := payments.OrderTotal()
	res := make(TopOrderPayments, len(payments))
	for i := range payments {
		res[i] = NewTopOrderPayment(payments[i], total)
	}
	return res
}

func (p TopOrderPayments) Response() []*types.TopOrderPayment {
	res := make([]*types.TopOrderPayment, len(p))
	for i := range p {
		res[i] = p[i].Response()
	}
	return res
}
