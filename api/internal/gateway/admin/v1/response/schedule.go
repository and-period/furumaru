package response

type Schedule struct {
	ID              string   `json:"id"`              // スケジュールID
	CoordinatorID   string   `json:"coordinatorId"`   // コーディネータID
	CoordinatorName string   `json:"coordinatorName"` // コーディネータ名
	ShippingID      string   `json:"shippingId"`      // 配送設定ID
	ShippingName    string   `json:"shippingName"`    // 配送設定名
	Status          int32    `json:"status"`          // 開催状況
	Title           string   `json:"title"`           // タイトル
	Description     string   `json:"description"`     // 説明
	ThumbnailURL    string   `json:"thumnailUrl"`     // サムネイルURL
	Thumbnails      []*Image `json:"thumbnails"`      // サムネイルURL(リサイズ済み)一覧
	ImageURL        string   `json:"imageUrl"`        // 蓋絵URL
	OpeningVideoURL string   `json:"openingVideoUrl"` // オープニング動画URL
	Public          bool     `json:"public"`          // 公開フラグ
	Approved        bool     `json:"approved"`        // 承認フラグ
	StartAt         int64    `json:"startAt"`         // 配信開始日時
	EndAt           int64    `json:"endAt"`           // 配信終了日時
	CreatedAt       int64    `json:"createdAt"`       // 登録日時
	UpdatedAt       int64    `json:"updatedAt"`       // 更新日時
}

type ScheduleResponse struct {
	*Schedule
	Lives []*Live `json:"lives"` // Deprecated
}

type SchedulesResponse struct {
	Schedules []*Schedule `json:"schedules"` // マルシェ開催情報一覧
	Total     int64       `json:"total"`     // 合計数
}
