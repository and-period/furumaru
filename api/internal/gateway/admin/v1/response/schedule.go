package response

// Schedule - マルシェ開催情報
type Schedule struct {
	ID              string `json:"id"`              // スケジュールID
	CoordinatorID   string `json:"coordinatorId"`   // コーディネータID
	Status          int32  `json:"status"`          // 開催状況
	Title           string `json:"title"`           // タイトル
	Description     string `json:"description"`     // 説明
	ThumbnailURL    string `json:"thumbnailUrl"`    // サムネイルURL
	ImageURL        string `json:"imageUrl"`        // 蓋絵URL
	OpeningVideoURL string `json:"openingVideoUrl"` // オープニング動画URL
	Public          bool   `json:"public"`          // 公開フラグ
	Approved        bool   `json:"approved"`        // 承認フラグ
	StartAt         int64  `json:"startAt"`         // 配信開始日時
	EndAt           int64  `json:"endAt"`           // 配信終了日時
	CreatedAt       int64  `json:"createdAt"`       // 登録日時
	UpdatedAt       int64  `json:"updatedAt"`       // 更新日時
}

type ScheduleResponse struct {
	Schedule    *Schedule    `json:"schedule"`    // マルシェ開催情報
	Coordinator *Coordinator `json:"coordinator"` // コーディネータ情報
}

type SchedulesResponse struct {
	Schedules    []*Schedule    `json:"schedules"`    // マルシェ開催一覧
	Coordinators []*Coordinator `json:"coordinators"` // コーディネータ一覧
	Total        int64          `json:"total"`        // 合計数
}

type AnalyzeScheduleResponse struct {
	ViewerLogs []*BroadcastViewerLog `json:"viewerLogs"` // 視聴者数ログ
}
