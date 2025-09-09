package types

// ProductReview - 商品レビュー
type ProductReview struct {
	ID           string `json:"id"`           // 商品レビューID
	ProductID    string `json:"productId"`    // 商品ID
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

// ProductReviewReaction - 商品レビューのリアクション
type ProductReviewReaction struct {
	ReviewID     string `json:"reviewId"`     // 商品レビューID
	ReactionType int32  `json:"reactionType"` // リアクションタイプ
}

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

type ProductReviewResponse struct {
	Review *ProductReview `json:"review"` // 商品レビュー
}

type ProductReviewsResponse struct {
	Reviews   []*ProductReview `json:"reviews"`   // 商品レビュー一覧
	NextToken string           `json:"nextToken"` // 次の取得開始位置
}

type UserProductReviewsResponse struct {
	Reviews   []*ProductReview         `json:"reviews"`   // 商品レビュー一覧
	Reactions []*ProductReviewReaction `json:"reactions"` // 商品レビューのリアクション一覧
}
