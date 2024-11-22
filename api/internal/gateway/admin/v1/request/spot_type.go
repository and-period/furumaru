package request

type CreateSpotTypeRequest struct {
	Name string `json:"name,omitempty"` // スポット種別名
}

type UpdateSpotTypeRequest struct {
	Name string `json:"name,omitempty"` // スポット種別名
}
