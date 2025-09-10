package types

// BroadcastStatus - ライブ配信状況
type BroadcastStatus int32

const (
	BroadcastStatusUnknown  BroadcastStatus = 0
	BroadcastStatusDisabled BroadcastStatus = 1 // リソース未作成
	BroadcastStatusWaiting  BroadcastStatus = 2 // リソース作成/削除中
	BroadcastStatusIdle     BroadcastStatus = 3 // 停止中
	BroadcastStatusActive   BroadcastStatus = 4 // 配信中
)

// BroadcastViewerLogInterval - ライブ配信視聴ログ取得間隔
type BroadcastViewerLogInterval string

const (
	BroadcastViewerLogIntervalSecond BroadcastViewerLogInterval = "second"
	BroadcastViewerLogIntervalMinute BroadcastViewerLogInterval = "minute"
	BroadcastViewerLogIntervalHour   BroadcastViewerLogInterval = "hour"
)

// Broadcast - ライブ配信情報
type Broadcast struct {
	ID               string          `json:"id"`               // ライブ配信ID
	ScheduleID       string          `json:"scheduleId"`       // 開催スケジュールID
	Status           BroadcastStatus `json:"status"`           // ライブ配信状況
	InputURL         string          `json:"inputUrl"`         // ライブ配信URL(入力)
	OutputURL        string          `json:"outputUrl"`        // ライブ配信URL(出力)
	ArchiveURL       string          `json:"archiveUrl"`       // オンデマンド配信URL
	YoutubeAccount   string          `json:"youtubeAccount"`   // Youtubeアカウント
	YoutubeViewerURL string          `json:"youtubeViewerUrl"` // Youtube視聴画面URL
	YoutubeAdminURL  string          `json:"youtubeAdminUrl"`  // Youtube管理画面URL
	CreatedAt        int64           `json:"createdAt"`        // 登録日時
	UpdatedAt        int64           `json:"updatedAt"`        // 更新日時
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

// BroadcastViewerLog - ライブ配信視聴ログ解析情報
type BroadcastViewerLog struct {
	BroadcastID string `json:"broadcastId"` // ライブ配信ID
	StartAt     int64  `json:"startAt"`     // 集計開始日時
	EndAt       int64  `json:"endAt"`       // 集計終了日時
	Total       int64  `json:"total"`       // 合計視聴者数
}

type UpdateBroadcastArchiveRequest struct {
	ArchiveURL string `json:"archiveUrl" validate:"required,url"` // アーカイブ動画URL
}

type ActivateBroadcastMP4Request struct {
	InputURL string `json:"inputUrl" validate:"required,url"` // 配信動画URL
}

type AuthYoutubeBroadcastRequest struct {
	YoutubeHandle string `json:"youtubeHandle" validate:"required"` // 連携先Youtubeアカウント
}

type CallbackAuthYoutubeBroadcastRequest struct {
	State    string `json:"state" validate:"required"`    // Google認証時に取得したstate
	AuthCode string `json:"authCode" validate:"required"` // Google認証時に取得したcode
}

type CreateYoutubeBroadcastRequest struct {
	Title       string `json:"title" validate:"required,max=100"`         // ライブ配信タイトル
	Description string `json:"description" validate:"omitempty,max=1000"` // ライブ配信説明
	Public      bool   `json:"public" validate:""`                        // 公開設定
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
