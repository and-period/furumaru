package request

type CreateNotificationRequest struct {
	Type        int32   `json:"type" binding:"required"`                     // お知らせ種別
	Title       string  `json:"title" binding:"required,max=128"`            // タイトル
	Body        string  `json:"body" binding:"required,max=2000"`            // 本文
	Note        string  `json:"note" binding:"omitempty,max=2000"`           // 備考
	Targets     []int32 `json:"targets" binding:"min=1,max=4,dive,required"` // 掲載対象一覧
	PublishedAt int64   `json:"publishedAt" binding:"required"`              // 掲載開始日
	PromotionID string  `json:"promotionId" binding:"omitempty"`             // プロモーションID
}

type UpdateNotificationRequest struct {
	Title       string  `json:"title" binding:"required,max=128"`            // タイトル
	Body        string  `json:"body" binding:"required,max=2000"`            // 本文
	Note        string  `json:"note" binding:"omitempty,max=2000"`           // 備考
	Targets     []int32 `json:"targets" binding:"min=1,max=4,dive,required"` // 掲載対象一覧
	PublishedAt int64   `json:"publishedAt" binding:"required"`              // 掲載開始日
}
