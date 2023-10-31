package response

// Live - ライブ配信情報
type Live struct {
	ID         string   `json:"id"`         // ライブ配信ID
	ScheduleID string   `json:"scheduleId"` // マルシェ開催スケジュールID
	ProducerID string   `json:"producerId"` // 生産者ID
	ProductIDs []string `json:"productIds"` // 商品ID一覧
	Comment    string   `json:"comment"`    // コメント
	StartAt    int64    `json:"startAt"`    // ライブ配信開始日時
	EndAt      int64    `json:"endAt"`      // ライブ配信終了日時
}
