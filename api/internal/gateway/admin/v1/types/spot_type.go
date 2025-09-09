package types

// SpotType - スポット種別情報
type SpotType struct {
	ID        string `json:"id"`        // スポット種別ID
	Name      string `json:"name"`      // スポット種別名
	CreatedAt int64  `json:"createdAt"` // 登録日時
	UpdatedAt int64  `json:"updatedAt"` // 更新日時
}

type CreateSpotTypeRequest struct {
	Name string `json:"name" validate:"required,max=32"` // スポット種別名
}

type UpdateSpotTypeRequest struct {
	Name string `json:"name" validate:"required,max=32"` // スポット種別名
}

type SpotTypeResponse struct {
	SpotType *SpotType `json:"spotType"` // 体験種別情報
}

type SpotTypesResponse struct {
	SpotTypes []*SpotType `json:"spotTypes"` // 体験種別一覧
	Total     int64       `json:"total"`     // 体験種別合計数
}
