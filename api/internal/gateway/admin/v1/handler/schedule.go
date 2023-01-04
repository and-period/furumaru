package handler

import (
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
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
	coordinatorID := req.CoordinatorID
	if getRole(ctx) == service.AdminRoleCoordinator {
		coordinatorID = getAdminID(ctx)
	}

	var (
		producers service.Producers
		shipping  *service.Shipping
		products  service.Products
		err       error
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		_, err = h.getCoordinator(ectx, coordinatorID)
		return
	})
	eg.Go(func() (err error) {
		set := set.NewEmpty[string](len(req.Lives))
		for i := range req.Lives {
			set.Add(req.Lives[i].ProducerID)
		}
		producerIDs := set.Slice()
		producers, err = h.multiGetProducers(ectx, producerIDs)
		if err != nil {
			return err
		}
		if len(producers) != len(producerIDs) {
			return errInvalidProducerID
		}
		return
	})
	eg.Go(func() (err error) {
		shipping, err = h.getShipping(ectx, req.ShippingID)
		return
	})
	eg.Go(func() error {
		set := set.NewEmpty[string](len(req.Lives))
		for i := range req.Lives {
			set.Add(req.Lives[i].ProductIDs...)
		}
		productIDs := set.Slice()
		products, err = h.multiGetProducts(ectx, productIDs)
		if err != nil {
			return err
		}
		if len(products) != len(productIDs) {
			return errInvalidProductID
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		switch {
		case errors.Is(err, exception.ErrNotFound),
			errors.Is(err, errInvalidProducerID),
			errors.Is(err, errInvalidProductID):
			badRequest(ctx, err)
		default:
			httpError(ctx, err)
		}
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
		CoordinatorID: coordinatorID,
		ShippingID:    req.ShippingID,
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

	schedule.Fill(shipping)
	lives.Fill(producers.Map(), products.Map())

	res := &response.ScheduleResponse{
		Schedule: schedule.Response(),
		Lives:    lives.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
