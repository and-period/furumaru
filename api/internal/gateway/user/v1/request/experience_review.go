package request

type CreateExperienceReviewRequest struct {
	Rate    int64  `json:"rate,omitempty"`    // 評価
	Title   string `json:"title,omitempty"`   // タイトル
	Comment string `json:"comment,omitempty"` // コメント
}

type UpdateExperienceReviewRequest struct {
	Rate    int64  `json:"rate,omitempty"`    // 評価
	Title   string `json:"title,omitempty"`   // タイトル
	Comment string `json:"comment,omitempty"` // コメント
}

type UpsertExperienceReviewReactionRequest struct {
	ReactionType int32 `json:"reactionType,omitempty"` // リアクション種別
}
