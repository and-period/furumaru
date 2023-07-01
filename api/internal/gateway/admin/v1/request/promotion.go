package request

type CreatePromotionRequest struct {
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	Public       bool   `json:"public,omitempty"`
	DiscountType int32  `json:"discountType,omitempty"`
	DiscountRate int64  `json:"discountRate,omitempty"`
	Code         string `json:"code,omitempty"`
	StartAt      int64  `json:"startAt,omitempty"`
	EndAt        int64  `json:"endAt,omitempty"`
}

type UpdatePromotionRequest struct {
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	Public       bool   `json:"public,omitempty"`
	DiscountType int32  `json:"discountType,omitempty"`
	DiscountRate int64  `json:"discountRate,omitempty"`
	Code         string `json:"code,omitempty"`
	StartAt      int64  `json:"startAt,omitempty"`
	EndAt        int64  `json:"endAt,omitempty"`
}
