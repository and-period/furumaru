package store

type ListStoresInput struct {
	Limit  int64 `validate:"required,max=200"`
	Offset int64 `validate:"min=0"`
}

type GetStoreInput struct {
	StoreID int64 `validate:"required"`
}

type CreateStoreInput struct {
	Name string `validate:"required,max=64"`
}

type UpdateStoreInput struct {
	StoreID      int64  `validate:"required"`
	Name         string `validate:"required,max=64"`
	ThumbnailURL string `validate:"omitempty,url"`
}

type ListStaffsByStoreIDInput struct {
	StoreID int64 `validate:"required"`
}
