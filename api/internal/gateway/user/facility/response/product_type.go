package response

// ProductType - 品目情報
type ProductType struct {
	ID         string `json:"id"`         // 品目ID
	Name       string `json:"name"`       // 品目名
	IconURL    string `json:"iconUrl"`    // アイコンURL
	CategoryID string `json:"categoryId"` // 商品種別ID
}
