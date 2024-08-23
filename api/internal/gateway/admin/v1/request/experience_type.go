package request

type CreateExperienceTypeRequest struct {
	Name string `json:"name,omitempty"` // 体験種別名
}

type UpdateExperienceTypeRequest struct {
	Name string `json:"name,omitempty"` // 体験種別名
}
