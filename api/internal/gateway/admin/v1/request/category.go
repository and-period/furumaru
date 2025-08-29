package request

type CreateCategoryRequest struct {
	Name string `json:"name"` // 商品種別名
}

type UpdateCategoryRequest struct {
	Name string `json:"name"` // 商品種別名
}
