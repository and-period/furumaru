package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/user"
)

func (h *handler) multiGetCoordinators(ctx context.Context, coordinatorIDs []string) (service.Coordinators, error) {
	if len(coordinatorIDs) == 0 {
		return service.Coordinators{}, nil
	}
	in := &user.MultiGetCoordinatorsInput{
		CoordinatorIDs: coordinatorIDs,
	}
	coordinators, err := h.user.MultiGetCoordinators(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewCoordinators(coordinators), nil
}

func (h *handler) getCoordinator(ctx context.Context, coordinatorID string) (*service.Coordinator, error) {
	in := &user.GetCoordinatorInput{
		CoordinatorID: coordinatorID,
	}
	coordinator, err := h.user.GetCoordinator(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewCoordinator(coordinator), nil
}
