package types

// ProductTag - 商品タグ情報
type ProductTag struct {
	ID        string `json:"id"`        // 商品タグID
	Name      string `json:"name"`      // 商品タグ名
	CreatedAt int64  `json:"createdAt"` // 登録日時
	UpdatedAt int64  `json:"updatedAt"` // 更新日時
}

type CreateProductTagRequest struct {
	Name string `json:"name" validate:"required,max=32"` // 商品タグ名
}

type UpdateProductTagRequest struct {
	Name string `json:"name" validate:"required,max=32"` // 商品タグ名
}

type ProductTagResponse struct {
	ProductTag *ProductTag `json:"productTag"` // 商品タグ情報
}

type ProductTagsResponse struct {
	ProductTags []*ProductTag `json:"productTags"` // 商品タグ一覧
	Total       int64         `json:"total"`       // 商品タグ合計数
}
