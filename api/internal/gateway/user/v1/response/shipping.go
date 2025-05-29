package response

// Shipping - 配送設定情報
type Shipping struct {
	ID                string          `json:"id"`                // 配送設定ID
	Box60Rates        []*ShippingRate `json:"box60Rates"`        // 箱サイズ60の通常（常温・冷蔵便）配送料一覧
	Box60Frozen       int64           `json:"box60Frozen"`       // 箱サイズ60の追加（冷凍便）追加配送料(税込)
	Box80Rates        []*ShippingRate `json:"box80Rates"`        // 箱サイズ80の通常（常温・冷蔵便）配送料一覧
	Box80Frozen       int64           `json:"box80Frozen"`       // 箱サイズ80の追加（冷凍便）追加配送料(税込)
	Box100Rates       []*ShippingRate `json:"box100Rates"`       // 箱サイズ100の通常（常温・冷蔵便）配送料一覧
	Box100Frozen      int64           `json:"box100Frozen"`      // 箱サイズ100の追加（冷凍便）追加配送料(税込)
	HasFreeShipping   bool            `json:"hasFreeShipping"`   // 送料無料オプションの有無
	FreeShippingRates int64           `json:"freeShippingRates"` // 送料無料になる金額(税込)
}

type ShippingRate struct {
	Number          int64    `json:"number"`          // No.
	Name            string   `json:"name"`            // 配送料金設定名
	Price           int64    `json:"price"`           // 配送料金(税込)
	Prefectures     []string `json:"prefectures"`     // 対象都道府県名
	PrefectureCodes []int32  `json:"prefectureCodes"` // 対象都道府県一覧
}
