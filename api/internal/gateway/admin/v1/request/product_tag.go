package request

type CreateProductTagRequest struct {
	Name string `json:"name"` // 商品タグ名
}

type UpdateProductTagRequest struct {
	Name string `json:"name"` // 商品タグ名
}
