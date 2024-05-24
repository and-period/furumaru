package response

// Broadcast - ライブ配信情報
type Broadcast struct {
	ID               string `json:"id"`               // ライブ配信ID
	ScheduleID       string `json:"scheduleId"`       // 開催スケジュールID
	Status           int32  `json:"status"`           // ライブ配信状況
	InputURL         string `json:"inputUrl"`         // ライブ配信URL(入力)
	OutputURL        string `json:"outputUrl"`        // ライブ配信URL(出力)
	ArchiveURL       string `json:"archiveUrl"`       // オンデマンド配信URL
	YoutubeAccount   string `json:"youtubeAccount"`   // Youtubeアカウント
	YoutubeViewerURL string `json:"youtubeViewerUrl"` // Youtube視聴画面URL
	YoutubeAdminURL  string `json:"youtubeAdminUrl"`  // Youtube管理画面URL
	CreatedAt        int64  `json:"createdAt"`        // 登録日時
	UpdatedAt        int64  `json:"updatedAt"`        // 更新日時
}

// GuestBroadcast - ゲスト用ライブ配信情報
type GuestBroadcast struct {
	Title             string `json:"title"`             // ライブ配信タイトル
	Description       string `json:"description"`       // ライブ配信説明
	CoordinatorMarche string `json:"coordinatorMarche"` // ライブ配信担当者(マルシェ)
	CoordinatorName   string `json:"coordinatorName"`   // ライブ配信担当者(名前)
	StartAt           int64  `json:"startAt"`           // ライブ配信開始日時
	EndAt             int64  `json:"endAt"`             // ライブ配信終了日時
}

type BroadcastResponse struct {
	Broadcast *Broadcast `json:"broadcast"` // ライブ配信情報
}

type GuestBroadcastResponse struct {
	Broadcast *GuestBroadcast `json:"broadcast"` // ゲスト用ライブ配信情報
}

type AuthYoutubeBroadcastResponse struct {
	URL string `json:"url"` // 認証URL
}
