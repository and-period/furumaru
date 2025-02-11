package response

import "time"

// Shop - 店舗情報
type Shop struct {
	ID             string         `json:"id"`             // 店舗ID
	Name           string         `json:"name"`           // 店舗名
	CoordinatorID  string         `json:"coordinatorId"`  // コーディネータID
	ProducerIDs    []string       `json:"producerIds"`    // 生産者ID一覧
	ProductTypeIDs []string       `json:"productTypeIds"` // 取り扱い品目一覧
	BusinessDays   []time.Weekday `json:"businessDays"`   // 営業曜日(発送可能日)
	CreatedAt      int64          `json:"createdAt"`      // 登録日時
	UpdatedAt      int64          `json:"updatedAt"`      // 更新日時
}

type ShopResponse struct {
	Shop         *Shop          `json:"shop"`         // 店舗情報
	Coordinator  *Coordinator   `json:"coordinator"`  // コーディネータ情報
	Producers    []*Producer    `json:"producers"`    // 生産者一覧
	ProductTypes []*ProductType `json:"productTypes"` // 品目一覧
}
