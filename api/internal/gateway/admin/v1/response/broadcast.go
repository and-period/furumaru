package response

// Broadcast - ライブ配信情報
type Broadcast struct {
	ID         string `json:"id"`         // ライブ配信ID
	ScheduleID string `json:"scheduleId"` // 開催スケジュールID
	Status     int32  `json:"status"`     // ライブ配信状況
	InputURL   string `json:"inputUrl"`   // ライブ配信URL(入力)
	OutputURL  string `json:"outputUrl"`  // ライブ配信URL(出力)
	CreatedAt  int64  `json:"createdAt"`  // 登録日時
	UpdatedAt  int64  `json:"updatedAt"`  // 更新日時
}

type BroadcastResponse struct {
	Broadcast *Broadcast `json:"broadcast"` // ライブ配信情報
}
