package request

type CreateProductReviewRequest struct {
	Rate    int64  `json:"rate"`    // 評価
	Title   string `json:"title"`   // タイトル
	Comment string `json:"comment"` // コメント
}

type UpdateProductReviewRequest struct {
	Rate    int64  `json:"rate"`    // 評価
	Title   string `json:"title"`   // タイトル
	Comment string `json:"comment"` // コメント
}

type UpsertProductReviewReactionRequest struct {
	ReactionType int32 `json:"reactionType"` // リアクション種別
}
