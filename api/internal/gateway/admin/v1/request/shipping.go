package request

type CreateShippingRequest struct {
	Box60Rates        []*CreateShippingRate `json:"box60Rates,omitempty"`        // 箱サイズ60の通常便配送料一覧
	Box60Frozen       int64                 `json:"box60Frozen,omitempty"`       // 箱サイズ60の冷凍便追加配送料(税込)
	Box80Rates        []*CreateShippingRate `json:"box80Rates,omitempty"`        // 箱サイズ80の通常便配送料一覧
	Box80Frozen       int64                 `json:"box80Frozen,omitempty"`       // 箱サイズ80の冷凍便追加配送料(税込)
	Box100Rates       []*CreateShippingRate `json:"box100Rates,omitempty"`       // 箱サイズ100の通常便配送料一覧
	Box100Frozen      int64                 `json:"box100Frozen,omitempty"`      // 箱サイズ100の冷凍便追加配送料(税込)
	HasFreeShipping   bool                  `json:"hasFreeShipping,omitempty"`   // 送料無料オプションの有無
	FreeShippingRates int64                 `json:"freeShippingRates,omitempty"` // 送料無料になる金額(税込)
}

type CreateShippingRate struct {
	Name            string  `json:"name,omitempty"`            // 配送料金設定名
	Price           int64   `json:"price,omitempty"`           // 配送料金(税込)
	PrefectureCodes []int32 `json:"prefectureCodes,omitempty"` // 対象都道府県一覧
}

type UpdateShippingRequest struct {
	Box60Rates        []*UpdateShippingRate `json:"box60Rates,omitempty"`        // 箱サイズ60の通常便配送料一覧
	Box60Frozen       int64                 `json:"box60Frozen,omitempty"`       // 箱サイズ60の冷凍便追加配送料(税込)
	Box80Rates        []*UpdateShippingRate `json:"box80Rates,omitempty"`        // 箱サイズ80の通常便配送料一覧
	Box80Frozen       int64                 `json:"box80Frozen,omitempty"`       // 箱サイズ80の冷凍便追加配送料(税込)
	Box100Rates       []*UpdateShippingRate `json:"box100Rates,omitempty"`       // 箱サイズ100の通常便配送料一覧
	Box100Frozen      int64                 `json:"box100Frozen,omitempty"`      // 箱サイズ100の冷凍便追加配送料(税込)
	HasFreeShipping   bool                  `json:"hasFreeShipping,omitempty"`   // 送料無料オプションの有無
	FreeShippingRates int64                 `json:"freeShippingRates,omitempty"` // 送料無料になる金額(税込)
}

type UpdateShippingRate struct {
	Name            string  `json:"name,omitempty"`            // 配送料金設定名
	Price           int64   `json:"price,omitempty"`           // 配送料金(税込)
	PrefectureCodes []int32 `json:"prefectureCodes,omitempty"` // 対象都道府県一覧
}

type UpsertShippingRequest struct {
	Box60Rates        []*UpsertShippingRate `json:"box60Rates,omitempty"`        // 箱サイズ60の通常便配送料一覧
	Box60Frozen       int64                 `json:"box60Frozen,omitempty"`       // 箱サイズ60の冷凍便追加配送料(税込)
	Box80Rates        []*UpsertShippingRate `json:"box80Rates,omitempty"`        // 箱サイズ80の通常便配送料一覧
	Box80Frozen       int64                 `json:"box80Frozen,omitempty"`       // 箱サイズ80の冷凍便追加配送料(税込)
	Box100Rates       []*UpsertShippingRate `json:"box100Rates,omitempty"`       // 箱サイズ100の通常便配送料一覧
	Box100Frozen      int64                 `json:"box100Frozen,omitempty"`      // 箱サイズ100の冷凍便追加配送料(税込)
	HasFreeShipping   bool                  `json:"hasFreeShipping,omitempty"`   // 送料無料オプションの有無
	FreeShippingRates int64                 `json:"freeShippingRates,omitempty"` // 送料無料になる金額(税込)
}

type UpsertShippingRate struct {
	Name            string  `json:"name,omitempty"`            // 配送料金設定名
	Price           int64   `json:"price,omitempty"`           // 配送料金(税込)
	PrefectureCodes []int32 `json:"prefectureCodes,omitempty"` // 対象都道府県一覧
}

type UpdateDefaultShippingRequest struct {
	Box60Rates        []*UpdateDefaultShippingRate `json:"box60Rates,omitempty"`        // 箱サイズ60の通常便配送料一覧
	Box60Frozen       int64                        `json:"box60Frozen,omitempty"`       // 箱サイズ60の冷凍便追加配送料(税込)
	Box80Rates        []*UpdateDefaultShippingRate `json:"box80Rates,omitempty"`        // 箱サイズ80の通常便配送料一覧
	Box80Frozen       int64                        `json:"box80Frozen,omitempty"`       // 箱サイズ80の冷凍便追加配送料(税込)
	Box100Rates       []*UpdateDefaultShippingRate `json:"box100Rates,omitempty"`       // 箱サイズ100の通常便配送料一覧
	Box100Frozen      int64                        `json:"box100Frozen,omitempty"`      // 箱サイズ100の冷凍便追加配送料(税込)
	HasFreeShipping   bool                         `json:"hasFreeShipping,omitempty"`   // 送料無料オプションの有無
	FreeShippingRates int64                        `json:"freeShippingRates,omitempty"` // 送料無料になる金額(税込)
}

type UpdateDefaultShippingRate struct {
	Name            string  `json:"name,omitempty"`            // 配送料金設定名
	Price           int64   `json:"price,omitempty"`           // 配送料金(税込)
	PrefectureCodes []int32 `json:"prefectureCodes,omitempty"` // 対象都道府県一覧
}
