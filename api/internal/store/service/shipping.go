package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) ListShippings(ctx context.Context, in *store.ListShippingsInput) (entity.Shippings, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	return nil, 0, exception.ErrNotImplemented
}

func (s *service) GetShipping(ctx context.Context, in *store.GetShippingInput) (*entity.Shipping, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	return nil, exception.ErrNotImplemented
}

func (s *service) CreateShipping(ctx context.Context, in *store.CreateShippingInput) (*entity.Shipping, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	return nil, exception.ErrNotImplemented
}

func (s *service) UpdateShipping(ctx context.Context, in *store.UpdateShippingInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	return exception.ErrNotImplemented
}

func (s *service) DeleteShipping(ctx context.Context, in *store.DeleteShippingInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	return exception.ErrNotImplemented
}
