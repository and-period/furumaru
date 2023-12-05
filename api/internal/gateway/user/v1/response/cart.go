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
	Carts       []*Cart      `json:"carts"`        // カート一覧
	Items       []*CartItem  `json:"items"`        // カート内の商品一覧(集計結果)
	Products    []*Product   `json:"products"`     // 商品一覧
	Coordinator *Coordinator `json:"coordinators"` // コーディネータ情報
	Promotion   *Promotion   `json:"promotion"`    // プロモーション情報
	SubTotal    int64        `json:"subtotal"`     // 購入金額
	Discount    int64        `json:"discount"`     // 割引金額
	ShippingFee int64        `json:"shippingFee"`  // 配送手数料
	Tax         int64        `json:"tax"`          // 消費税
	TaxRate     int64        `json:"taxRate"`      // 消費税率(%)
	Total       int64        `json:"total"`        // 合計金額
}
