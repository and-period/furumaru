package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) liveRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication, h.filterAccessSchedule)
	arg.GET("", h.ListLives)
	arg.POST("", h.CreateLive)
	arg.GET("/:liveId", h.filterAccessLive, h.GetLive)
	arg.PATCH("/:liveId", h.filterAccessLive, h.UpdateLive)
	arg.DELETE("/:liveId", h.filterAccessLive, h.DeleteLive)
}

func (h *handler) filterAccessLive(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			live, err := h.getLive(ctx, util.GetParam(ctx, "liveId"))
			if err != nil {
				return false, err
			}
			producer, err := h.getProducer(ctx, live.ProducerID)
			if err != nil {
				return false, err
			}
			return currentAdmin(ctx, producer.CoordinatorID), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.Next()
}

func (h *handler) ListLives(ctx *gin.Context) {
	scheduleID := util.GetParam(ctx, "scheduleId")
	in := &store.ListLivesInput{
		ScheduleIDs: []string{scheduleID},
	}
	lives, total, err := h.store.ListLives(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	var (
		producers service.Producers
		products  service.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		producers, err = h.multiGetProducers(ectx, lives.ProducerIDs())
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, lives.ProductIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.LivesResponse{
		Lives:     service.NewLives(lives).Response(),
		Producers: producers.Response(),
		Products:  products.Response(),
		Total:     total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetLive(ctx *gin.Context) {
	liveID := util.GetParam(ctx, "liveId")
	live, err := h.getLive(ctx, liveID)
	if err != nil {
		httpError(ctx, err)
		return
	}
	var (
		producer *service.Producer
		products service.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, live.ProducerID)
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, live.ProductIDs)
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.LiveResponse{
		Live:     live.Response(),
		Producer: producer.Response(),
		Products: products.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateLive(ctx *gin.Context) {
	req := &request.CreateLiveRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	var (
		producer *service.Producer
		products service.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, req.ProducerID)
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, req.ProductIDs)
		if len(products) != len(req.ProductIDs) {
			return fmt.Errorf("handler: unmatch products length: %w", exception.ErrInvalidArgument)
		}
		return
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		badRequest(ctx, err)
		return
	}
	if err != nil {
		httpError(ctx, err)
		return
	}

	if getRole(ctx) == service.AdminRoleCoordinator {
		if !currentAdmin(ctx, producer.CoordinatorID) {
			forbidden(ctx, errors.New("handler: invalid coordinator id"))
			return
		}
	}

	in := &store.CreateLiveInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
		ProducerID: req.ProducerID,
		ProductIDs: req.ProductIDs,
		Comment:    req.Comment,
		StartAt:    jst.ParseFromUnix(req.StartAt),
		EndAt:      jst.ParseFromUnix(req.EndAt),
	}
	live, err := h.store.CreateLive(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.LiveResponse{
		Live:     service.NewLive(live).Response(),
		Producer: producer.Response(),
		Products: products.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateLive(ctx *gin.Context) {
	req := &request.UpdateLiveRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	products, err := h.multiGetProducts(ctx, req.ProductIDs)
	if err != nil {
		httpError(ctx, err)
		return
	}
	if len(products) != len(req.ProductIDs) {
		badRequest(ctx, errors.New("handler: unmatch products length"))
		return
	}

	in := &store.UpdateLiveInput{
		LiveID:     util.GetParam(ctx, "liveId"),
		ProductIDs: req.ProductIDs,
		Comment:    req.Comment,
		StartAt:    jst.ParseFromUnix(req.StartAt),
		EndAt:      jst.ParseFromUnix(req.EndAt),
	}
	if err := h.store.UpdateLive(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) DeleteLive(ctx *gin.Context) {
	in := &store.DeleteLiveInput{
		LiveID: util.GetParam(ctx, "liveId"),
	}
	if err := h.store.DeleteLive(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) getLive(ctx context.Context, liveID string) (*service.Live, error) {
	in := &store.GetLiveInput{
		LiveID: liveID,
	}
	live, err := h.store.GetLive(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewLive(live), nil
}
