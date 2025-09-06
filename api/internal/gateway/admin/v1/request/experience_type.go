package request

type CreateExperienceTypeRequest struct {
	Name string `json:"name" binding:"required,max=32"` // 体験種別名
}

type UpdateExperienceTypeRequest struct {
	Name string `json:"name" binding:"required,max=32"` // 体験種別名
}
