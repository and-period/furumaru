package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
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
	if len(coordinators) == 0 {
		return service.Coordinators{}, nil
	}
	shops, err := h.listShopsByCoordinatorIDs(ctx, coordinatorIDs)
	if err != nil {
		return nil, err
	}
	return service.NewCoordinators(coordinators, shops.MapByCoordinatorID()), nil
}

func (h *handler) multiGetCoordinatorsWithDeleted(ctx context.Context, coordinatorIDs []string) (service.Coordinators, error) {
	if len(coordinatorIDs) == 0 {
		return service.Coordinators{}, nil
	}
	in := &user.MultiGetCoordinatorsInput{
		CoordinatorIDs: coordinatorIDs,
		WithDeleted:    true,
	}
	coordinators, err := h.user.MultiGetCoordinators(ctx, in)
	if err != nil {
		return nil, err
	}
	if len(coordinators) == 0 {
		return service.Coordinators{}, nil
	}
	shops, err := h.listShopsByCoordinatorIDs(ctx, coordinatorIDs)
	if err != nil {
		return nil, err
	}
	return service.NewCoordinators(coordinators, shops.MapByCoordinatorID()), nil
}

func (h *handler) getCoordinator(ctx context.Context, coordinatorID string) (*service.Coordinator, error) {
	in := &user.GetCoordinatorInput{
		CoordinatorID: coordinatorID,
	}
	coordinator, err := h.user.GetCoordinator(ctx, in)
	if err != nil {
		return nil, err
	}
	shop, err := h.getShopByCoordinatorID(ctx, coordinatorID)
	if err != nil {
		return nil, err
	}
	return service.NewCoordinator(coordinator, shop), nil
}

func (h *handler) getCoordinatorWithDeleted(ctx context.Context, coordinatorID string) (*service.Coordinator, error) {
	in := &user.GetCoordinatorInput{
		CoordinatorID: coordinatorID,
		WithDeleted:   true,
	}
	coordinator, err := h.user.GetCoordinator(ctx, in)
	if err != nil {
		return nil, err
	}
	shop, err := h.getShopByCoordinatorID(ctx, coordinatorID)
	if err != nil {
		return nil, err
	}
	return service.NewCoordinator(coordinator, shop), nil
}
