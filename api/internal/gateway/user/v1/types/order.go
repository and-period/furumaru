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
	ID              string              `json:"id"`              // 注文履歴ID
	CoordinatorID   string              `json:"coordinatorId"`   // コーディネータID
	PromotionID     string              `json:"promotionId"`     // プロモーションID
	Type            OrderType           `json:"type"`            // 注文種別
	Status          OrderStatus         `json:"status"`          // 注文ステータス
	Payment         *OrderPayment       `json:"payment"`         // 支払い情報
	Refund          *OrderRefund        `json:"refund"`          // 注文キャンセル情報
	Fulfillments    []*OrderFulfillment `json:"fulfillments"`    // 配送情報一覧
	Items           []*OrderItem        `json:"items"`           // 注文商品一覧
	Experience      *OrderExperience    `json:"experience"`      // 注文体験情報
	BillingAddress  *Address            `json:"shippingAddress"` // 請求先情報
	ShippingAddress *Address            `json:"billingAddress"`  // 配送先情報
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

// OrderFulfillment - 配送情報
type OrderFulfillment struct {
	FulfillmentID   string            `json:"fulfillmentId"`   // 配送情報ID
	TrackingNumber  string            `json:"trackingNumber"`  // 伝票番号
	Status          FulfillmentStatus `json:"status"`          // 配送状況
	ShippingCarrier ShippingCarrier   `json:"shippingCarrier"` // 配送会社
	ShippingType    ShippingType      `json:"shippingType"`    // 配送方法
	BoxNumber       int64             `json:"boxNumber"`       // 箱の通番
	BoxSize         ShippingSize      `json:"boxSize"`         // 箱の大きさ
	BoxRate         int64             `json:"boxRate"`         // 箱の占有率
	ShippedAt       int64             `json:"shippedAt"`       // 配送日時
}

// OrderItem - 注文商品情報
type OrderItem struct {
	FulfillmentID string `json:"fulfillmentId"` // 配送情報ID
	ProductID     string `json:"productId"`     // 商品ID
	Price         int64  `json:"price"`         // 購入価格(税込)
	Quantity      int64  `json:"quantity"`      // 購入数量
}

// OrderExperience - 注文体験情報
type OrderExperience struct {
	ExperienceID          string `json:"experienceId"`          // 体験ID
	AdultCount            int64  `json:"adultCount"`            // 大人人数
	AdultPrice            int64  `json:"adultPrice"`            // 大人価格
	JuniorHighSchoolCount int64  `json:"juniorHighSchoolCount"` // 中学生人数
	JuniorHighSchoolPrice int64  `json:"juniorHighSchoolPrice"` // 中学生価格
	ElementarySchoolCount int64  `json:"elementarySchoolCount"` // 小学生人数
	ElementarySchoolPrice int64  `json:"elementarySchoolPrice"` // 小学生価格
	PreschoolCount        int64  `json:"preschoolCount"`        // 幼児人数
	PreschoolPrice        int64  `json:"preschoolPrice"`        // 幼児価格
	SeniorCount           int64  `json:"seniorCount"`           // シニア人数
	SeniorPrice           int64  `json:"seniorPrice"`           // シニア価格
	Transportation        string `json:"transportation"`        // 交通手段
	RequestedDate         string `json:"requestedDate"`         // 体験希望日
	RequestedTime         string `json:"requestedTime"`         // 体験希望時間
}

type OrderResponse struct {
	Order       *Order       `json:"order"`       // 注文履歴情報
	Coordinator *Coordinator `json:"coordinator"` // コーディネータ情報
	Promotion   *Promotion   `json:"promotion"`   // プロモーション情報
	Products    []*Product   `json:"products"`    // 商品一覧
	Experience  *Experience  `json:"experience"`  // 体験情報
}

type OrdersResponse struct {
	Order        []*Order       `json:"orders"`       // 注文履歴一覧
	Coordinators []*Coordinator `json:"coordinators"` // コーディネータ一覧
	Promotions   []*Promotion   `json:"promotions"`   // プロモーション一覧
	Products     []*Product     `json:"products"`     // 商品一覧
	Experiences  []*Experience  `json:"experiences"`  // 体験一覧
	Total        int64          `json:"total"`        // 合計数
}
