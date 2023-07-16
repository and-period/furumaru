package request

type CreateLiveRequest struct {
	ProducerID string   `json:"producerId,omitempty"` // 生産者ID
	ProductIDs []string `json:"productIds,omitempty"` // 商品ID一覧
	Comment    string   `json:"comment,omitempty"`    // コメント
	StartAt    int64    `json:"startAt,omitempty"`    // 配信開始日時
	EndAt      int64    `json:"endAt,omitempty"`      // 配信終了日時
}

type UpdateLiveRequest struct {
	ProductIDs []string `json:"productIds,omitempty"` // 商品ID一覧
	Comment    string   `json:"comment,omitempty"`    // コメント
	StartAt    int64    `json:"startAt,omitempty"`    // 配信開始日時
	EndAt      int64    `json:"endAt,omitempty"`      // 配信終了日時
}
