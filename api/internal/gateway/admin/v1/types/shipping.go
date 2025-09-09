package types

// Shipping - 配送設定情報
type Shipping struct {
	ID                string          `json:"id"`                // 配送設定ID
	Name              string          `json:"name"`              // 配送設定名
	IsDefault         bool            `json:"isDefault"`         // デフォルト設定フラグ
	Box60Rates        []*ShippingRate `json:"box60Rates"`        // 箱サイズ60の通常（常温・冷蔵便）配送料一覧
	Box60Frozen       int64           `json:"box60Frozen"`       // 箱サイズ60の追加（冷凍便）追加配送料(税込)
	Box80Rates        []*ShippingRate `json:"box80Rates"`        // 箱サイズ80の通常（常温・冷蔵便）配送料一覧
	Box80Frozen       int64           `json:"box80Frozen"`       // 箱サイズ80の追加（冷凍便）追加配送料(税込)
	Box100Rates       []*ShippingRate `json:"box100Rates"`       // 箱サイズ100の通常（常温・冷蔵便）配送料一覧
	Box100Frozen      int64           `json:"box100Frozen"`      // 箱サイズ100の追加（冷凍便）追加配送料(税込)
	HasFreeShipping   bool            `json:"hasFreeShipping"`   // 送料無料オプションの有無
	FreeShippingRates int64           `json:"freeShippingRates"` // 送料無料になる金額(税込)
	CreatedAt         int64           `json:"createdAt"`         // 登録日時
	UpdatedAt         int64           `json:"updatedAt"`         // 更新日時
}

type ShippingRate struct {
	Number          int64   `json:"number"`          // No.
	Name            string  `json:"name"`            // 配送料金設定名
	Price           int64   `json:"price"`           // 配送料金(税込)
	PrefectureCodes []int32 `json:"prefectureCodes"` // 対象都道府県一覧
}

type CreateShippingRequest struct {
	Name              string                `json:"name" validate:"required,max=64"`                       // 配送設定名
	Box60Rates        []*CreateShippingRate `json:"box60Rates" validate:"required,dive"`                   // 箱サイズ60の通常便配送料一覧
	Box60Frozen       int64                 `json:"box60Frozen" validate:"required,min=0,lt=10000000000"`  // 箱サイズ60の冷凍便追加配送料(税込)
	Box80Rates        []*CreateShippingRate `json:"box80Rates" validate:"required,dive"`                   // 箱サイズ80の通常便配送料一覧
	Box80Frozen       int64                 `json:"box80Frozen" validate:"required,min=0,lt=10000000000"`  // 箱サイズ80の冷凍便追加配送料(税込)
	Box100Rates       []*CreateShippingRate `json:"box100Rates" validate:"required,dive"`                  // 箱サイズ100の通常便配送料一覧
	Box100Frozen      int64                 `json:"box100Frozen" validate:"required,min=0,lt=10000000000"` // 箱サイズ100の冷凍便追加配送料(税込)
	HasFreeShipping   bool                  `json:"hasFreeShipping"`                                       // 送料無料オプションの有無
	FreeShippingRates int64                 `json:"freeShippingRates" validate:"min=0,lt=10000000000"`     // 送料無料になる金額(税込)
}

type CreateShippingRate struct {
	Name            string  `json:"name" validate:"required,max=64"`                       // 配送料金設定名
	Price           int64   `json:"price" validate:"required,min=0,lt=10000000000"`        // 配送料金(税込)
	PrefectureCodes []int32 `json:"prefectureCodes" validate:"required,dive,min=1,max=47"` // 対象都道府県一覧
}

type UpdateShippingRequest struct {
	Name              string                `json:"name" validate:"required,max=64"`                       // 配送設定名
	Box60Rates        []*UpdateShippingRate `json:"box60Rates" validate:"required,dive"`                   // 箱サイズ60の通常便配送料一覧
	Box60Frozen       int64                 `json:"box60Frozen" validate:"required,min=0,lt=10000000000"`  // 箱サイズ60の冷凍便追加配送料(税込)
	Box80Rates        []*UpdateShippingRate `json:"box80Rates" validate:"required,dive"`                   // 箱サイズ80の通常便配送料一覧
	Box80Frozen       int64                 `json:"box80Frozen" validate:"required,min=0,lt=10000000000"`  // 箱サイズ80の冷凍便追加配送料(税込)
	Box100Rates       []*UpdateShippingRate `json:"box100Rates" validate:"required,dive"`                  // 箱サイズ100の通常便配送料一覧
	Box100Frozen      int64                 `json:"box100Frozen" validate:"required,min=0,lt=10000000000"` // 箱サイズ100の冷凍便追加配送料(税込)
	HasFreeShipping   bool                  `json:"hasFreeShipping"`                                       // 送料無料オプションの有無
	FreeShippingRates int64                 `json:"freeShippingRates" validate:"min=0,lt=10000000000"`     // 送料無料になる金額(税込)
}

