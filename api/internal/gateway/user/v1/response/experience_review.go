package response

// ExperienceReview - 体験レビュー
type ExperienceReview struct {
	ID           string `json:"id"`           // 体験レビューID
	ExperienceID string `json:"experienceId"` // 体験ID
	UserID       string `json:"userId"`       // ユーザーID
	Username     string `json:"username"`     // ユーザー名
	AccountID    string `json:"accountId"`    // アカウントID
	ThumbnailURL string `json:"thumbnailUrl"` // サムネイルURL
	Rate         int64  `json:"rate"`         // 評価
	Title        string `json:"title"`        // タイトル
	Comment      string `json:"comment"`      // コメント
	PublishedAt  int64  `json:"publishedAt"`  // 投稿日時
	LikeTotal    int64  `json:"likeTotal"`    // いいね数
	DislikeTotal int64  `json:"dislikeTotal"` // いまいち数
}

// ExperienceReviewReaction - 体験レビューのリアクション
type ExperienceReviewReaction struct {
	ReviewID     string `json:"reviewId"`     // 体験レビューID
	ReactionType int32  `json:"reactionType"` // リアクションタイプ
}

type ExperienceReviewResponse struct {
	Review *ExperienceReview `json:"review"` // 体験レビュー
}

type ExperienceReviewsResponse struct {
	Reviews   []*ExperienceReview `json:"reviews"`   // 体験レビュー一覧
	NextToken string              `json:"nextToken"` // 次の取得開始位置
}

type UserExperienceReviewsResponse struct {
	Reviews   []*ExperienceReview         `json:"reviews"`   // 体験レビュー一覧
	Reactions []*ExperienceReviewReaction `json:"reactions"` // 体験レビューのリアクション一覧
}
