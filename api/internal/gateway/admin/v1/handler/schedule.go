package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) scheduleRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("", h.ListSchedules)
	arg.POST("", h.CreateSchedule)
	arg.GET("/:scheduleId", h.GetSchedule)
}

func (h *handler) ListSchedules(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		badRequest(ctx, err)
		return
	}

	in := &store.ListSchedulesInput{
		Limit:  limit,
		Offset: offset,
	}
	schedules, total, err := h.store.ListSchedules(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	if len(schedules) == 0 {
		res := &response.SchedulesResponse{
			Schedules: []*response.Schedule{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		shippings    service.Shippings
		coordinators service.Coordinators
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		shippings, err = h.multiGetShippings(ectx, schedules.ShippingIDs())
		return
	})
	eg.Go(func() (err error) {
		coordinators, err = h.multiGetCoordinators(ectx, schedules.CoordinatorIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	sschedules := service.NewSchedules(schedules)
	sschedules.Fill(shippings.Map(), coordinators.Map())
	res := &response.SchedulesResponse{
		Schedules: sschedules.Response(),
		Total:     total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetSchedule(ctx *gin.Context) {
	scheduleID := util.GetParam(ctx, "scheduleId")
	schedule, err := h.getSchedule(ctx, scheduleID)
	if err != nil {
		httpError(ctx, err)
		return
	}
	lives, err := h.getScheduleDetailsByScheduleID(ctx, schedule)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := response.ScheduleResponse{
		Schedule: schedule.Response(),
		Lives:    lives.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateSchedule(ctx *gin.Context) {
	req := &request.CreateScheduleRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}
	if getRole(ctx).IsCoordinator() {
		if !currentAdmin(ctx, req.CoordinatorID) {
			forbidden(ctx, errors.New("handler: invalid coordinator id"))
			return
		}
	}

	var (
		coordinator *service.Coordinator
		shipping    *service.Shipping
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ectx, req.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		shipping, err = h.getShipping(ectx, req.ShippingID)
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

	var thumbnailURL, imageURL, openingVideoURL string
	eg, ectx = errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if req.ThumbnailURL == "" {
			return
		}
		in := &media.UploadFileInput{
			URL: req.ThumbnailURL,
		}
		thumbnailURL, err = h.media.UploadScheduleThumbnail(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		if req.ImageURL == "" {
			return
		}
		in := &media.UploadFileInput{
			URL: req.ImageURL,
		}
		imageURL, err = h.media.UploadScheduleImage(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		if req.OpeningVideoURL == "" {
			return
		}
		in := &media.UploadFileInput{
			URL: req.OpeningVideoURL,
		}
		openingVideoURL, err = h.media.UploadScheduleOpeningVideo(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	in := &store.CreateScheduleInput{
		CoordinatorID:   req.CoordinatorID,
		ShippingID:      req.ShippingID,
		Title:           req.Title,
		Description:     req.Description,
		ThumbnailURL:    thumbnailURL,
		ImageURL:        imageURL,
		OpeningVideoURL: openingVideoURL,
		Public:          req.Public,
		StartAt:         jst.ParseFromUnix(req.StartAt),
		EndAt:           jst.ParseFromUnix(req.EndAt),
	}
	schedule, err := h.store.CreateSchedule(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	sschedule := service.NewSchedule(schedule)
	sschedule.Fill(shipping, coordinator)

	res := &response.ScheduleResponse{
		Schedule: sschedule.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) getSchedule(ctx context.Context, scheduleID string) (*service.Schedule, error) {
	in := &store.GetScheduleInput{
		ScheduleID: scheduleID,
	}
	sschedule, err := h.store.GetSchedule(ctx, in)
	if err != nil {
		return nil, err
	}
	schedule := service.NewSchedule(sschedule)
	return schedule, nil
}

func (h *handler) getScheduleDetailsByScheduleID(ctx context.Context, schedule *service.Schedule) (service.Lives, error) {
	in := &store.ListLivesByScheduleIDInput{
		ScheduleID: schedule.ID,
	}
	slives, err := h.store.ListLivesByScheduleID(ctx, in)
	if err != nil {
		return nil, err
	}
	var (
		coordinator *service.Coordinator
		producers   service.Producers
		shipping    *service.Shipping
		products    service.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ectx, schedule.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		set := set.NewEmpty[string](len(slives))
		for i := range slives {
			set.Add(slives[i].ProducerID)
		}
		producerIDs := set.Slice()
		producers, err = h.multiGetProducers(ectx, producerIDs)
		if len(producers) != len(producerIDs) {
			return errInvalidProducerID
		}
		return
	})
	eg.Go(func() (err error) {
		shipping, err = h.getShipping(ectx, schedule.ShippingID)
		return
	})
	eg.Go(func() (err error) {
		set := set.NewEmpty[string](len(slives))
		for i := range slives {
			set.Add(slives[i].ProductIDs()...)
		}
		productIDs := set.Slice()
		products, err = h.multiGetProducts(ectx, productIDs)
		if len(products) != len(productIDs) {
			return errInvalidProductID
		}
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	lives := service.NewLives(slives)
	schedule.Fill(shipping, coordinator)
	lives.Fill(producers.Map(), products.Map())

	return lives, nil
}
