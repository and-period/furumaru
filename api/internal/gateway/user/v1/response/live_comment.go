package response

// LiveComment - ライブ配信のコメント
type LiveComment struct {
	UserID       string `json:"userId"`       // ユーザーID
	Username     string `json:"username"`     // ユーザー名
	AccountID    string `json:"accountId"`    // アカウントID
	ThumbnailURL string `json:"thumbnailUrl"` // サムネイルURL
	Comment      string `json:"comment"`      // コメント
	PublishedAt  int64  `json:"publishedAt"`  // 投稿日時
}

type LiveCommentsResponse struct {
	Comments  []*LiveComment `json:"comments"`  // コメント一覧
	NextToken string         `json:"nextToken"` // 次の取得開始位置
}
