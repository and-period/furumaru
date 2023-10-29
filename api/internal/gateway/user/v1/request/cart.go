package request

type AddCartItemRequest struct {
	ProductID string `json:"productId,omitempty"` // 商品ID
	Quantity  int64  `json:"quantity,omitempty"`  // 数量
}
