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
	r := rg.Group("/schedules", h.authentication)

	r.GET("", h.ListSchedules)
	r.POST("", h.CreateSchedule)
	r.GET("/:scheduleId", h.filterAccessSchedule, h.GetSchedule)
	r.PATCH("/:scheduleId", h.filterAccessSchedule, h.UpdateSchedule)
	r.PATCH("/:scheduleId/approval", h.filterAccessSchedule, h.ApproveSchedule)
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
		h.httpError(ctx, err)
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
		h.badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.ListSchedulesInput{
		Limit:  limit,
		Offset: offset,
	}
	if getRole(ctx) == service.AdminRoleCoordinator {
		in.CoordinatorID = getAdminID(ctx)
	}
	schedules, total, err := h.store.ListSchedules(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(schedules) == 0 {
		res := &response.SchedulesResponse{
			Schedules: []*response.Schedule{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	coordinators, err := h.multiGetCoordinators(ctx, schedules.CoordinatorIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.SchedulesResponse{
		Schedules:    service.NewSchedules(schedules).Response(),
		Coordinators: coordinators.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetSchedule(ctx *gin.Context) {
	scheduleID := util.GetParam(ctx, "scheduleId")
	schedule, err := h.getSchedule(ctx, scheduleID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	coordinator, err := h.getCoordinator(ctx, schedule.CoordinatorID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := response.ScheduleResponse{
		Schedule:    schedule.Response(),
		Coordinator: coordinator.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateSchedule(ctx *gin.Context) {
	req := &request.CreateScheduleRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	if getRole(ctx).IsCoordinator() {
		if !currentAdmin(ctx, req.CoordinatorID) {
			h.forbidden(ctx, errors.New("handler: invalid coordinator id"))
			return
		}
	}

	coordinator, err := h.getCoordinator(ctx, req.CoordinatorID)
	if errors.Is(err, exception.ErrNotFound) {
		h.badRequest(ctx, err)
		return
	}
	if err != nil {
		h.httpError(ctx, err)
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
		h.httpError(ctx, err)
		return
	}

	in := &store.CreateScheduleInput{
		CoordinatorID:   req.CoordinatorID,
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
		h.httpError(ctx, err)
		return
	}
	sschedule := service.NewSchedule(schedule)

	res := &response.ScheduleResponse{
		Schedule:    sschedule.Response(),
		Coordinator: coordinator.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateSchedule(ctx *gin.Context) {
	req := &request.UpdateScheduleRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
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
		h.httpError(ctx, err)
		return
	}

	in := &store.UpdateScheduleInput{
		ScheduleID:      util.GetParam(ctx, "scheduleId"),
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
		h.httpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) ApproveSchedule(ctx *gin.Context) {
	req := &request.ApproveScheduleRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.ApproveScheduleInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
		AdminID:    getAdminID(ctx),
		Approved:   req.Approved,
	}
	if err := h.store.ApproveSchedule(ctx, in); err != nil {
		h.httpError(ctx, err)
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
