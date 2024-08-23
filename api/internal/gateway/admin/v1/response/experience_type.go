package response

// ExperienceType - 体験種別情報
type ExperienceType struct {
	ID        string `json:"id"`        // 体験種別ID
	Name      string `json:"name"`      // 体験種別名
	CreatedAt int64  `json:"createdAt"` // 登録日時
	UpdatedAt int64  `json:"updatedAt"` // 更新日時
}

type ExperienceTypeResponse struct {
	ExperienceType *ExperienceType `json:"experienceType"` // 体験種別情報
}

type ExperienceTypesResponse struct {
	ExperienceTypes []*ExperienceType `json:"experienceTypes"` // 体験種別一覧
	Total           int64             `json:"total"`           // 体験種別合計数
}
