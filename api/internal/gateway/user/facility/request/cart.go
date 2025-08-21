package request

type AddCartItemRequest struct {
	ProductID string `json:"productId" binding:"required"` // 商品ID
	Quantity  int64  `json:"quantity" binding:"min=1"`     // 数量
}
