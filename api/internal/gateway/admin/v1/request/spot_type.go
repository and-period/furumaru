package request

type CreateSpotTypeRequest struct {
	Name string `json:"name" binding:"required,max=32"` // スポット種別名
}

type UpdateSpotTypeRequest struct {
	Name string `json:"name" binding:"required,max=32"` // スポット種別名
}
