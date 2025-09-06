package request

type CreateProductReviewRequest struct {
	Rate    int64  `json:"rate" binding:"min=1,max=5"`          // 評価
	Title   string `json:"title" binding:"required,max=64"`     // タイトル
	Comment string `json:"comment" binding:"required,max=2000"` // コメント
}

type UpdateProductReviewRequest struct {
	Rate    int64  `json:"rate" binding:"min=1,max=5"`          // 評価
	Title   string `json:"title" binding:"required,max=64"`     // タイトル
	Comment string `json:"comment" binding:"required,max=2000"` // コメント
}

type UpsertProductReviewReactionRequest struct {
	ReactionType int32 `json:"reactionType" binding:"required"` // リアクション種別
}
