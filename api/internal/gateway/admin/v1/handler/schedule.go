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
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) scheduleRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("", h.ListSchedules)
	arg.POST("", h.CreateSchedule)
	arg.GET("/:scheduleId", h.filterAccessSchedule, h.GetSchedule)
	arg.PATCH("/:scheduleId", h.filterAccessSchedule, h.UpdateSchedule)
	arg.PATCH("/:scheduleId/approval", h.filterAccessSchedule, h.ApproveSchedule)
}

func (h *handler) filterAccessSchedule(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			schedule, err := h.getSchedule(ctx, util.GetParam(ctx, "scheduleId"))
			if err != nil {
				return false, err
			}
			return schedule.CoordinatorID == getAdminID(ctx), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.Next()
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

	res := &response.SchedulesResponse{
		Schedules:    service.NewSchedules(schedules).Response(),
		Coordinators: coordinators.Response(),
		Shippings:    shippings.Response(),
		Total:        total,
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

	var (
		coordinator *service.Coordinator
		shipping    *service.Shipping
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ectx, schedule.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		shipping, err = h.getShipping(ectx, schedule.ShippingID)
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	res := response.ScheduleResponse{
		Schedule:    schedule.Response(),
		Coordinator: coordinator.Response(),
		Shipping:    shipping.Response(),
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

	res := &response.ScheduleResponse{
		Schedule:    sschedule.Response(),
		Coordinator: coordinator.Response(),
		Shipping:    shipping.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateSchedule(ctx *gin.Context) {
	req := &request.UpdateScheduleRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	_, err := h.getShipping(ctx, req.ShippingID)
	if errors.Is(err, exception.ErrNotFound) {
		badRequest(ctx, err)
		return
	}
	if err != nil {
		httpError(ctx, err)
		return
	}

	var thumbnailURL, imageURL, openingVideoURL string
	eg, ectx := errgroup.WithContext(ctx)
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

	in := &store.UpdateScheduleInput{
		ScheduleID:      util.GetParam(ctx, "scheduleId"),
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
	if err := h.store.UpdateSchedule(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) ApproveSchedule(ctx *gin.Context) {
	req := &request.ApproveScheduleRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &store.ApproveScheduleInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
		AdminID:    getAdminID(ctx),
		Approved:   req.Approved,
	}
	if err := h.store.ApproveSchedule(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
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
