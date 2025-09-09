package types

type ContactCategory struct {
	ID        string `json:"id"`        // お問い合わせ種別ID
	Title     string `json:"title"`     // お問い合わせ種別名
	CreatedAt int64  `json:"createdAt"` // 登録日時
	UpdatedAt int64  `json:"updatedAt"` // 更新日時
}

type ContactCategoryResponse struct {
	*ContactCategory
}

type ContactCategoriesResponse struct {
	ContactCategories []*ContactCategory `json:"contactCategories"` // お問い合わせ種別一覧
}
