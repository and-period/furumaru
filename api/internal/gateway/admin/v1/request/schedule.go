package request

type CreateScheduleRequest struct {
	CoordinatorID   string `json:"coordinatorId"`   // コーディネータID
	Title           string `json:"title"`           // タイトル
	Description     string `json:"description"`     // 説明
	ThumbnailURL    string `json:"thumbnailUrl"`    // サムネイルURL
	ImageURL        string `json:"imageUrl"`        // 蓋絵URL
	OpeningVideoURL string `json:"openingVideoUrl"` // オープニング動画URL
	Public          bool   `json:"public"`          // 公開フラグ
	StartAt         int64  `json:"startAt"`         // 配信開始日時
	EndAt           int64  `json:"endAt"`           // 配信終了日時
}

type UpdateScheduleRequest struct {
	Title           string `json:"title"`           // タイトル
	Description     string `json:"description"`     // 説明
	ThumbnailURL    string `json:"thumbnailUrl"`    // サムネイルURL
	ImageURL        string `json:"imageUrl"`        // 蓋絵URL
	OpeningVideoURL string `json:"openingVideoUrl"` // オープニング動画URL
	StartAt         int64  `json:"startAt"`         // 配信開始日時
	EndAt           int64  `json:"endAt"`           // 配信終了日時
}

type ApproveScheduleRequest struct {
	Approved bool `json:"approved"` // 承認フラグ
}

type PublishScheduleRequest struct {
	Public bool `json:"public"` // 公開フラグ
}
