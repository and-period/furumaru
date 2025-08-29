package request

type AddCartItemRequest struct {
	ProductID string `json:"productId"` // 商品ID
	Quantity  int64  `json:"quantity"`  // 数量
}
