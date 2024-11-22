package response

// SpotType - 体験種別情報
type SpotType struct {
	ID   string `json:"id"`   // 体験種別ID
	Name string `json:"name"` // 体験種別名
}

type SpotTypesResponse struct {
	SpotTypes []*SpotType `json:"spotTypes"` // 体験種別一覧
	Total     int64       `json:"total"`     // 体験種別合計数
}
