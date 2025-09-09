package types

// Category - 商品種別情報
type Category struct {
	ID        string `json:"id"`        // 商品種別ID
	Name      string `json:"name"`      // 商品種別名
	CreatedAt int64  `json:"createdAt"` // 登録日時
	UpdatedAt int64  `json:"updatedAt"` // 更新日時
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,max=32"` // 商品種別名
}

type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,max=32"` // 商品種別名
}

type CategoryResponse struct {
	Category *Category `json:"category"` // 商品種別情報
}

type CategoriesResponse struct {
	Categories []*Category `json:"categories"` // 商品種別一覧
	Total      int64       `json:"total"`      // 商品種別合計数
}
