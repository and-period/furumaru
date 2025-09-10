package types

// お知らせ種別
type NotificationType int32

const (
	NotificationTypeUnknown   NotificationType = 0
	NotificationTypeOther     NotificationType = 1 // その他
	NotificationTypeSystem    NotificationType = 2 // システム関連
	NotificationTypeLive      NotificationType = 3 // ライブ関連
	NotificationTypePromotion NotificationType = 4 // セール関連
)

// お知らせ通知先
type NotificationTarget int32

const (
	NotificationTargetUnknown        NotificationTarget = 0
	NotificationTargetUsers          NotificationTarget = 1 // ユーザー
	NotificationTargetProducers      NotificationTarget = 2 // 生産者
	NotificationTargetCoordinators   NotificationTarget = 3 // コーディネータ
	NotificationTargetAdministrators NotificationTarget = 4 // 管理者
)

// お知らせ状態
type NotificationStatus int32

const (
	NotificationStatusUnknown  NotificationStatus = 0
	NotificationStatusWaiting  NotificationStatus = 1 // 投稿前
	NotificationStatusNotified NotificationStatus = 2 // 投稿済み
)

// Notification - お知らせ情報
type Notification struct {
	ID          string               `json:"id"`          // お知らせID
	Status      NotificationStatus   `json:"status"`      // お知らせ状態
	Type        NotificationType     `json:"type"`        // お知らせ種別
	Title       string               `json:"title"`       // タイトル
	Body        string               `json:"body"`        // 本文
	Note        string               `json:"note"`        // 備考
	Targets     []NotificationTarget `json:"targets"`     // 掲載対象一覧
	PublishedAt int64                `json:"publishedAt"` // 掲載開始日時
	PromotionID string               `json:"promotionId"` // プロモーションID
	CreatedBy   string               `json:"createdBy"`   // 登録者ID
	UpdatedBy   string               `json:"updatedBy"`   // 更新者ID
	CreatedAt   int64                `json:"createdAt"`   // 作成日時
	UpdatedAt   int64                `json:"updatedAt"`   // 更新日時
}

type CreateNotificationRequest struct {
	Type        int32   `json:"type" validate:"required"`                     // お知らせ種別
	Title       string  `json:"title" validate:"required,max=128"`            // タイトル
	Body        string  `json:"body" validate:"required,max=2000"`            // 本文
	Note        string  `json:"note" validate:"omitempty,max=2000"`           // 備考
	Targets     []int32 `json:"targets" validate:"min=1,max=4,dive,required"` // 掲載対象一覧
	PublishedAt int64   `json:"publishedAt" validate:"required"`              // 掲載開始日
	PromotionID string  `json:"promotionId" validate:"omitempty"`             // プロモーションID
}

type UpdateNotificationRequest struct {
	Title       string  `json:"title" validate:"required,max=128"`            // タイトル
	Body        string  `json:"body" validate:"required,max=2000"`            // 本文
	Note        string  `json:"note" validate:"omitempty,max=2000"`           // 備考
	Targets     []int32 `json:"targets" validate:"min=1,max=4,dive,required"` // 掲載対象一覧
	PublishedAt int64   `json:"publishedAt" validate:"required"`              // 掲載開始日
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
