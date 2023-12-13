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

	r.GET("/:coordinatorId", h.GetCoordinator)
}

func (h *handler) GetCoordinator(ctx *gin.Context) {
	coordinator, err := h.getCoordinator(ctx, util.GetParam(ctx, "coordinatorId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		lives    service.LiveSummaries
		archives service.ArchiveSummaries
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
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		productTypes service.ProductTypes
		producers    service.Producers
		products     service.Products
	)
	eg, ectx = errgroup.WithContext(ctx)
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
			CoordinatorID: coordinator.ID,
			EndAtGte:      h.now(),
			OnlyPublished: true,
			NoLimit:       true,
		}
		sproducts, _, err := h.store.ListProducts(ectx, in)
		if err != nil {
			return err
		}
		products = service.NewProducts(sproducts)
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
		Products:     products.Response(),
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
