package response

// Promotion - プロモーション情報
type Promotion struct {
	ID           string `json:"id"`           // プロモーションID
	Title        string `json:"title"`        // タイトル
	Description  string `json:"description"`  // 詳細説明
	DiscountType int32  `json:"discountType"` // 割引計算方法
	DiscountRate int64  `json:"discountRate"` // 割引額(%/円)
	Code         string `json:"code"`         // クーポンコード
	StartAt      int64  `json:"startAt"`      // クーポン使用可能日時(開始)
	EndAt        int64  `json:"endAt"`        // クーポン使用可能日時(終了)
}
