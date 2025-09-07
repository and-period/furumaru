package request

type CreateScheduleRequest struct {
	CoordinatorID   string `json:"coordinatorId" validate:"required"`         // コーディネータID
	Title           string `json:"title" validate:"required,max=64"`          // タイトル
	Description     string `json:"description" validate:"required,max=2000"`  // 説明
	ThumbnailURL    string `json:"thumbnailUrl" validate:"omitempty,url"`     // サムネイルURL
	ImageURL        string `json:"imageUrl" validate:"omitempty,url"`         // 蓋絵URL
	OpeningVideoURL string `json:"openingVideoUrl" validate:"omitempty,url"`  // オープニング動画URL
	Public          bool   `json:"public" validate:""`                        // 公開フラグ
	StartAt         int64  `json:"startAt" validate:"required"`               // 配信開始日時
	EndAt           int64  `json:"endAt" validate:"required,gtfield=StartAt"` // 配信終了日時
}

type UpdateScheduleRequest struct {
	Title           string `json:"title" validate:"required,max=64"`          // タイトル
	Description     string `json:"description" validate:"required,max=2000"`  // 説明
	ThumbnailURL    string `json:"thumbnailUrl" validate:"omitempty,url"`     // サムネイルURL
	ImageURL        string `json:"imageUrl" validate:"omitempty,url"`         // 蓋絵URL
	OpeningVideoURL string `json:"openingVideoUrl" validate:"omitempty,url"`  // オープニング動画URL
	StartAt         int64  `json:"startAt" validate:"required"`               // 配信開始日時
	EndAt           int64  `json:"endAt" validate:"required,gtfield=StartAt"` // 配信終了日時
}

type ApproveScheduleRequest struct {
	Approved bool `json:"approved" validate:""` // 承認フラグ
}

type PublishScheduleRequest struct {
	Public bool `json:"public" validate:""` // 公開フラグ
}
