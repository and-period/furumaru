package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
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

	productTypes, err := h.multiGetProductTypes(ctx, coordinators.ProductTypeIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.CoordinatorsResponse{
		Coordinators: service.NewCoordinators(coordinators).Response(),
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
		lives        service.LiveSummaries
		archives     service.ArchiveSummaries
		productTypes service.ProductTypes
		producers    service.Producers
		products     entity.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		params := &listLiveSummariesParams{
			coordinatorID: coordinator.ID,
		}
		lives, err = h.listLiveSummaries(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		params := &listArchiveSummariesParams{
			coordinatorID: coordinator.ID,
			noLimit:       true,
		}
		archives, err = h.listArchiveSummaries(ectx, params)
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
	eg.Go(func() error {
		in := &store.ListProductsInput{
			CoordinatorID:    coordinator.ID,
			OnlyPublished:    true,
			ExcludeOutOfSale: true,
			NoLimit:          true,
		}
		products, _, err = h.store.ListProducts(ectx, in)
		return nil
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.CoordinatorResponse{
		Coordinator:  coordinator.Response(),
		Lives:        lives.Response(),
		Archives:     archives.Response(),
		ProductTypes: productTypes.Response(),
		Producers:    producers.Response(),
		Products:     service.NewProducts(products).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

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
