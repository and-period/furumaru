package request

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
