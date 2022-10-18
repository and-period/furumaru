package handler

import (
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) scheduleRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.POST("", h.CreateSchedule)
}

func (h *handler) CreateSchedule(ctx *gin.Context) {
	req := &request.CreateScheduleRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	var (
		producers service.Producers
		products  service.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		in := &user.GetCoordinatorInput{
			CoordinatorID: req.CoordinatorID,
		}
		_, err := h.user.GetCoordinator(ectx, in)
		if err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		var productIDs []string
		for i := range req.Lives {
			productIDs = append(productIDs, req.Lives[i].ProductIDs...)
		}
		in := &store.MultiGetProductsInput{
			ProductIDs: productIDs,
		}
		sproducts, err := h.store.MultiGetProducts(ectx, in)
		if err != nil {
			return err
		}
		products = service.NewProducts(sproducts)
		return nil
	})
	eg.Go(func() error {
		var producerIDs []string
		for i := range req.Lives {
			producerIDs = append(producerIDs, req.Lives[i].ProducerID)
		}
		in := &user.MultiGetProducersInput{
			ProducerIDs: producerIDs,
		}
		sproducers, err := h.user.MultiGetProducers(ectx, in)
		if err != nil {
			return nil
		}
		producers = service.NewProducers(sproducers)
		return nil
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrInvalidArgument) {
		badRequest(ctx, err)
	}
	if err != nil {
		httpError(ctx, err)
		return
	}

	unixLives := make([]*store.CreateScheduleLive, len(req.Lives))
	for i := range req.Lives {
		unixLives[i] = &store.CreateScheduleLive{
			Title:       req.Lives[i].Title,
			Description: req.Lives[i].Description,
			ProducerID:  req.Lives[i].ProducerID,
			ProductIDs:  req.Lives[i].ProductIDs,
			StartAt:     jst.ParseFromUnix(req.Lives[i].StartAt),
			EndAt:       jst.ParseFromUnix(req.Lives[i].EndAt),
		}
	}

	in := &store.CreateScheduleInput{
		CoordinatorID: req.CoordinatorID,
		Title:         req.Title,
		Description:   req.Description,
		ThumbnailURL:  req.ThumbnailURL,
		StartAt:       jst.ParseFromUnix(req.StartAt),
		EndAt:         jst.ParseFromUnix(req.EndAt),
		Lives:         unixLives,
	}
	sschedule, slives, err := h.store.CreateSchedule(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	schedule := service.NewSchedule(sschedule)
	lives := service.NewLives(slives)

	lives.Fill(producers, products.Map())

	res := &response.ScheduleResponse{
		Schedule: schedule.Response(),
		Lives:    lives.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
