package response

// ProductType - 品目情報
type ProductType struct {
	ID           string `json:"id"`           // 品目ID
	Name         string `json:"name"`         // 品目名
	CategoryID   string `json:"categoryId"`   // 商品種別ID
	CategoryName string `json:"categoryName"` // 商品種別名
	CreatedAt    int64  `json:"createdAt"`    // 登録日時
	UpdatedAt    int64  `json:"updatedAt"`    // 更新日時
}

type ProductTypeResponse struct {
	*ProductType
}

type ProductTypesResponse struct {
	ProductTypes []*ProductType `json:"productTypes"` // 品目一覧
	Total        int64          `json:"total"`        // 品目合計数
}
