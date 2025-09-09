package types

// VideoComment - オンデマンド配信のコメント
type VideoComment struct {
	ID           string `json:"id"`           // コメントID
	UserID       string `json:"userId"`       // ユーザーID
	Username     string `json:"username"`     // ユーザー名
	AccountID    string `json:"accountId"`    // アカウントID
	ThumbnailURL string `json:"thumbnailUrl"` // サムネイルURL
	Comment      string `json:"comment"`      // コメント
	Disabled     bool   `json:"disabled"`     // コメント無効フラグ
	PublishedAt  int64  `json:"publishedAt"`  // 投稿日時
}

type UpdateVideoCommentRequest struct {
	Disabled bool `json:"disabled" validate:""` // コメント無効フラグ
}

type VideoCommentsResponse struct {
	Comments  []*VideoComment `json:"comments"`  // コメント一覧
	NextToken string          `json:"nextToken"` // 次の取得開始位置
}
