package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) coordinatorRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/coordinators")

	r.GET("", h.ListCoordinators)
	r.GET("/:coordinatorId", h.GetCoordinator)
}

func (h *handler) ListCoordinators(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	coordinatorsIn := &user.ListCoordinatorsInput{
		Limit:  limit,
		Offset: offset,
	}
	coordinators, total, err := h.user.ListCoordinators(ctx, coordinatorsIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(coordinators) == 0 {
		res := &response.CoordinatorsResponse{
			Coordinators: []*response.Coordinator{},
			ProductTypes: []*response.ProductType{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	shops, err := h.listShopsByCoordinatorIDs(ctx, coordinators.IDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	productTypes, err := h.multiGetProductTypes(ctx, shops.ProductTypeIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.CoordinatorsResponse{
		Coordinators: service.NewCoordinators(coordinators, shops.MapByCoordinatorID()).Response(),
		ProductTypes: productTypes.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetCoordinator(ctx *gin.Context) {
	coordinator, err := h.getCoordinator(ctx, util.GetParam(ctx, "coordinatorId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		shipping     *service.Shipping
		lives        service.LiveSummaries
		archives     service.ArchiveSummaries
		productTypes service.ProductTypes
		producers    service.Producers
		products     service.Products
		experiences  service.Experiences
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		shipping, err = h.getShippingByCoordinatorID(ectx, coordinator.ID)
		return
	})
	eg.Go(func() (err error) {
		params := &listLiveSummariesParams{
			shopID:  coordinator.ShopID,
			noLimit: true,
		}
		lives, _, err = h.listLiveSummaries(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		params := &listArchiveSummariesParams{
			shopID:  coordinator.ShopID,
			noLimit: true,
		}
		archives, _, err = h.listArchiveSummaries(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		productTypes, err = h.multiGetProductTypes(ectx, coordinator.ProductTypeIDs)
		return
	})
	eg.Go(func() (err error) {
		producers, err = h.listProducersByCoordinatorID(ectx, coordinator.ID)
		return
	})
	eg.Go(func() (err error) {
		in := &store.ListProductsInput{
			ShopID:           coordinator.ShopID,
			OnlyPublished:    true,
			ExcludeOutOfSale: true,
			NoLimit:          true,
		}
		products, err = h.listProducts(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &store.ListExperiencesInput{
			ShopID:          coordinator.ShopID,
			OnlyPublished:   true,
			ExcludeFinished: true,
			NoLimit:         true,
		}
		experiences, err = h.listExperiences(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.CoordinatorResponse{
		Coordinator:  coordinator.Response(),
		Shipping:     shipping.Response(),
		Lives:        lives.Response(),
		Archives:     archives.Response(),
		ProductTypes: productTypes.Response(),
		Producers:    producers.Response(),
		Products:     products.Response(),
		Experiences:  experiences.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) multiGetCoordinators(
	ctx context.Context,
	coordinatorIDs []string,
) (service.Coordinators, error) {
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

func (h *handler) multiGetCoordinatorsWithDeleted(
	ctx context.Context,
	coordinatorIDs []string,
) (service.Coordinators, error) {
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

func (h *handler) getCoordinator(
	ctx context.Context,
	coordinatorID string,
) (*service.Coordinator, error) {
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

func (h *handler) getCoordinatorWithDeleted(
	ctx context.Context,
	coordinatorID string,
) (*service.Coordinator, error) {
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
