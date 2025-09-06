package request

import "time"

type UpdateShopRequest struct {
	Name           string         `json:"name" binding:"required,max=64"`                  // 店舗名
	ProductTypeIDs []string       `json:"productTypeIds" binding:"required,dive,required"` // 取り扱い品目一覧
	BusinessDays   []time.Weekday `json:"businessDays" binding:"max=7,unique"`             // 営業曜日(発送可能日)
}
