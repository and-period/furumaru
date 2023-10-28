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
