package request

type CreateProductTagRequest struct {
	Name string `json:"name" validate:"required,max=32"` // 商品タグ名
}

type UpdateProductTagRequest struct {
	Name string `json:"name" validate:"required,max=32"` // 商品タグ名
}
