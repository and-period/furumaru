package request

type CreateNotificationRequest struct {
	Type        int32   `json:"type"`        // お知らせ種別
	Title       string  `json:"title"`       // タイトル
	Body        string  `json:"body"`        // 本文
	Note        string  `json:"note"`        // 備考
	Targets     []int32 `json:"targets"`     // 掲載対象一覧
	PublishedAt int64   `json:"publishedAt"` // 掲載開始日
	PromotionID string  `json:"promotionId"` // プロモーションID
}

type UpdateNotificationRequest struct {
	Title       string  `json:"title"`       // タイトル
	Body        string  `json:"body"`        // 本文
	Note        string  `json:"note"`        // 備考
	Targets     []int32 `json:"targets"`     // 掲載対象一覧
	PublishedAt int64   `json:"publishedAt"` // 掲載開始日
}
