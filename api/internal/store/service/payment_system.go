package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) MultiGetPaymentSystems(
	ctx context.Context,
	in *store.MultiGetPaymentSystemsInput,
) (entity.PaymentSystems, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	systems, err := s.db.PaymentSystem.MultiGet(ctx, in.MethodTypes)
	return systems, internalError(err)
}

func (s *service) GetPaymentSystem(
	ctx context.Context,
	in *store.GetPaymentSystemInput,
) (*entity.PaymentSystem, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	system, err := s.db.PaymentSystem.Get(ctx, in.MethodType)
	return system, internalError(err)
}

func (s *service) UpdatePaymentSystem(
	ctx context.Context,
	in *store.UpdatePaymentStatusInput,
) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.PaymentSystem.Update(ctx, in.MethodType, in.Status)
	return internalError(err)
}
