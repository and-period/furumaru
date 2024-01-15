package response

// Cart - カート情報
type Cart struct {
	Number        int64       `json:"number"`        // 箱の通番
	Type          int32       `json:"type"`          // 箱の種別
	Size          int32       `json:"size"`          // 箱のサイズ
	Rate          int64       `json:"rate"`          // 箱の占有率
	Items         []*CartItem `json:"items"`         // 箱の商品一覧
	CoordinatorID string      `json:"coordinatorId"` // コーディネータID
}

// CartItem - カート内の商品情報
type CartItem struct {
	ProductID string `json:"productId"` // 商品ID
	Quantity  int64  `json:"quantity"`  // 数量
}

type CartResponse struct {
	Carts        []*Cart        `json:"carts"`        // カート一覧
	Coordinators []*Coordinator `json:"coordinators"` // コーディネータ一覧
	Products     []*Product     `json:"products"`     // 商品一覧
}

type CalcCartResponse struct {
	Carts       []*Cart      `json:"carts"`       // カート一覧
	Items       []*CartItem  `json:"items"`       // カート内の商品一覧(集計結果)
	Products    []*Product   `json:"products"`    // 商品一覧
	Coordinator *Coordinator `json:"coordinator"` // コーディネータ情報
	Promotion   *Promotion   `json:"promotion"`   // プロモーション情報
	SubTotal    int64        `json:"subtotal"`    // 購入金額(税込)
	Discount    int64        `json:"discount"`    // 割引金額(税込)
	ShippingFee int64        `json:"shippingFee"` // 配送手数料(税込)
	Total       int64        `json:"total"`       // 合計金額(税込)
	RequestID   string       `json:"requestId"`   // 支払い時にAPIへ送信するキー(重複判定用)
}
