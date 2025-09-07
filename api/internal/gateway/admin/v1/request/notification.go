package request

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
