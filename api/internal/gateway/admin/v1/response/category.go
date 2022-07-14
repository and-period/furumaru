package response

// Category - 商品種別情報
type Category struct {
	ID        string `json:"id"`        // 商品種別ID
	Name      string `json:"name"`      // 商品種別名
	CreatedAt int64  `json:"createdAt"` // 登録日時
	UpdatedAt int64  `json:"updatedAt"` // 更新日時
}

type CategoryResponse struct {
	*Category
}

type CategoriesResponse struct {
	Categories []*Category `json:"categories"` // 商品種別一覧
	Total      int64       `json:"total"`      // 商品種別合計数
}
