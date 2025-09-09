package types

type TopOrderValue struct {
	Value      int64   `json:"value"`      // 値
	Comparison float64 `json:"comparison"` // 比較値（％：前日比など）
}

type TopOrderSalesTrend struct {
	Period     string `json:"period"`     // 期間
	SalesTotal int64  `json:"salesTotal"` // 売上合計
}

type TopOrderPayment struct {
	PaymentMethodType int32   `json:"paymentMethodType"` // 支払い方法
	OrderCount        int64   `json:"orderCount"`        // 注文数
	UserCount         int64   `json:"userCount"`         // ユーザー数
	SalesTotal        int64   `json:"salesTotal"`        // 売上合計
	Rate              float64 `json:"rate"`              // 割合（支払い方法別注文数 / 注文数）
}

type TopOrdersResponse struct {
	StartAt     int64                 `json:"startAt"`     // 開始日時
	EndAt       int64                 `json:"endAt"`       // 終了日時
	PeriodType  string                `json:"periodType"`  // 期間種別
	Orders      *TopOrderValue        `json:"orders"`      // 注文数
	Users       *TopOrderValue        `json:"users"`       // ユーザー数
	Sales       *TopOrderValue        `json:"sales"`       // 売上合計
	Payments    []*TopOrderPayment    `json:"payments"`    // 支払い方法別情報
	SalesTrends []*TopOrderSalesTrend `json:"salesTrends"` // 売上推移（グラフ描画用）
}
