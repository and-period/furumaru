package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

func (s *service) ListCoordinators(
	ctx context.Context, in *user.ListCoordinatorsInput,
) (entity.Coordinators, error) {
	// TODO: 詳細の実装
	return nil, exception.ErrNotImplemented
}

func (s *service) GetCoordinator(
	ctx context.Context, in *user.GetCoordinatorInput,
) (*entity.Coordinator, error) {
	// TODO: 詳細の実装
	return nil, exception.ErrNotImplemented
}

func (s *service) CreateCoordinator(
	ctx context.Context, in *user.CreateCoordinatorInput,
) (*entity.Coordinator, error) {
	// TODO: 詳細の実装
	return nil, exception.ErrNotImplemented
}
