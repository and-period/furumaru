package request

type CreateScheduleRequest struct {
	CoordinatorID   string `json:"coordinatorId" binding:"required"`         // コーディネータID
	Title           string `json:"title" binding:"required,max=64"`          // タイトル
	Description     string `json:"description" binding:"required,max=2000"`  // 説明
	ThumbnailURL    string `json:"thumbnailUrl" binding:"omitempty,url"`     // サムネイルURL
	ImageURL        string `json:"imageUrl" binding:"omitempty,url"`         // 蓋絵URL
	OpeningVideoURL string `json:"openingVideoUrl" binding:"omitempty,url"`  // オープニング動画URL
	Public          bool   `json:"public" binding:""`                        // 公開フラグ
	StartAt         int64  `json:"startAt" binding:"required"`               // 配信開始日時
	EndAt           int64  `json:"endAt" binding:"required,gtfield=StartAt"` // 配信終了日時
}

type UpdateScheduleRequest struct {
	Title           string `json:"title" binding:"required,max=64"`          // タイトル
	Description     string `json:"description" binding:"required,max=2000"`  // 説明
	ThumbnailURL    string `json:"thumbnailUrl" binding:"omitempty,url"`     // サムネイルURL
	ImageURL        string `json:"imageUrl" binding:"omitempty,url"`         // 蓋絵URL
	OpeningVideoURL string `json:"openingVideoUrl" binding:"omitempty,url"`  // オープニング動画URL
	StartAt         int64  `json:"startAt" binding:"required"`               // 配信開始日時
	EndAt           int64  `json:"endAt" binding:"required,gtfield=StartAt"` // 配信終了日時
}

type ApproveScheduleRequest struct {
	Approved bool `json:"approved" binding:""` // 承認フラグ
}

type PublishScheduleRequest struct {
	Public bool `json:"public" binding:""` // 公開フラグ
}
