package request

type CreateProductTagRequest struct {
	Name string `json:"name,omitempty"` // 商品タグ名
}

type UpdateProductTagRequest struct {
	Name string `json:"name,omitempty"` // 商品タグ名
}
