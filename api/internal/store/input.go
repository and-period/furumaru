package store

import (
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type ListCategoriesInput struct {
	Name   string                 `validate:"omitempty,max=32"`
	Limit  int64                  `validate:"required,max=200"`
	Offset int64                  `validate:"min=0"`
	Orders []*ListCategoriesOrder `validate:"omitempty,dive,required"`
}

type ListCategoriesOrder struct {
	Key        entity.CategoryOrderBy `validate:"required"`
	OrderByASC bool                   `validate:""`
}

type MultiGetCategoriesInput struct {
	CategoryIDs []string `validate:"omitempty,dive,required"`
}

type GetCategoryInput struct {
	CategoryID string `validate:"required"`
}

type CreateCategoryInput struct {
	Name string `validate:"required,max=32"`
}

type UpdateCategoryInput struct {
	CategoryID string `validate:"required"`
	Name       string `validate:"required,max=32"`
}

type DeleteCategoryInput struct {
	CategoryID string `validate:"required"`
}

type ListProductTypesInput struct {
	Name       string                   `validate:"omitempty,max=32"`
	CategoryID string                   `validate:"omitempty"`
	Limit      int64                    `validate:"required,max=200"`
	Offset     int64                    `validate:"min=0"`
	Orders     []*ListProductTypesOrder `validate:"omitempty,dive,required"`
}

type ListProductTypesOrder struct {
	Key        entity.ProductTypeOrderBy `validate:"required"`
	OrderByASC bool                      `validate:""`
}

type MultiGetProductTypesInput struct {
	ProductTypeIDs []string `validate:"omitempty,dive,required"`
}

type GetProductTypeInput struct {
	ProductTypeID string `validate:"required"`
}

type CreateProductTypeInput struct {
	Name       string `validate:"required,max=32"`
	IconURL    string `validate:"required"`
	CategoryID string `validate:"required"`
}

type UpdateProductTypeInput struct {
	ProductTypeID string `validate:"required"`
	Name          string `validate:"required,max=32"`
	IconURL       string `validate:"required"`
}

type UpdateProductTypeIconsInput struct {
	ProductTypeID string        `validate:"required"`
	Icons         common.Images `validate:""`
}

type DeleteProductTypeInput struct {
	ProductTypeID string `validate:"required"`
}

type ListShippingsInput struct {
	Limit  int64                 `validate:"required,max=200"`
	Offset int64                 `validate:"min=0"`
	Orders []*ListShippingsOrder `validate:"omitempty,dive,required"`
}

type ListShippingsOrder struct {
	Key        entity.ShippingOrderBy `validate:"required"`
	OrderByASC bool                   `validate:""`
}

type MultiGetShippingsInput struct {
	ShippingIDs []string `validate:"omitempty,dive,required"`
}

type GetShippingInput struct {
	ShippingID string `validate:"required"`
}

type CreateShippingInput struct {
	Name               string                `validate:"required,max=64"`
	Box60Rates         []*CreateShippingRate `validate:"required,dive,required"`
	Box60Refrigerated  int64                 `validate:"min=0,lt=10000000000"`
	Box60Frozen        int64                 `validate:"min=0,lt=10000000000"`
	Box80Rates         []*CreateShippingRate `validate:"required,dive,required"`
	Box80Refrigerated  int64                 `validate:"min=0,lt=10000000000"`
	Box80Frozen        int64                 `validate:"min=0,lt=10000000000"`
	Box100Rates        []*CreateShippingRate `validate:"required,dive,required"`
	Box100Refrigerated int64                 `validate:"min=0,lt=10000000000"`
	Box100Frozen       int64                 `validate:"min=0,lt=10000000000"`
	HasFreeShipping    bool                  `validate:""`
	FreeShippingRates  int64                 `validate:"min=0,lt=10000000000"`
}

type CreateShippingRate struct {
	Name        string  `validate:"required"`
	Price       int64   `validate:"required,lt=10000000000"`
	Prefectures []int64 `validate:"required,dive,min=0"`
}

type UpdateShippingInput struct {
	ShippingID         string                `validate:"required"`
	Name               string                `validate:"required,max=64"`
	Box60Rates         []*UpdateShippingRate `validate:"required,dive,required"`
	Box60Refrigerated  int64                 `validate:"min=0,lt=10000000000"`
	Box60Frozen        int64                 `validate:"min=0,lt=10000000000"`
	Box80Rates         []*UpdateShippingRate `validate:"required,dive,required"`
	Box80Refrigerated  int64                 `validate:"min=0,lt=10000000000"`
	Box80Frozen        int64                 `validate:"min=0,lt=10000000000"`
	Box100Rates        []*UpdateShippingRate `validate:"required,dive,required"`
	Box100Refrigerated int64                 `validate:"min=0,lt=10000000000"`
	Box100Frozen       int64                 `validate:"min=0,lt=10000000000"`
	HasFreeShipping    bool                  `validate:""`
	FreeShippingRates  int64                 `validate:"min=0,lt=10000000000"`
}

type UpdateShippingRate struct {
	Name        string  `validate:"required"`
	Price       int64   `validate:"required,lt=10000000000"`
	Prefectures []int64 `validate:"required,dive,min=0"`
}

type DeleteShippingInput struct {
	ShippingID string `validate:"required"`
}

type ListProductsInput struct {
	Name        string               `validate:"omitempty,max=128"`
	ProducerID  string               `validate:"omitempty"`
	ProducerIDs []string             `validate:"dive,required"`
	Limit       int64                `validate:"required,max=200"`
	Offset      int64                `validate:"min=0"`
	Orders      []*ListProductsOrder `validate:"omitempty,dive,required"`
}

type ListProductsOrder struct {
	Key        entity.ProductOrderBy `validate:"required"`
	OrderByASC bool                  `validate:""`
}

type MultiGetProductsInput struct {
	ProductIDs []string `validate:"omitempty,dive,required"`
}

type GetProductInput struct {
	ProductID string `validate:"required"`
}

type CreateProductInput struct {
	ProducerID       string                `validate:"required"`
	TypeID           string                `validate:"required"`
	Name             string                `validate:"required,max=128"`
	Description      string                `validate:"required,max=20000"`
	Public           bool                  `validate:""`
	Inventory        int64                 `validate:"min=0"`
	Weight           int64                 `validate:"min=0"`
	WeightUnit       entity.WeightUnit     `validate:"required,oneof=1 2"`
	Item             int64                 `validate:"min=1"`
	ItemUnit         string                `validate:"required,max=16"`
	ItemDescription  string                `validate:"required,max=64"`
	Media            []*CreateProductMedia `validate:"max=8,unique=URL"`
	Price            int64                 `validate:"min=0"`
	DeliveryType     entity.DeliveryType   `validate:"required,oneof=1 2 3"`
	Box60Rate        int64                 `validate:"min=0,max=100"`
	Box80Rate        int64                 `validate:"min=0,max=100"`
	Box100Rate       int64                 `validate:"min=0,max=100"`
	OriginPrefecture string                `validate:"omitempty,max=32"`
	OriginCity       string                `validate:"omitempty,max=32"`
}

type CreateProductMedia struct {
	URL         string `validate:"required,url"`
	IsThumbnail bool   `validate:""`
}

type UpdateProductInput struct {
	ProductID        string                `validate:"required"`
	ProducerID       string                `validate:"required"`
	TypeID           string                `validate:"required"`
	Name             string                `validate:"required,max=128"`
	Description      string                `validate:"required,max=20000"`
	Public           bool                  `validate:""`
	Inventory        int64                 `validate:"min=0"`
	Weight           int64                 `validate:"min=0"`
	WeightUnit       entity.WeightUnit     `validate:"required,oneof=1 2"`
	Item             int64                 `validate:"min=1"`
	ItemUnit         string                `validate:"required,max=16"`
	ItemDescription  string                `validate:"required,max=64"`
	Media            []*UpdateProductMedia `validate:"max=8,unique=URL"`
	Price            int64                 `validate:"min=0"`
	DeliveryType     entity.DeliveryType   `validate:"required,oneof=1 2 3"`
	Box60Rate        int64                 `validate:"min=0,max=100"`
	Box80Rate        int64                 `validate:"min=0,max=100"`
	Box100Rate       int64                 `validate:"min=0,max=100"`
	OriginPrefecture string                `validate:"omitempty,max=32"`
	OriginCity       string                `validate:"omitempty,max=32"`
}

type UpdateProductMedia struct {
	URL         string `validate:"required,url"`
	IsThumbnail bool   `validate:""`
}

type UpdateProductMediaInput struct {
	ProductID string                     `validate:"required"`
	Images    []*UpdateProductMediaImage `validate:"omitempty,dive,required,unique=URL"`
}

type UpdateProductMediaImage struct {
	OriginURL string        `validate:"required"`
	Images    common.Images `validate:""`
}

type DeleteProductInput struct {
	ProductID string `validate:"required"`
}

type ListPromotionsInput struct {
	Limit  int64                  `validate:"required,max=200"`
	Offset int64                  `validate:"min=0"`
	Orders []*ListPromotionsOrder `validate:"omitempty,dive,required"`
}

type ListPromotionsOrder struct {
	Key        entity.PromotionOrderBy `validate:"required"`
	OrderByASC bool                    `validate:""`
}

type GetPromotionInput struct {
	PromotionID string `validate:"required"`
}

type CreatePromotionInput struct {
	Title        string                   `validate:"required,max=64"`
	Description  string                   `validate:"required,max=2000"`
	Public       bool                     `validate:""`
	PublishedAt  time.Time                `validate:"required"`
	DiscountType entity.DiscountType      `validate:"required,oneof=1 2 3"`
	DiscountRate int64                    `validate:"min=0"`
	Code         string                   `validate:"len=8"`
	CodeType     entity.PromotionCodeType `validate:"required,oneof=1 2"`
	StartAt      time.Time                `validate:"required"`
	EndAt        time.Time                `validate:"required"`
}

type UpdatePromotionInput struct {
	PromotionID  string                   `validate:"required"`
	Title        string                   `validate:"required,max=64"`
	Description  string                   `validate:"required,max=2000"`
	Public       bool                     `validate:""`
	PublishedAt  time.Time                `validate:"required"`
	DiscountType entity.DiscountType      `validate:"required,oneof=1 2 3"`
	DiscountRate int64                    `validate:"min=0"`
	Code         string                   `validate:"len=8"`
	CodeType     entity.PromotionCodeType `validate:"required,oneof=1 2"`
	StartAt      time.Time                `validate:"required"`
	EndAt        time.Time                `validate:"required"`
}

type DeletePromotionInput struct {
	PromotionID string `validate:"required"`
}

type CreateScheduleInput struct {
	CoordinatorID string                `validate:"required"`
	ShippingID    string                `validate:"required"`
	Title         string                `validate:"required,max=64"`
	Description   string                `validate:"required,max=2000"`
	ThumbnailURL  string                `validate:"required"`
	StartAt       time.Time             `validate:"required"`
	EndAt         time.Time             `validate:"required"`
	Lives         []*CreateScheduleLive `validate:"required"`
}

type CreateScheduleLive struct {
	Title       string    `validate:"required,max=64"`
	Description string    `validate:"required,max=2000"`
	ProducerID  string    `validate:"required"`
	ProductIDs  []string  `validate:"dive,required,unique"`
	StartAt     time.Time `validate:"required"`
	EndAt       time.Time `validate:"required"`
}

type ListOrdersInput struct {
	CoordinatorID string             `validate:"omitempty"`
	Limit         int64              `validate:"required,max=200"`
	Offset        int64              `validate:"min=0"`
	Orders        []*ListOrdersOrder `validate:"omitempty,dive,required"`
}

type ListOrdersOrder struct {
	Key        entity.OrderOrderBy `validate:"required"`
	OrderByASC bool                `validate:""`
}

type GetOrderInput struct {
	OrderID string `validate:"required"`
}

type AggregateOrdersInput struct {
	UserIDs []string `validate:"omitempty,dive,required"`
}

type MultiGetAddressesInput struct {
	AddressIDs []string `validate:"omitempty,dive,required"`
}

type SearchPostalCodeInput struct {
	PostlCode string `validate:"required,numeric,len=7"`
}
