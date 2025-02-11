package request

import "time"

type UpdateShopRequest struct {
	Name           string         `json:"name,omitempty"`           // 店舗名
	ProductTypeIDs []string       `json:"productTypeIds,omitempty"` // 取り扱い品目一覧
	BusinessDays   []time.Weekday `json:"businessDays,omitempty"`   // 営業曜日(発送可能日)
}
