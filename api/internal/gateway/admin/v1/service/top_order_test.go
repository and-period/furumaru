package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestTopOrderPeriodType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		periodType entity.AggregateOrderPeriodType
		expect     TopOrderPeriodType
		response   string
	}{
		{
			name:       "day",
			periodType: entity.AggregateOrderPeriodTypeDay,
			expect:     TopOrderPeriodTypeDay,
			response:   "day",
		},
		{
			name:       "week",
			periodType: entity.AggregateOrderPeriodTypeWeek,
			expect:     TopOrderPeriodTypeWeek,
			response:   "week",
		},
		{
			name:       "month",
			periodType: entity.AggregateOrderPeriodTypeMonth,
			expect:     TopOrderPeriodTypeMonth,
			response:   "month",
		},
		{
			name:       "default",
			periodType: entity.AggregateOrderPeriodType(""),
			expect:     "",
			response:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewTopOrderPeriodType(tt.periodType)
			assert.Equal(t, tt.expect, actual)

			entity := actual.StoreEntity()
			assert.Equal(t, tt.periodType, entity)

			res := actual.Response()
			assert.Equal(t, tt.response, res)
		})
	}
}

func TestTopOrderPeriodTypeFromRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		periodType string
		expect     TopOrderPeriodType
	}{
		{
			name:       "day",
			periodType: "day",
			expect:     TopOrderPeriodTypeDay,
		},
		{
			name:       "week",
			periodType: "week",
			expect:     TopOrderPeriodTypeWeek,
		},
		{
			name:       "month",
			periodType: "month",
			expect:     TopOrderPeriodTypeMonth,
		},
		{
			name:       "default",
			periodType: "",
			expect:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewTopOrderPeriodTypeFromRequest(tt.periodType)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestTopOrderPeriodType_String(t *testing.T) {
	t.Parallel()
	now := jst.Date(2025, 1, 18, 18, 30, 0, 0)
	tests := []struct {
		name   string
		period TopOrderPeriodType
		expect string
	}{
		{
			name:   "day",
			period: TopOrderPeriodTypeDay,
			expect: "2025-01-18",
		},
		{
			name:   "week",
			period: TopOrderPeriodTypeWeek,
			expect: "2025-01-12",
		},
		{
			name:   "month",
			period: TopOrderPeriodTypeMonth,
			expect: "2025-01-01",
		},
		{
			name:   "default",
			period: "",
			expect: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.period.String(now)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestTopOrderValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		current int64
		before  int64
		expect  *TopOrderValue
	}{
		{
			name:    "default",
			current: 0,
			before:  0,
			expect: &TopOrderValue{
				TopOrderValue: response.TopOrderValue{
					Value:      0,
					Comparison: 0,
				},
			},
		},
		{
			name:    "current > before",
			current: 100,
			before:  50,
			expect: &TopOrderValue{
				TopOrderValue: response.TopOrderValue{
					Value:      100,
					Comparison: 100,
				},
			},
		},
		{
			name:    "current < before",
			current: 50,
			before:  100,
			expect: &TopOrderValue{
				TopOrderValue: response.TopOrderValue{
					Value:      50,
					Comparison: -50,
				},
			},
		},
		{
			name:    "current = before",
			current: 100,
			before:  100,
			expect: &TopOrderValue{
				TopOrderValue: response.TopOrderValue{
					Value:      100,
					Comparison: 0,
				},
			},
		},
		{
			name:    "before = 0",
			current: 50,
			before:  0,
			expect: &TopOrderValue{
				TopOrderValue: response.TopOrderValue{
					Value:      50,
					Comparison: 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewTopOrderValue(tt.current, tt.before)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestTopOrderValue_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		value  *TopOrderValue
		expect *response.TopOrderValue
	}{
		{
			name: "success",
			value: &TopOrderValue{
				TopOrderValue: response.TopOrderValue{
					Value:      100,
					Comparison: 50,
				},
			},
			expect: &response.TopOrderValue{
				Value:      100,
				Comparison: 50,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.value.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestTopOrderSalesTrends(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		periodType TopOrderPeriodType
		startAt    time.Time
		endAt      time.Time
		aggregated entity.AggregatedPeriodOrders
		expect     TopOrderSalesTrends
	}{
		{
			name:       "success day",
			periodType: TopOrderPeriodTypeDay,
			startAt:    jst.Date(2025, 1, 17, 18, 30, 0, 0),
			endAt:      jst.Date(2025, 1, 18, 18, 30, 0, 0),
			aggregated: entity.AggregatedPeriodOrders{
				{
					Period:        jst.Date(2025, 1, 18, 0, 0, 0, 0),
					OrderCount:    2,
					UserCount:     1,
					SalesTotal:    100,
					DiscountTotal: 0,
				},
			},
			expect: TopOrderSalesTrends{
				{
					TopOrderSalesTrend: response.TopOrderSalesTrend{
						Period:     "2025-01-17",
						SalesTotal: 0,
					},
				},
				{
					TopOrderSalesTrend: response.TopOrderSalesTrend{
						Period:     "2025-01-18",
						SalesTotal: 100,
					},
				},
			},
		},
		{
			name:       "success week",
			periodType: TopOrderPeriodTypeWeek,
			startAt:    jst.Date(2025, 1, 1, 18, 30, 0, 0),
			endAt:      jst.Date(2025, 1, 18, 18, 30, 0, 0),
			aggregated: entity.AggregatedPeriodOrders{
				{
					Period:        jst.Date(2025, 1, 12, 0, 0, 0, 0),
					OrderCount:    2,
					UserCount:     1,
					SalesTotal:    100,
					DiscountTotal: 0,
				},
			},
			expect: TopOrderSalesTrends{
				{
					TopOrderSalesTrend: response.TopOrderSalesTrend{
						Period:     "2024-12-29",
						SalesTotal: 0,
					},
				},
				{
					TopOrderSalesTrend: response.TopOrderSalesTrend{
						Period:     "2025-01-05",
						SalesTotal: 0,
					},
				},
				{
					TopOrderSalesTrend: response.TopOrderSalesTrend{
						Period:     "2025-01-12",
						SalesTotal: 100,
					},
				},
			},
		},
		{
			name:       "success month",
			periodType: TopOrderPeriodTypeMonth,
			startAt:    jst.Date(2025, 1, 1, 18, 30, 0, 0),
			endAt:      jst.Date(2025, 1, 18, 18, 30, 0, 0),
			aggregated: entity.AggregatedPeriodOrders{
				{
					Period:        jst.Date(2025, 1, 1, 0, 0, 0, 0),
					OrderCount:    2,
					UserCount:     1,
					SalesTotal:    100,
					DiscountTotal: 0,
				},
			},
			expect: TopOrderSalesTrends{
				{
					TopOrderSalesTrend: response.TopOrderSalesTrend{
						Period:     "2025-01-01",
						SalesTotal: 100,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewTopOrderSalesTrends(tt.periodType, tt.startAt, tt.endAt, tt.aggregated)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestTopOrderSalesTrends_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		trends TopOrderSalesTrends
		expect []*response.TopOrderSalesTrend
	}{
		{
			name: "success",
			trends: TopOrderSalesTrends{
				{
					TopOrderSalesTrend: response.TopOrderSalesTrend{
						Period:     "2025-01-18",
						SalesTotal: 100,
					},
				},
			},
			expect: []*response.TopOrderSalesTrend{
				{
					Period:     "2025-01-18",
					SalesTotal: 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.trends.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
