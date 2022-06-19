package store

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
