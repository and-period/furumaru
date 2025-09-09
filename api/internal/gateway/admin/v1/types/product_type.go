package types

// ProductType - 品目情報
type ProductType struct {
	ID         string `json:"id"`         // 品目ID
	Name       string `json:"name"`       // 品目名
	IconURL    string `json:"iconUrl"`    // アイコンURL
	CategoryID string `json:"categoryId"` // 商品種別ID
	CreatedAt  int64  `json:"createdAt"`  // 登録日時
	UpdatedAt  int64  `json:"updatedAt"`  // 更新日時
}

type CreateProductTypeRequest struct {
	Name    string `json:"name" validate:"required,max=32"` // 品目名
	IconURL string `json:"iconUrl" validate:"required,url"` // アイコンURL
}

type UpdateProductTypeRequest struct {
	Name    string `json:"name" validate:"required,max=32"` // 品目名
	IconURL string `json:"iconUrl" validate:"required,url"` // アイコンURL
}

type ProductTypeResponse struct {
	ProductType *ProductType `json:"productType"` // 品目情報
	Category    *Category    `json:"category"`    // 商品種別情報
}

type ProductTypesResponse struct {
	ProductTypes []*ProductType `json:"productTypes"` // 品目一覧
	Categories   []*Category    `json:"categories"`   // 商品種別一覧
	Total        int64          `json:"total"`        // 品目合計数
}
