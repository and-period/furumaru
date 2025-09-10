package types

// PromotionStatus - プロモーションの状態
type PromotionStatus int32

const (
	PromotionStatusUnknown  PromotionStatus = 0
	PromotionStatusPrivate  PromotionStatus = 1 // 非公開
	PromotionStatusWaiting  PromotionStatus = 2 // 利用開始前
	PromotionStatusEnabled  PromotionStatus = 3 // 利用可能
	PromotionStatusFinished PromotionStatus = 4 // 利用終了
)

// DiscountType - 割引計算方法
type DiscountType int32

const (
	DiscountTypeUnknown      DiscountType = 0
	DiscountTypeAmount       DiscountType = 1 // 固定額(円)
	DiscountTypeRate         DiscountType = 2 // 料率計算(%)
	DiscountTypeFreeShipping DiscountType = 3 // 送料無料
)

// Promotion - プロモーション情報
type Promotion struct {
	ID           string          `json:"id"`           // プロモーションID
	Title        string          `json:"title"`        // タイトル
	Description  string          `json:"description"`  // 詳細説明
	Status       PromotionStatus `json:"status"`       // ステータス
	DiscountType DiscountType    `json:"discountType"` // 割引計算方法
	DiscountRate int64           `json:"discountRate"` // 割引額(%/円)
	Code         string          `json:"code"`         // クーポンコード
	StartAt      int64           `json:"startAt"`      // クーポン使用可能日時(開始)
	EndAt        int64           `json:"endAt"`        // クーポン使用可能日時(終了)
}

type PromotionResponse struct {
	Promotion *Promotion `json:"promotion"` // プロモーション情報
}