type UpdateShippingRate struct {
	Name            string  `json:"name" validate:"required,max=64"`                       // 配送料金設定名
	Price           int64   `json:"price" validate:"required,min=0,lt=10000000000"`        // 配送料金(税込)
	PrefectureCodes []int32 `json:"prefectureCodes" validate:"required,dive,min=1,max=47"` // 対象都道府県一覧
}

type UpsertShippingRequest struct {
	Box60Rates        []*UpsertShippingRate `json:"box60Rates" validate:"required,dive"`                   // 箱サイズ60の通常便配送料一覧
	Box60Frozen       int64                 `json:"box60Frozen" validate:"required,min=0,lt=10000000000"`  // 箱サイズ60の冷凍便追加配送料(税込)
	Box80Rates        []*UpsertShippingRate `json:"box80Rates" validate:"required,dive"`                   // 箱サイズ80の通常便配送料一覧
	Box80Frozen       int64                 `json:"box80Frozen" validate:"required,min=0,lt=10000000000"`  // 箱サイズ80の冷凍便追加配送料(税込)
	Box100Rates       []*UpsertShippingRate `json:"box100Rates" validate:"required,dive"`                  // 箱サイズ100の通常便配送料一覧
	Box100Frozen      int64                 `json:"box100Frozen" validate:"required,min=0,lt=10000000000"` // 箱サイズ100の冷凍便追加配送料(税込)
	HasFreeShipping   bool                  `json:"hasFreeShipping"`                                       // 送料無料オプションの有無
	FreeShippingRates int64                 `json:"freeShippingRates" validate:"min=0,lt=10000000000"`     // 送料無料になる金額(税込)
}

type UpsertShippingRate struct {
	Name            string  `json:"name" validate:"required,max=64"`                       // 配送料金設定名
	Price           int64   `json:"price" validate:"required,min=0,lt=10000000000"`        // 配送料金(税込)
	PrefectureCodes []int32 `json:"prefectureCodes" validate:"required,dive,min=1,max=47"` // 対象都道府県一覧
}

type UpdateDefaultShippingRequest struct {
	Box60Rates        []*UpdateDefaultShippingRate `json:"box60Rates" validate:"required,dive"`                   // 箱サイズ60の通常便配送料一覧
	Box60Frozen       int64                        `json:"box60Frozen" validate:"required,min=0,lt=10000000000"`  // 箱サイズ60の冷凍便追加配送料(税込)
	Box80Rates        []*UpdateDefaultShippingRate `json:"box80Rates" validate:"required,dive"`                   // 箱サイズ80の通常便配送料一覧
	Box80Frozen       int64                        `json:"box80Frozen" validate:"required,min=0,lt=10000000000"`  // 箱サイズ80の冷凍便追加配送料(税込)
	Box100Rates       []*UpdateDefaultShippingRate `json:"box100Rates" validate:"required,dive"`                  // 箱サイズ100の通常便配送料一覧
	Box100Frozen      int64                        `json:"box100Frozen" validate:"required,min=0,lt=10000000000"` // 箱サイズ100の冷凍便追加配送料(税込)
	HasFreeShipping   bool                         `json:"hasFreeShipping"`                                       // 送料無料オプションの有無
	FreeShippingRates int64                        `json:"freeShippingRates" validate:"min=0,lt=10000000000"`     // 送料無料になる金額(税込)
}

type UpdateDefaultShippingRate struct {
	Name            string  `json:"name" validate:"required,max=64"`                       // 配送料金設定名
	Price           int64   `json:"price" validate:"required,min=0,lt=10000000000"`        // 配送料金(税込)
	PrefectureCodes []int32 `json:"prefectureCodes" validate:"required,dive,min=1,max=47"` // 対象都道府県一覧
}

type ShippingResponse struct {
	Shipping    *Shipping    `json:"shipping"`    // 配送設定情報
	Coordinator *Coordinator `json:"coordinator"` // コーディネータ情報
}

type ShippingsResponse struct {
	Shippings    []*Shipping    `json:"shippings"`    // 配送設定一覧
	Coordinators []*Coordinator `json:"coordinators"` // コーディネータ一覧
	Total        int64          `json:"total"`        // 合計数
}
