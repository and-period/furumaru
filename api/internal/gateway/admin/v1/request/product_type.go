package request

type CreateProductTypeRequest struct {
	Name string `json:"name,omitempty"` // 品目名
}

type UpdateProductTypeRequest struct {
	Name string `json:"name,omitempty"` // 品目名
}
