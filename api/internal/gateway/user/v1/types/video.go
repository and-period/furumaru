package types

// Video - オンデマンド配信情報
type Video struct {
	ID            string   `json:"id"`            // オンデマンド動画ID
	CoordinatorID string   `json:"coordinatorId"` // コーディネータID
	ProductIDs    []string `json:"productIds"`    // 商品ID一覧
	ExperienceIDs []string `json:"experienceIds"` // 体験ID一覧
	Title         string   `json:"title"`         // タイトル
	Description   string   `json:"description"`   // 説明
	ThumbnailURL  string   `json:"thumbnailUrl"`  // サムネイルURL
	VideoURL      string   `json:"videoUrl"`      // 動画URL
	PublishedAt   int64    `json:"publishedAt"`   // 公開日時
}

// VideoSummary - オンデマンド配信の概要
type VideoSummary struct {
	ID            string `json:"id"`            // オンデマンド動画ID
	CoordinatorID string `json:"coordinatorId"` // コーディネータID
	Title         string `json:"title"`         // タイトル
	ThumbnailURL  string `json:"thumbnailUrl"`  // サムネイルURL
	PublishedAt   int64  `json:"publishedAt"`   // 公開日時
}

type VideoResponse struct {
	Video       *Video        `json:"video"`       // オンデマンド動画情報
	Coordinator *Coordinator  `json:"coordinator"` // コーディネータ情報
	Products    []*Product    `json:"products"`    // 商品一覧
	Experiences []*Experience `json:"experiences"` // 体験一覧
}

type VideosResponse struct {
	Videos       []*VideoSummary `json:"videos"`       // オンデマンド動画一覧
	Coordinators []*Coordinator  `json:"coordinators"` // コーディネータ一覧
	Total        int64           `json:"total"`        // オンデマンド動画合計数
}
