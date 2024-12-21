package response

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
}

type ProductReviewResponse struct {
	Review *ProductReview `json:"review"` // 商品レビュー
}

type ProductReviewsResponse struct {
	Reviews   []*ProductReview `json:"reviews"`   // 商品レビュー一覧
	NextToken string           `json:"nextToken"` // 次の取得開始位置
}
