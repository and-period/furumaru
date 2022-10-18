package response

type Schedule struct {
	ID            string `json:"id"`            // スケジュールID
	CoordinatorID string `json:"coordinatorId"` // コーディネーターID
	Title         string `json:"title"`         // タイトル
	Description   string `json:"description"`   // 説明
	ThumbnailURL  string `json:"thumnailURL"`   // サムネイルURL
	StartAt       int64  `json:"startAt"`       // 配信開始日時
	EndAt         int64  `json:"endAt"`         // 配信終了日時
	Canceled      bool   `json:"canceled"`      // 開催中止フラグ
	CreatedAt     int64  `json:"createdAt"`     // 登録日時
	UpdatedAt     int64  `json:"updatedAt"`     // 更新日時
}

type ScheduleResponse struct {
	*Schedule
	Lives []*Live
}
