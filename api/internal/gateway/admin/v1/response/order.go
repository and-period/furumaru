package response

// Order - 注文履歴情報
type Order struct {
	ID              string              `json:"id"`              // 注文履歴ID
	UserID          string              `json:"userId"`          // ユーザーID
	CoordinatorID   string              `json:"coordinatorId"`   // コーディネータID
	PromotionID     string              `json:"promotionId"`     // プロモーションID
	ManagementID    int64               `json:"managementId"`    // 注文管理用ID
	ShippingMessage string              `json:"shippingMessage"` // 発送連絡時のメッセージ
	Type            int32               `json:"type"`            // 注文種別
	Status          int32               `json:"status"`          // 注文ステータス
	Payment         *OrderPayment       `json:"payment"`         // 支払い情報
	Refund          *OrderRefund        `json:"refund"`          // 注文キャンセル情報
	Fulfillments    []*OrderFulfillment `json:"fulfillments"`    // 配送情報一覧
	Items           []*OrderItem        `json:"items"`           // 注文商品一覧
	CreatedAt       int64               `json:"createdAt"`       // 登録日時
	UpdatedAt       int64               `json:"updatedAt"`       // 更新日時
	CompletedAt     int64               `json:"completedAt"`     // 対応完了日時
}

// OrderPayment - 支払い情報
type OrderPayment struct {
	TransactionID string `json:"transactionId"` // 取引ID
	MethodType    int32  `json:"methodType"`    // 決済手段種別
	Status        int32  `json:"status"`        // 支払い状況
	Subtotal      int64  `json:"subtotal"`      // 購入金額(税込)
	Discount      int64  `json:"discount"`      // 割引金額(税込)
	ShippingFee   int64  `json:"shippingFee"`   // 配送手数料(税込)
	Total         int64  `json:"total"`         // 合計金額(税込)
	OrderedAt     int64  `json:"orderedAt"`     // 注文日時
	PaidAt        int64  `json:"paidAt"`        // 支払日時
	*Address             // 請求先情報
}

// OrderRefund - 注文キャンセル情報
type OrderRefund struct {
	Total      int64  `json:"total"`      // 返金金額
	Type       int32  `json:"type"`       // 注文キャンセル種別
	Reason     string `json:"reason"`     // 注文キャンセル理由
	Canceled   bool   `json:"canceled"`   // 注文キャンセルフラグ
	CanceledAt int64  `json:"canceledAt"` // 注文キャンセル日時
}

// OrderFulfillment - 配送情報
type OrderFulfillment struct {
	FulfillmentID   string `json:"fulfillmentId"`   // 配送情報ID
	TrackingNumber  string `json:"trackingNumber"`  // 伝票番号
	Status          int32  `json:"status"`          // 配送状況
	ShippingCarrier int32  `json:"shippingCarrier"` // 配送会社
	ShippingType    int32  `json:"shippingType"`    // 配送方法
	BoxNumber       int64  `json:"boxNumber"`       // 箱の通番
	BoxSize         int32  `json:"boxSize"`         // 箱の大きさ
	BoxRate         int64  `json:"boxRate"`         // 箱の占有率
	ShippedAt       int64  `json:"shippedAt"`       // 配送日時
	*Address               // 配送先情報
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
	User        *User        `json:"user"`        // 購入者情報
	Coordinator *Coordinator `json:"coordinator"` // コーディネータ情報
	Promotion   *Promotion   `json:"promotion"`   // プロモーション情報
	Products    []*Product   `json:"products"`    // 商品一覧
}

type OrdersResponse struct {
	Orders       []*Order       `json:"orders"`       // 注文履歴一覧
	Users        []*User        `json:"users"`        // 購入者一覧
	Coordinators []*Coordinator `json:"coordinators"` // コーディネータ一覧
	Promotions   []*Promotion   `json:"promotions"`   // プロモーション一覧
	Total        int64          `json:"total"`        // 注文履歴合計数
}
