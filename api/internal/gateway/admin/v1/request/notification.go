package request

type CreateNotificationRequest struct {
	Type        int32   `json:"type,omitempty"`        // お知らせ種別
	Title       string  `json:"title,omitempty"`       // タイトル
	Body        string  `json:"body,omitempty"`        // 本文
	Note        string  `json:"note,omitempty"`        // 備考
	Targets     []int32 `json:"targets,omitempty"`     // 掲載対象一覧
	PublishedAt int64   `json:"publishedAt,omitempty"` // 掲載開始日
	PromotionID string  `json:"promotionId,omitempty"` // プロモーションID
}

type UpdateNotificationRequest struct {
	Title       string  `json:"title,omitempty"`       // タイトル
	Body        string  `json:"body,omitempty"`        // 本文
	Note        string  `json:"note,omitempty"`        // 備考
	Targets     []int32 `json:"targets,omitempty"`     // 掲載対象一覧
	PublishedAt int64   `json:"publishedAt,omitempty"` // 掲載開始日
}
