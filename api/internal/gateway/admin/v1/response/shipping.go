package response

// Shipping - 配送設定情報
type Shipping struct {
	ID                 string          `json:"id"`                 // 配送設定ID
	Name               string          `json:"name"`               // 配送設定名
	Box60Rates         []*ShippingRate `json:"box60Rates"`         // 箱サイズ60の通常便配送料一覧
	Box60Refrigerated  int64           `json:"box60Refrigerated"`  // 箱サイズ60の冷蔵便追加配送料
	Box60Frozen        int64           `json:"box60Frozen"`        // 箱サイズ60の冷凍便追加配送料
	Box80Rates         []*ShippingRate `json:"box80Rates"`         // 箱サイズ80の通常便配送料一覧
	Box80Refrigerated  int64           `json:"box80Refrigerated"`  // 箱サイズ80の冷蔵便追加配送料
	Box80Frozen        int64           `json:"box80Frozen"`        // 箱サイズ80の冷凍便追加配送料
	Box100Rates        []*ShippingRate `json:"box100Rates"`        // 箱サイズ100の通常便配送料一覧
	Box100Refrigerated int64           `json:"box100Refrigerated"` // 箱サイズ100の冷蔵便追加配送料
	Box100Frozen       int64           `json:"box100Frozen"`       // 箱サイズ100の冷凍便追加配送料
	HasFreeShipping    bool            `json:"hasFreeShipping"`    // 送料無料オプションの有無
	FreeShippingRates  int64           `json:"freeShippingRates"`  // 送料無料になる金額
	CreatedAt          int64           `json:"createdAt"`          // 登録日時
	UpdatedAt          int64           `json:"updatedAt"`          // 更新日時
}

type ShippingRate struct {
	Number      int64    `json:"number"`      // No.
	Name        string   `json:"name"`        // 配送料金設定名
	Price       int64    `json:"price"`       // 配送料金
	Prefectures []string `json:"prefectures"` // 対象都道府県一覧
}

type ShippingResponse struct {
	*Shipping
}

type ShippingsResponse struct {
	Shippings []*Shipping `json:"shippings"` // 配送設定一覧
	Total     int64       `json:"total"`     // 合計数
}
