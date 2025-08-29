package request

type CreatePromotionRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Public       bool   `json:"public"`
	DiscountType int32  `json:"discountType"`
	DiscountRate int64  `json:"discountRate"`
	Code         string `json:"code"`
	StartAt      int64  `json:"startAt"`
	EndAt        int64  `json:"endAt"`
}

type UpdatePromotionRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Public       bool   `json:"public"`
	DiscountType int32  `json:"discountType"`
	DiscountRate int64  `json:"discountRate"`
	Code         string `json:"code"`
	StartAt      int64  `json:"startAt"`
	EndAt        int64  `json:"endAt"`
}
