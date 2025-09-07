package request

type CreateProductReviewRequest struct {
	Rate    int64  `json:"rate" validate:"min=1,max=5"`          // 評価
	Title   string `json:"title" validate:"required,max=64"`     // タイトル
	Comment string `json:"comment" validate:"required,max=2000"` // コメント
}

type UpdateProductReviewRequest struct {
	Rate    int64  `json:"rate" validate:"min=1,max=5"`          // 評価
	Title   string `json:"title" validate:"required,max=64"`     // タイトル
	Comment string `json:"comment" validate:"required,max=2000"` // コメント
}

type UpsertProductReviewReactionRequest struct {
	ReactionType int32 `json:"reactionType" validate:"required"` // リアクション種別
}
