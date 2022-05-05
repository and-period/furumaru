package service

import validator "github.com/go-playground/validator/v10"

//nolint:errcheck
func newValidator() *validator.Validate {
	return validator.New()
}

type ListStoresInput struct {
	Limit  int64 `validate:"required,max=200"`
	Offset int64 `validate:"min=0"`
}

type GetStoreInput struct {
	StoreID int64 `validate:"required"`
}

type ListStaffsByStoreIDInput struct {
	StoreID int64 `validate:"required"`
}
