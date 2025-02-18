package response

// Promotion - プロモーション情報
type Promotion struct {
	ID           string `json:"id"`           // プロモーションID
	ShopID       string `json:"shopId"`       // 店舗ID
	Title        string `json:"title"`        // タイトル
	Description  string `json:"description"`  // 詳細説明
	Status       int32  `json:"status"`       // ステータス
	Public       bool   `json:"public"`       // 公開フラグ
	TargetType   int32  `json:"targetType"`   // 対象商品
	DiscountType int32  `json:"discountType"` // 割引計算方法
	DiscountRate int64  `json:"discountRate"` // 割引額(%/円)
	Code         string `json:"code"`         // クーポンコード
	StartAt      int64  `json:"startAt"`      // クーポン使用可能日時(開始)
	EndAt        int64  `json:"endAt"`        // クーポン使用可能日時(終了)
	UsedCount    int64  `json:"usedCount"`    // 使用回数
	UsedAmount   int64  `json:"usedAmount"`   // 使用による割引合計額
	CreatedAt    int64  `json:"createdAt"`    // 登録日時
	UpdatedAt    int64  `json:"updatedAt"`    // 更新日時
}

type PromotionResponse struct {
	Promotion *Promotion `json:"promotion"` // プロモーション情報
}

type PromotionsResponse struct {
	Promotions []*Promotion `json:"promotions"` // プロモーション情報一覧
	Total      int64        `json:"total"`      // プロモーション合計数
}
