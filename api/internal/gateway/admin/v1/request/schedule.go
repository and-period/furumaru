package request

type CreateScheduleRequest struct {
	CoordinatorID string                `json:"coordinatorId,omitempty"` // コーディネーターID
	ShippingID    string                `json:"shippingId,omitempty"`    // 配送設定ID
	Title         string                `json:"title,omitempty"`         // タイトル
	Description   string                `json:"description,omitempty"`   // 説明
	ThumbnailURL  string                `json:"thumnailUrl,omitempty"`   // サムネイルURL
	StartAt       int64                 `json:"startAt,omitempty"`       // 配信開始日時
	EndAt         int64                 `json:"endAt,omitempty"`         // 配信終了日時
	Lives         []*CreateScheduleLive `json:"lives,omitempty"`         // ライブ配信一覧
}

type CreateScheduleLive struct {
	Title       string   `json:"title,omitempty"`       // タイトル
	Description string   `json:"description,omitempty"` // 説明
	ProducerID  string   `json:"producerId,omitempty"`  // 生産者ID
	ProductIDs  []string `json:"productIds,omitempty"`  // 商品ID一覧
	StartAt     int64    `json:"startAt,omitempty"`     // 配信開始日時
	EndAt       int64    `json:"endAt,omitempty"`       // 配信終了日時
}
