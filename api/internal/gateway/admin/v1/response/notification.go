package response

// Notification - お知らせ情報
type Notification struct {
	ID          string       `json:"id"`          // お知らせID
	CreatedBy   string       `json:"createdBy"`   // 登録者ID
	CreatorName string       `json:"creatorName"` // 登録者名
	UpdatedBy   string       `json:"updatedBy"`   // 更新者ID
	Title       string       `json:"title"`       // タイトル
	Body        string       `json:"body"`        // 本文
	Targets     []TargetType `json:"targets"`     // 掲載対象一覧
	PublishedAt int64        `json:"publishedAt"` // 掲載開始日時
	Public      bool         `json:"public"`      // 公開フラグ
	CreatedAt   int64        `json:"createdAt"`   // 作成日時
	UpdatedAt   int64        `json:"updatedAt"`   // 更新日時
}

// 掲載対象
type TargetType int32

const (
	PostTargetUnknown      TargetType = 0 // 対象不明
	PostTargetUsers        TargetType = 1 // ユーザー対象
	PostTargetProducers    TargetType = 2 // 生産者対象
	PostTargetCoordinators TargetType = 3 // コーディネーター対象
)

type NotificationResponse struct {
	*Notification
}

type NotificationsResponse struct {
	Notifications []*Notification `json:"notifications"` // お知らせ一覧
	Total         int64           `json:"total"`         // お知らせ合計数
}
