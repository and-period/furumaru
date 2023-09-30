package response

// Order - 注文履歴情報
type Order struct {
	ID          string            `json:"id"`          // 注文履歴ID
	ScheduleID  string            `json:"scheduleId"`  // 開催スケジュールID
	PromotionID string            `json:"promotionId"` // プロモーションID
	UserID      string            `json:"userId"`      // ユーザーID
	UserName    string            `json:"userName"`    // 注文者名
	Payment     *OrderPayment     `json:"payment"`     // 支払い情報
	Fulfillment *OrderFulfillment `json:"fulfillment"` // 配送情報
	Refund      *OrderRefund      `json:"refund"`      // 注文キャンセル情報
	Items       []*OrderItem      `json:"items"`       // 注文商品一覧
	OrderedAt   int64             `json:"orderedAt"`   // 注文日時
	PaidAt      int64             `json:"paidAt"`      // 支払日時
	DeliveredAt int64             `json:"deliveredAt"` // 配送日時
	CanceledAt  int64             `json:"canceledAt"`  // 注文キャンセル日時
	CreatedAt   int64             `json:"createdAt"`   // 登録日時
	UpdatedAt   int64             `json:"updatedAt"`   // 更新日時
}

// OrderPayment - 支払い情報
type OrderPayment struct {
	TransactionID string `json:"transactionId"` // 取引ID
	MethodType    int32  `json:"methodType"`    // 決済手段種別
	Status        int32  `json:"status"`        // 支払い状況
	Subtotal      int64  `json:"subtotal"`      // 購入金額
	Discount      int64  `json:"discount"`      // 割引金額
	ShippingFee   int64  `json:"shippingFee"`   // 配送手数料
	Tax           int64  `json:"tax"`           // 消費税
	Total         int64  `json:"total"`         // 合計金額
	AddressID     string `json:"addressId"`     // 請求先情報ID
	*Address             // 請求先情報
}

// OrderFulfillment - 配送情報
type OrderFulfillment struct {
	TrackingNumber  string `json:"trackingNumber"`  // 伝票番号
	Status          int32  `json:"status"`          // 配送状況
	ShippingCarrier int32  `json:"shippingCarrier"` // 配送会社
	ShippingMethod  int32  `json:"shippingMethod"`  // 配送方法
	BoxSize         int32  `json:"boxSize"`         // 箱の大きさ
	AddressID       string `json:"addressId"`       // 配送先情報ID
	*Address               // 配送先情報
}

// OrderRefund - 注文キャンセル情報
type OrderRefund struct {
	Canceled bool   `json:"canceled"` // 注文キャンセルフラグ
	Reason   string `json:"reason"`   // 注文キャンセル理由詳細
	Total    int64  `json:"total"`    // 返金金額
}

// OrderItem - 注文商品情報
type OrderItem struct {
	ProductID string          `json:"productId"` // 商品ID
	Name      string          `json:"name"`      // 商品名
	Price     int64           `json:"price"`     // 購入価格
	Quantity  int64           `json:"quantity"`  // 購入数量
	Weight    float64         `json:"weight"`    // 重量(kg,少数第一位まで)
	Media     []*ProductMedia `json:"media"`     // メディア一覧
}

type OrderResponse struct {
	*Order
}

type OrdersResponse struct {
	Orders []*Order `json:"orders"` // 注文履歴一覧
	Total  int64    `json:"total"`  // 注文履歴合計数
}
