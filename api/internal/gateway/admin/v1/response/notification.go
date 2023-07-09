package response

// Notification - お知らせ情報
type Notification struct {
	ID          string  `json:"id"`          // お知らせID
	Status      int32   `json:"status"`      // お知らせ状態
	Type        int32   `json:"type"`        // お知らせ種別
	Title       string  `json:"title"`       // タイトル
	Body        string  `json:"body"`        // 本文
	Note        string  `json:"note"`        // 備考
	Targets     []int32 `json:"targets"`     // 掲載対象一覧
	PublishedAt int64   `json:"publishedAt"` // 掲載開始日時
	PromotionID string  `json:"promotionId"` // プロモーションID
	CreatedBy   string  `json:"createdBy"`   // 登録者ID
	UpdatedBy   string  `json:"updatedBy"`   // 更新者ID
	CreatedAt   int64   `json:"createdAt"`   // 作成日時
	UpdatedAt   int64   `json:"updatedAt"`   // 更新日時
}

type NotificationResponse struct {
	Notification *Notification `json:"notification"` // お知らせ情報
	Admin        *Admin        `json:"admin"`        // 登録者情報
}

type NotificationsResponse struct {
	Notifications []*Notification `json:"notifications"` // お知らせ一覧
	Admins        []*Admin        `json:"admins"`        // 登録者情報一覧
	Total         int64           `json:"total"`         // お知らせ合計数
}
