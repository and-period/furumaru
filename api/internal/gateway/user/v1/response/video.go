package response

// VideoSummary - オンデマンド配信情報
type VideoSummary struct {
	ID            string `json:"id"`            // オンデマンド動画ID
	CoordinatorID string `json:"coordinatorId"` // コーディネータID
	Title         string `json:"title"`         // タイトル
	ThumbnailURL  string `json:"thumbnailUrl"`  // サムネイルURL
	PublishedAt   int64  `json:"publishedAt"`   // 公開日時
}
