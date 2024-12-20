package response

// Video - オンデマンド配信情報
type Video struct {
	ID                string   `json:"id"`                // オンデマンド動画ID
	CoordinatorID     string   `json:"coordinatorId"`     // コーディネータID
	ProductIDs        []string `json:"productIds"`        // 商品ID一覧
	ExperienceIDs     []string `json:"experienceIds"`     // 体験ID一覧
	Title             string   `json:"title"`             // タイトル
	Description       string   `json:"description"`       // 説明
	Status            int32    `json:"status"`            // 配信状況
	ThumbnailURL      string   `json:"thumbnailUrl"`      // サムネイルURL
	VideoURL          string   `json:"videoUrl"`          // 動画URL
	Public            bool     `json:"public"`            // 公開設定
	Limited           bool     `json:"limited"`           // 限定公開設定
	DisplayProduct    bool     `json:"displayProduct"`    // 商品への表示設定
	DisplayExperience bool     `json:"displayExperience"` // 体験への表示設定
	PublishedAt       int64    `json:"publishedAt"`       // 公開日時
	CreatedAt         int64    `json:"createdAt"`         // 作成日時
	UpdatedAt         int64    `json:"updatedAt"`         // 更新日時
}

// VideoViewerLog - オンデマンド配信視聴ログ解析情報
type VideoViewerLog struct {
	VideoID string `json:"videoId"` // オンデマンド動画ID
	StartAt int64  `json:"startAt"` // 集計開始日時
	EndAt   int64  `json:"endAt"`   // 集計終了日時
	Total   int64  `json:"total"`   // 合計視聴者数
}

type VideoResponse struct {
	Video       *Video        `json:"video"`       // オンデマンド動画情報
	Coordinator *Coordinator  `json:"coordinator"` // コーディネータ情報
	Products    []*Product    `json:"products"`    // 商品一覧
	Experiences []*Experience `json:"experiences"` // 体験一覧
}

type VideosResponse struct {
	Videos       []*Video       `json:"videos"`       // オンデマンド動画一覧
	Coordinators []*Coordinator `json:"coordinators"` // コーディネータ一覧
	Products     []*Product     `json:"products"`     // 商品一覧
	Experiences  []*Experience  `json:"experiences"`  // 体験一覧
	Total        int64          `json:"total"`        // オンデマンド動画合計数
}

type AnalyzeVideoResponse struct {
	ViewerLogs   []*VideoViewerLog `json:"viewerLogs"`   // 視聴者数ログ
	TotalViewers int64             `json:"totalViewers"` // 合計視聴者数
}
