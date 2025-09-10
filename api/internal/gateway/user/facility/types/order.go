package types

// OrderType - 注文種別
type OrderType int32

const (
	OrderTypeUnknown    OrderType = 0
	OrderTypeProduct    OrderType = 1 // 商品
	OrderTypeExperience OrderType = 2 // 体験
)

// OrderStatus - 注文ステータス
type OrderStatus int32

const (
	OrderStatusUnknown   OrderStatus = 0
	OrderStatusUnpaid    OrderStatus = 1 // 支払い待ち
	OrderStatusPreparing OrderStatus = 2 // 発送対応中
	OrderStatusCompleted OrderStatus = 3 // 完了
	OrderStatusCanceled  OrderStatus = 4 // キャンセル
	OrderStatusRefunded  OrderStatus = 5 // 返金
	OrderStatusFailed    OrderStatus = 6 // 失敗
)

// Order - 注文履歴情報
type Order struct {
	ID             string        `json:"id"`             // 注文履歴ID
	CoordinatorID  string        `json:"coordinatorId"`  // コーディネータID
	PromotionID    string        `json:"promotionId"`    // プロモーションID
	Type           OrderType     `json:"type"`           // 注文種別
	Status         OrderStatus   `json:"status"`         // 注文ステータス
	PickupAt       int64         `json:"pickupAt"`       // 受け取り日時
	PickupLocation string        `json:"pickupLocation"` // 受け取り場所
	Payment        *OrderPayment `json:"payment"`        // 支払い情報
	Refund         *OrderRefund  `json:"refund"`         // 注文キャンセル情報
	Items          []*OrderItem  `json:"items"`          // 注文商品一覧
}

// OrderPayment - 支払い情報
type OrderPayment struct {
	TransactionID string            `json:"transactionId"` // 取引ID
	MethodType    PaymentMethodType `json:"methodType"`    // 決済手段種別
	Status        PaymentStatus     `json:"status"`        // 支払い状況
	Subtotal      int64             `json:"subtotal"`      // 購入金額(税込)
	Discount      int64             `json:"discount"`      // 割引金額(税込)
	ShippingFee   int64             `json:"shippingFee"`   // 配送手数料(税込)
	Total         int64             `json:"total"`         // 合計金額(税込)
	OrderedAt     int64             `json:"orderedAt"`     // 注文日時
	PaidAt        int64             `json:"paidAt"`        // 支払日時
}

// OrderRefund - 注文キャンセル情報
type OrderRefund struct {
	Total      int64      `json:"total"`      // 返金金額
	Type       RefundType `json:"type"`       // 注文キャンセル種別
	Reason     string     `json:"reason"`     // 注文キャンセル理由
	Canceled   bool       `json:"canceled"`   // 注文キャンセルフラグ
	CanceledAt int64      `json:"canceledAt"` // 注文キャンセル日時
}

// OrderItem - 注文商品情報
type OrderItem struct {
	FulfillmentID string `json:"fulfillmentId"` // 配送情報ID
	ProductID     string `json:"productId"`     // 商品ID
	Price         int64  `json:"price"`         // 購入価格(税込)
	Quantity      int64  `json:"quantity"`      // 購入数量
}

type OrderResponse struct {
	Order       *Order       `json:"order"`       // 注文履歴情報
	Coordinator *Coordinator `json:"coordinator"` // コーディネータ情報
	Promotion   *Promotion   `json:"promotion"`   // プロモーション情報
	Products    []*Product   `json:"products"`    // 商品一覧
}

type OrdersResponse struct {
	Order        []*Order       `json:"orders"`       // 注文履歴一覧
	Coordinators []*Coordinator `json:"coordinators"` // コーディネータ一覧
	Promotions   []*Promotion   `json:"promotions"`   // プロモーション一覧
	Products     []*Product     `json:"products"`     // 商品一覧
	Total        int64          `json:"total"`        // 合計数
}
