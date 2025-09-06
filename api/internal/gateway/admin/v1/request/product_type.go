package request

type CreateProductTypeRequest struct {
	Name    string `json:"name" binding:"required,max=32"` // 品目名
	IconURL string `json:"iconUrl" binding:"required,url"` // アイコンURL
}

type UpdateProductTypeRequest struct {
	Name    string `json:"name" binding:"required,max=32"` // 品目名
	IconURL string `json:"iconUrl" binding:"required,url"` // アイコンURL
}
