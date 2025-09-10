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

// PromotionTargetType - プロモーションの対象
type PromotionTargetType int32

const (
	PromotionTargetTypeUnknown      PromotionTargetType = 0
	PromotionTargetTypeAllShop      PromotionTargetType = 1 // すべての店舗
	PromotionTargetTypeSpecificShop PromotionTargetType = 2 // 特定の店舗のみ
)

// Promotion - プロモーション情報
type Promotion struct {
	ID           string              `json:"id"`           // プロモーションID
	ShopID       string              `json:"shopId"`       // 店舗ID
	Title        string              `json:"title"`        // タイトル
	Description  string              `json:"description"`  // 詳細説明
	Status       PromotionStatus     `json:"status"`       // ステータス
	Public       bool                `json:"public"`       // 公開フラグ
	TargetType   PromotionTargetType `json:"targetType"`   // 対象商品
	DiscountType DiscountType        `json:"discountType"` // 割引計算方法
	DiscountRate int64               `json:"discountRate"` // 割引額(%/円)
	Code         string              `json:"code"`         // クーポンコード
	StartAt      int64               `json:"startAt"`      // クーポン使用可能日時(開始)
	EndAt        int64               `json:"endAt"`        // クーポン使用可能日時(終了)
	UsedCount    int64               `json:"usedCount"`    // 使用回数
	UsedAmount   int64               `json:"usedAmount"`   // 使用による割引合計額
	CreatedAt    int64               `json:"createdAt"`    // 登録日時
	UpdatedAt    int64               `json:"updatedAt"`    // 更新日時
}

type CreatePromotionRequest struct {
	Title        string `json:"title" validate:"required,max=64"`
	Description  string `json:"description" validate:"required,max=2000"`
	Public       bool   `json:"public" validate:""`
	DiscountType int32  `json:"discountType" validate:"required"`
	DiscountRate int64  `json:"discountRate" validate:"min=0"`
	Code         string `json:"code" validate:"len=8"`
	StartAt      int64  `json:"startAt" validate:"required"`
	EndAt        int64  `json:"endAt" validate:"required,gtfield=StartAt"`
}

type UpdatePromotionRequest struct {
	Title        string `json:"title" validate:"required,max=64"`
	Description  string `json:"description" validate:"required,max=2000"`
	Public       bool   `json:"public" validate:""`
	DiscountType int32  `json:"discountType" validate:"required"`
	DiscountRate int64  `json:"discountRate" validate:"min=0"`
	Code         string `json:"code" validate:"len=8"`
	StartAt      int64  `json:"startAt" validate:"required"`
	EndAt        int64  `json:"endAt" validate:"required,gtfield=StartAt"`
}

type PromotionResponse struct {
	Promotion *Promotion `json:"promotion"` // プロモーション情報
	Shop      *Shop      `json:"shop"`      // 店舗情報
}

type PromotionsResponse struct {
	Promotions []*Promotion `json:"promotions"` // プロモーション情報一覧
	Shops      []*Shop      `json:"shops"`      // 店舗情報一覧
	Total      int64        `json:"total"`      // プロモーション合計数
}
