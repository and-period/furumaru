package request

type CreateNotificationRequest struct {
	Title       string       `json:"title,omitempty"`       // タイトル
	Body        string       `json:"body,omitempty"`        // 本文
	Targets     []TargetType `json:"targets,omitempty"`     // 掲載対象一覧
	PublishedAt int64        `json:"publishedAt,omitempty"` // 掲載開始日
	Public      bool         `json:"public,omitempty"`      // 公開フラグ
}

type TargetType int32

const (
	PostTargetUnknown      TargetType = 0 // 対象不明
	PostTargetUsers        TargetType = 1 // ユーザー対象
	PostTargetProducers    TargetType = 2 // 生産者対象
	PostTargetCoordinators TargetType = 3 // コーディネーター対象
)

type UpdateNotificationRequest struct {
	Title       string       `json:"title,omitempty"`       // タイトル
	Body        string       `json:"body,omitempty"`        // 本文
	Targets     []TargetType `json:"targets,omitempty"`     // 掲載対象一覧
	PublishedAt int64        `json:"publishedAt,omitempty"` // 掲載開始日
	Public      bool         `json:"public,omitempty"`      // 公開フラグ
}
