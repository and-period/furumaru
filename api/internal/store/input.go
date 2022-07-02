package store

import "github.com/and-period/furumaru/api/internal/store/entity"

type ListCategoriesInput struct {
	Name   string `validate:"omitempty,max=32"`
	Limit  int64  `validate:"required,max=200"`
	Offset int64  `validate:"min=0"`
}

type MultiGetCategoriesInput struct {
	CategoryIDs []string `validate:"omitempty,dive,required"`
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
	Name       string `validate:"omitempty,max=32"`
	CategoryID string `validate:"omitempty"`
	Limit      int64  `validate:"required,max=200"`
	Offset     int64  `validate:"min=0"`
}

type MultiGetProductTypesInput struct {
	ProductTypeIDs []string `validate:"omitempty,dive,required"`
}

type CreateProductTypeInput struct {
	Name       string `validate:"required,max=32"`
	CategoryID string `validate:"required"`
}

type UpdateProductTypeInput struct {
	ProductTypeID string `validate:"required"`
	Name          string `validate:"required,max=32"`
}

type DeleteProductTypeInput struct {
	ProductTypeID string `validate:"required"`
}

// TODO: ソート周りの対応
type ListProductsInput struct {
	Name          string `validate:"omitempty,max=128"`
	CoordinatorID string `validate:"omitempty"`
	ProducerID    string `validate:"omitempty"`
	Limit         int64  `validate:"required,max=200"`
	Offset        int64  `validate:"min=0"`
}

type GetProductInput struct {
	ProductID string `validate:"required"`
}

type CreateProductInput struct {
	CoordinatorID    string                `validate:"required"`
	ProducerID       string                `validate:"required"`
	CategoryID       string                `validate:"required"`
	TypeID           string                `validate:"required"`
	Name             string                `validate:"required,max=128"`
	Description      string                `validate:"required,max=2000"`
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
	CoordinatorID    string                `validate:"required"`
	ProducerID       string                `validate:"required"`
	CategoryID       string                `validate:"required"`
	TypeID           string                `validate:"required"`
	Name             string                `validate:"required,max=128"`
	Description      string                `validate:"required,max=2000"`
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

type UpdateProductMedia struct {
	URL         string `validate:"required,url"`
	IsThumbnail bool   `validate:""`
}

type DeleteProductInput struct {
	ProductID string `validate:"required"`
}
