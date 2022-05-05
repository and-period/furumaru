package service

import (
	"context"

	"github.com/and-period/marche/api/internal/store/entity"
)

func (s *storeService) ListStaffsByStoreID(ctx context.Context, in *ListStaffsByStoreIDInput) (entity.Staffs, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, storeError(err)
	}
	staffs, err := s.db.Staff.ListByStoreID(ctx, in.StoreID)
	return staffs, storeError(err)
}
