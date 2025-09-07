package request

type AddCartItemRequest struct {
	ProductID string `json:"productId" validate:"required"` // 商品ID
	Quantity  int64  `json:"quantity" validate:"min=1"`     // 数量
}
