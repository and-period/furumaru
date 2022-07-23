package request

type CreateShippingRequest struct {
	Name               string                `json:"name,omitempty"`               // 配送設定名
	Box60Rates         []*CreateShippingRate `json:"box60Rates,omitempty"`         // 箱サイズ60の通常便配送料一覧
	Box60Refrigerated  int64                 `json:"box60Refrigerated,omitempty"`  // 箱サイズ60の冷蔵便追加配送料
	Box60Frozen        int64                 `json:"box60Frozen,omitempty"`        // 箱サイズ60の冷凍便追加配送料
	Box80Rates         []*CreateShippingRate `json:"box80Rates,omitempty"`         // 箱サイズ80の通常便配送料一覧
	Box80Refrigerated  int64                 `json:"box80Refrigerated,omitempty"`  // 箱サイズ80の冷蔵便追加配送料
	Box80Frozen        int64                 `json:"box80Frozen,omitempty"`        // 箱サイズ80の冷凍便追加配送料
	Box100Rates        []*CreateShippingRate `json:"box100Rates,omitempty"`        // 箱サイズ100の通常便配送料一覧
	Box100Refrigerated int64                 `json:"box100Refrigerated,omitempty"` // 箱サイズ100の冷蔵便追加配送料
	Box100Frozen       int64                 `json:"box100Frozen,omitempty"`       // 箱サイズ100の冷凍便追加配送料
	HasFreeShipping    bool                  `json:"hasFreeShipping,omitempty"`    // 送料無料オプションの有無
	FreeShippingRates  int64                 `json:"freeShippingRates,omitempty"`  // 送料無料になる金額
}

type CreateShippingRate struct {
	Name        string   `json:"name,omitempty"`        // 配送料金設定名
	Price       int64    `json:"price,omitempty"`       // 配送料金
	Prefectures []string `json:"prefectures,omitempty"` // 対象都道府県一覧
}

type UpdateShippingRequest struct {
	Name               string                `json:"name,omitempty"`               // 配送設定名
	Box60Rates         []*UpdateShippingRate `json:"box60Rates,omitempty"`         // 箱サイズ60の通常便配送料一覧
	Box60Refrigerated  int64                 `json:"box60Refrigerated,omitempty"`  // 箱サイズ60の冷蔵便追加配送料
	Box60Frozen        int64                 `json:"box60Frozen,omitempty"`        // 箱サイズ60の冷凍便追加配送料
	Box80Rates         []*UpdateShippingRate `json:"box80Rates,omitempty"`         // 箱サイズ80の通常便配送料一覧
	Box80Refrigerated  int64                 `json:"box80Refrigerated,omitempty"`  // 箱サイズ80の冷蔵便追加配送料
	Box80Frozen        int64                 `json:"box80Frozen,omitempty"`        // 箱サイズ80の冷凍便追加配送料
	Box100Rates        []*UpdateShippingRate `json:"box100Rates,omitempty"`        // 箱サイズ100の通常便配送料一覧
	Box100Refrigerated int64                 `json:"box100Refrigerated,omitempty"` // 箱サイズ100の冷蔵便追加配送料
	Box100Frozen       int64                 `json:"box100Frozen,omitempty"`       // 箱サイズ100の冷凍便追加配送料
	HasFreeShipping    bool                  `json:"hasFreeShipping,omitempty"`    // 送料無料オプションの有無
	FreeShippingRates  int64                 `json:"freeShippingRates,omitempty"`  // 送料無料になる金額
}

type UpdateShippingRate struct {
	Name        string   `json:"name,omitempty"`        // 配送料金設定名
	Price       int64    `json:"price,omitempty"`       // 配送料金
	Prefectures []string `json:"prefectures,omitempty"` // 対象都道府県一覧
}
