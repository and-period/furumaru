package response

// 配信情報
type Live struct {
	ID           string     `json:"id"`           // ライブ配信ID
	ScheduleID   string     `json:"scheduleId"`   // スケジュールID
	Title        string     `json:"title"`        // タイトル
	Description  string     `json:"description"`  // 説明
	ProducerID   string     `json:"producerId"`   // 生産者ID
	ProducerName string     `json:"producerName"` // 生産者名
	StartAt      int64      `json:"startAt"`      // 配信開始日時
	EndAt        int64      `json:"endAt"`        // 配信終了日時
	Published    bool       `json:"published"`    // 配信公開フラグ
	Canceled     bool       `json:"canceled"`     // 配信中止
	Status       int32      `json:"status"`       // 配信ステータス
	Products     []*Product `json:"products"`     // 商品一覧
	CreatedAt    int64      `json:"createdAt"`    // 作成日時
	UpdatedAt    int64      `json:"UpdatedAt"`    // 作成日時
}

type LiveResponse struct {
	*Live
}

type LivesResponse struct {
	Lives []*Live `json:"lives"` // 配信情報一覧
}
