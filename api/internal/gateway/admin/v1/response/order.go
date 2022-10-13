package response

// Order - 注文履歴情報
type Order struct {
	ID          string            `json:"id"`          // 注文履歴ID
	UserID      string            `json:"userId"`      // ユーザーID
	ScheduleID  string            `json:"scheduleId"`  // 開催スケジュールID
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
	TransactionID  string `json:"transactionId"`  // 取引ID
	PromotionID    string `json:"promotionId"`    // プロモーションID
	PaymentID      string `json:"paymentId"`      // 決済手段ID
	PaymentType    int32  `json:"paymentType"`    // 決済手段
	Status         int32  `json:"status"`         // 支払い状況
	Subtotal       int64  `json:"subtotal"`       // 購入金額
	Discount       int64  `json:"discount"`       // 割引金額
	ShippingCharge int64  `json:"shippingCharge"` // 配送料金
	Tax            int64  `json:"tax"`            // 消費税
	Total          int64  `json:"total"`          // 支払合計
	Lastname       string `json:"lastname"`       // 請求先情報 姓
	Firstname      string `json:"firstname"`      // 請求先情報 名
	PostalCode     string `json:"postalCode"`     // 請求先情報 郵便番号
	Prefecture     string `json:"prefecture"`     // 請求先情報 都道府県
	City           string `json:"city"`           // 請求先情報 市区町村
	AddressLine1   string `json:"addressLine1"`   // 請求先情報 町名・番地
	AddressLine2   string `json:"addressLine2"`   // 請求先情報 ビル名・号室など
	PhoneNumber    string `json:"phoneNumber"`    // 請求先情報 電話番号
}

// OrderFulfillment - 配送情報
type OrderFulfillment struct {
	TrackingNumber  string  `json:"trackingNumber"`  // 伝票番号
	Status          int32   `json:"status"`          // 配送状況
	ShippingCarrier int32   `json:"shippingCarrier"` // 配送会社
	ShippingMethod  int32   `json:"shippingMethod"`  // 配送方法
	BoxSize         int32   `json:"boxSize"`         // 箱の大きさ
	BoxCount        int64   `json:"boxCount"`        // 箱の個数
	WeightTotal     float64 `json:"weight"`          // 合計重量(kg,少数第一位まで)
	Lastname        string  `json:"lastname"`        // 配送先情報 姓
	Firstname       string  `json:"firstname"`       // 配送先情報 名
	PostalCode      string  `json:"postalCode"`      // 配送先情報 郵便番号
	Prefecture      string  `json:"prefecture"`      // 配送先情報 都道府県
	City            string  `json:"city"`            // 配送先情報 市区町村
	AddressLine1    string  `json:"addressLine1"`    // 配送先情報 町名・番地
	AddressLine2    string  `json:"addressLine2"`    // 配送先情報 ビル名・号室など
	PhoneNumber     string  `json:"phoneNumber"`     // 配送先情報 電話番号
}

// OrderRefund - 注文キャンセル情報
type OrderRefund struct {
	Canceled bool   `json:"canceled"` // 注文キャンセルフラグ
	Type     int32  `json:"type"`     // 注文キャンセル種別
	Reason   string `json:"reason"`   // 注文キャンセル理由詳細
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
