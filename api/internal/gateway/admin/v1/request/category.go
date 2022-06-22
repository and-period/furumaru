package request

type CreateCategoryRequest struct {
	Name string `json:"name,omitempty"` // 商品種別名
}

type UpdateCategoryRequest struct {
	Name string `json:"name,omitempty"` // 商品種別名
}
