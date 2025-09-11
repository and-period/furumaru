package types

// ExperienceReviewReactionType - 体験レビューのリアクション種別
type ExperienceReviewReactionType int32

const (
	ExperienceReviewReactionTypeUnknown ExperienceReviewReactionType = 0
	ExperienceReviewReactionTypeLike    ExperienceReviewReactionType = 1 // いいね
	ExperienceReviewReactionTypeDislike ExperienceReviewReactionType = 2 // いまいち
)

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
	ReviewID     string                       `json:"reviewId"`     // 体験レビューID
	ReactionType ExperienceReviewReactionType `json:"reactionType"` // リアクションタイプ
}

type CreateExperienceReviewRequest struct {
	Rate    int64  `json:"rate" validate:"min=1,max=5"`          // 評価
	Title   string `json:"title" validate:"required,max=64"`     // タイトル
	Comment string `json:"comment" validate:"required,max=2000"` // コメント
}

type UpdateExperienceReviewRequest struct {
	Rate    int64  `json:"rate" validate:"min=1,max=5"`          // 評価
	Title   string `json:"title" validate:"required,max=64"`     // タイトル
	Comment string `json:"comment" validate:"required,max=2000"` // コメント
}

type UpsertExperienceReviewReactionRequest struct {
	ReactionType ExperienceReviewReactionType `json:"reactionType" validate:"required"` // リアクション種別
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
