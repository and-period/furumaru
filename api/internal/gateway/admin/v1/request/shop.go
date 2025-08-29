package request

import "time"

type UpdateShopRequest struct {
	Name           string         `json:"name"`           // 店舗名
	ProductTypeIDs []string       `json:"productTypeIds"` // 取り扱い品目一覧
	BusinessDays   []time.Weekday `json:"businessDays"`   // 営業曜日(発送可能日)
}
