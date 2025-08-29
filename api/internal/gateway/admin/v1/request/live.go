package request

type CreateLiveRequest struct {
	ProducerID string   `json:"producerId"` // 生産者ID
	ProductIDs []string `json:"productIds"` // 商品ID一覧
	Comment    string   `json:"comment"`    // コメント
	StartAt    int64    `json:"startAt"`    // 配信開始日時
	EndAt      int64    `json:"endAt"`      // 配信終了日時
}

type UpdateLiveRequest struct {
	ProductIDs []string `json:"productIds"` // 商品ID一覧
	Comment    string   `json:"comment"`    // コメント
	StartAt    int64    `json:"startAt"`    // 配信開始日時
	EndAt      int64    `json:"endAt"`      // 配信終了日時
}
