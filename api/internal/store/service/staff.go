package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *storeService) ListStaffsByStoreID(
	ctx context.Context, in *store.ListStaffsByStoreIDInput,
) (entity.Staffs, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	staffs, err := s.db.Staff.ListByStoreID(ctx, in.StoreID)
	return staffs, exception.InternalError(err)
}
