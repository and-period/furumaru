package request

type CreateScheduleRequest struct {
	CoordinatorID   string `json:"coordinatorId,omitempty"`   // コーディネータID
	Title           string `json:"title,omitempty"`           // タイトル
	Description     string `json:"description,omitempty"`     // 説明
	ThumbnailURL    string `json:"thumbnailUrl,omitempty"`    // サムネイルURL
	ImageURL        string `json:"imageUrl,omitempty"`        // 蓋絵URL
	OpeningVideoURL string `json:"openingVideoUrl,omitempty"` // オープニング動画URL
	Public          bool   `json:"public,omitempty"`          // 公開フラグ
	StartAt         int64  `json:"startAt,omitempty"`         // 配信開始日時
	EndAt           int64  `json:"endAt,omitempty"`           // 配信終了日時
}

type UpdateScheduleRequest struct {
	Title           string `json:"title,omitempty"`           // タイトル
	Description     string `json:"description,omitempty"`     // 説明
	ThumbnailURL    string `json:"thumbnailUrl,omitempty"`    // サムネイルURL
	ImageURL        string `json:"imageUrl,omitempty"`        // 蓋絵URL
	OpeningVideoURL string `json:"openingVideoUrl,omitempty"` // オープニング動画URL
	StartAt         int64  `json:"startAt,omitempty"`         // 配信開始日時
	EndAt           int64  `json:"endAt,omitempty"`           // 配信終了日時
}

type ApproveScheduleRequest struct {
	Approved bool `json:"approved,omitempty"` // 承認フラグ
}

type PublishScheduleRequest struct {
	Public bool `json:"public,omitempty"` // 公開フラグ
}
