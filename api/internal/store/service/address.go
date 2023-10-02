package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) MultiGetAddresses(ctx context.Context, in *store.MultiGetAddressesInput) (entity.Addresses, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	addresses, err := s.db.Address.MultiGet(ctx, in.AddressIDs)
	return addresses, internalError(err)
}
