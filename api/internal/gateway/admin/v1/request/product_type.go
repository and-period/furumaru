package request

type CreateProductTypeRequest struct {
	Name    string `json:"name" validate:"required,max=32"` // 品目名
	IconURL string `json:"iconUrl" validate:"required,url"` // アイコンURL
}

type UpdateProductTypeRequest struct {
	Name    string `json:"name" validate:"required,max=32"` // 品目名
	IconURL string `json:"iconUrl" validate:"required,url"` // アイコンURL
}
