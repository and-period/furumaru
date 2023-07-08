package request

type CreateScheduleRequest struct {
	CoordinatorID   string `json:"coordinatorId,omitempty"`   // コーディネータID
	ShippingID      string `json:"shippingId,omitempty"`      // 配送設定ID
	Title           string `json:"title,omitempty"`           // タイトル
	Description     string `json:"description,omitempty"`     // 説明
	ThumbnailURL    string `json:"thumnailUrl,omitempty"`     // サムネイルURL
	ImageURL        string `json:"imageUrl,omitempty"`        // 蓋絵URL
	OpeningVideoURL string `json:"openingVideoUrl,omitempty"` // オープニング動画URL
	Public          bool   `json:"public,omitempty"`          // 公開フラグ
	StartAt         int64  `json:"startAt,omitempty"`         // 配信開始日時
	EndAt           int64  `json:"endAt,omitempty"`           // 配信終了日時
}
