package request

type CreateProductTypeRequest struct {
	Name    string `json:"name,omitempty"`    // 品目名
	IconURL string `json:"iconUrl,omitempty"` // アイコンURL
}

type UpdateProductTypeRequest struct {
	Name    string `json:"name,omitempty"`    // 品目名
	IconURL string `json:"iconUrl,omitempty"` // アイコンURL
}
