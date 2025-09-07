package request

type CreateLiveRequest struct {
	ProducerID string   `json:"producerId" validate:"required"`               // 生産者ID
	ProductIDs []string `json:"productIds" validate:"required,dive,required"` // 商品ID一覧
	Comment    string   `json:"comment" validate:"required,max=2000"`         // コメント
	StartAt    int64    `json:"startAt" validate:"required"`                  // 配信開始日時
	EndAt      int64    `json:"endAt" validate:"required,gtfield=StartAt"`    // 配信終了日時
}

type UpdateLiveRequest struct {
	ProductIDs []string `json:"productIds" validate:"required,dive,required"` // 商品ID一覧
	Comment    string   `json:"comment" validate:"required,max=2000"`         // コメント
	StartAt    int64    `json:"startAt" validate:"required"`                  // 配信開始日時
	EndAt      int64    `json:"endAt" validate:"required,gtfield=StartAt"`    // 配信終了日時
}
