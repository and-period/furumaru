package request

type CreateProductReviewRequest struct {
	Rate    int64  `json:"rate,omitempty"`    // 評価
	Title   string `json:"title,omitempty"`   // タイトル
	Comment string `json:"comment,omitempty"` // コメント
}

type UpdateProductReviewRequest struct {
	Rate    int64  `json:"rate,omitempty"`    // 評価
	Title   string `json:"title,omitempty"`   // タイトル
	Comment string `json:"comment,omitempty"` // コメント
}
