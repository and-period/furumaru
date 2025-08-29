package request

type CreateProductTypeRequest struct {
	Name    string `json:"name"`    // 品目名
	IconURL string `json:"iconUrl"` // アイコンURL
}

type UpdateProductTypeRequest struct {
	Name    string `json:"name"`    // 品目名
	IconURL string `json:"iconUrl"` // アイコンURL
}
