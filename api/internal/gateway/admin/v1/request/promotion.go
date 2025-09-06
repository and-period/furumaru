package request

type CreatePromotionRequest struct {
	Title        string `json:"title" binding:"required,max=64"`
	Description  string `json:"description" binding:"required,max=2000"`
	Public       bool   `json:"public" binding:""`
	DiscountType int32  `json:"discountType" binding:"required"`
	DiscountRate int64  `json:"discountRate" binding:"min=0"`
	Code         string `json:"code" binding:"len=8"`
	StartAt      int64  `json:"startAt" binding:"required"`
	EndAt        int64  `json:"endAt" binding:"required,gtfield=StartAt"`
}

type UpdatePromotionRequest struct {
	Title        string `json:"title" binding:"required,max=64"`
	Description  string `json:"description" binding:"required,max=2000"`
	Public       bool   `json:"public" binding:""`
	DiscountType int32  `json:"discountType" binding:"required"`
	DiscountRate int64  `json:"discountRate" binding:"min=0"`
	Code         string `json:"code" binding:"len=8"`
	StartAt      int64  `json:"startAt" binding:"required"`
	EndAt        int64  `json:"endAt" binding:"required,gtfield=StartAt"`
}
