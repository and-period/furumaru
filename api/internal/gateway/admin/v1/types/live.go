package types

// Live - ライブ配信情報
type Live struct {
	ID         string   `json:"id"`         // ライブ配信ID
	ScheduleID string   `json:"scheduleId"` // マルシェ開催スケジュールID
	ProducerID string   `json:"producerId"` // 生産者ID
	ProductIDs []string `json:"productIds"` // 商品ID一覧
	Comment    string   `json:"comment"`    // コメント
	StartAt    int64    `json:"startAt"`    // ライブ配信開始日時
	EndAt      int64    `json:"endAt"`      // ライブ配信終了日時
	CreatedAt  int64    `json:"createdAt"`  // 登録日時
	UpdatedAt  int64    `json:"updatedAt"`  // 更新日時
}

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

type LiveResponse struct {
	Live     *Live      `json:"live"`     // ライブ配信情報
	Producer *Producer  `json:"producer"` // 生産者情報
	Products []*Product `json:"products"` // 商品一覧
}

type LivesResponse struct {
	Lives     []*Live     `json:"lives"`     // ライブ配信一覧
	Producers []*Producer `json:"producers"` // 生産者一覧
	Products  []*Product  `json:"products"`  // 商品一覧
	Total     int64       `json:"total"`     // 合計数
}
