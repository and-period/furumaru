package request

type CreateExperienceReviewRequest struct {
	Rate    int64  `json:"rate"`    // 評価
	Title   string `json:"title"`   // タイトル
	Comment string `json:"comment"` // コメント
}

type UpdateExperienceReviewRequest struct {
	Rate    int64  `json:"rate"`    // 評価
	Title   string `json:"title"`   // タイトル
	Comment string `json:"comment"` // コメント
}

type UpsertExperienceReviewReactionRequest struct {
	ReactionType int32 `json:"reactionType"` // リアクション種別
}
