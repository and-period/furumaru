package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/gin-gonic/gin"
)

func (h *handler) uploadRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/upload", h.authentication)

	r.GET("/state", h.GetUploadState)
	r.POST("/coordinators/thumbnail", h.CreateCoordinatorThumbnailUploadURL)
	r.POST("/coordinators/header", h.CreateCoordinatorHeaderUploadURL)
	r.POST("/coordinators/promotion-video", h.CreateCoordinatorPromotionVideoUploadURL)
	r.POST("/coordinators/bonus-video", h.CreateCoordinatorBonusVideoUploadURL)
	r.POST("/experiences/image", h.CreateExperienceImageUploadURL)
	r.POST("/experiences/video", h.CreateExperienceVideoUploadURL)
	r.POST("/experiences/promotion-video", h.CreateExperiencePromotionVideoUploadURL)
	r.POST("/producers/thumbnail", h.CreateProducerThumbnailUploadURL)
	r.POST("/producers/header", h.CreateProducerHeaderUploadURL)
	r.POST("/producers/promotion-video", h.CreateProducerPromotionVideoUploadURL)
	r.POST("/producers/bonus-video", h.CreateProducerBonusVideoUploadURL)
	r.POST("/products/image", h.CreateProductImageUploadURL)
	r.POST("/products/video", h.CreateProductVideoUploadURL)
	r.POST("/product-types/icon", h.CreateProductTypeIconUploadURL)
	r.POST("/schedules/thumbnail", h.CreateScheduleThumbnailUploadURL)
	r.POST("/schedules/image", h.CreateScheduleImageUploadURL)
	r.POST("/schedules/opening-video", h.CreateScheduleOpeningVideoUploadURL)
	r.POST("/schedules/:scheduleId/broadcasts/archive", h.CreateBroadcastArchiveMP4UploadURL)
	r.POST("/schedules/-/broadcasts/live", h.CreateBroadcastLiveMP4UploadURL)
	r.POST("/videos/thumbnail", h.CreateVideoThumbnailUploadURL)
	r.POST("/videos/file", h.CreateVideoFileUploadURL)
	r.POST("/spots/thumbnail", h.CreateSpotThumbnailURL)
}

func (h *handler) GetUploadState(ctx *gin.Context) {
	in := &media.GetUploadEventInput{
		Key: util.GetQuery(ctx, "key", ""),
	}
	event, err := h.media.GetUploadEvent(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.UploadStateResponse{
		URL:    event.ReferenceURL,
		Status: service.NewUploadStatus(event.Status).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateBroadcastArchiveMP4UploadURL(ctx *gin.Context) {
	req := &request.GetUploadURLRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.GenerateBroadcastArchiveMP4UploadInput{
		GenerateUploadURLInput: media.GenerateUploadURLInput{
			FileType: req.FileType,
		},
		ScheduleID: util.GetParam(ctx, "scheduleId"),
	}
	event, err := h.media.GetBroadcastArchiveMP4UploadURL(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.UploadURLResponse{
		Key: event.Key,
		URL: event.UploadURL,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateBroadcastLiveMP4UploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetBroadcastLiveMP4UploadURL)
}

func (h *handler) CreateCoordinatorThumbnailUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetCoordinatorThumbnailUploadURL)
}

func (h *handler) CreateCoordinatorHeaderUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetCoordinatorHeaderUploadURL)
}

func (h *handler) CreateCoordinatorPromotionVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetCoordinatorPromotionVideoUploadURL)
}

func (h *handler) CreateCoordinatorBonusVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetCoordinatorBonusVideoUploadURL)
}

func (h *handler) CreateExperienceImageUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetExperienceMediaImageUploadURL)
}

func (h *handler) CreateExperienceVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetExperienceMediaVideoUploadURL)
}

func (h *handler) CreateExperiencePromotionVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetExperiencePromotionVideoUploadURL)
}

func (h *handler) CreateProducerThumbnailUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetProducerThumbnailUploadURL)
}

func (h *handler) CreateProducerHeaderUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetProducerHeaderUploadURL)
}

func (h *handler) CreateProducerPromotionVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetProducerPromotionVideoUploadURL)
}

func (h *handler) CreateProducerBonusVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetProducerBonusVideoUploadURL)
}

func (h *handler) CreateProductImageUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetProductMediaImageUploadURL)
}

func (h *handler) CreateProductVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetProductMediaVideoUploadURL)
}

func (h *handler) CreateProductTypeIconUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetProductTypeIconUploadURL)
}

func (h *handler) CreateScheduleThumbnailUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetScheduleThumbnailUploadURL)
}

func (h *handler) CreateScheduleImageUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetScheduleImageUploadURL)
}

func (h *handler) CreateScheduleOpeningVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetScheduleOpeningVideoUploadURL)
}

func (h *handler) CreateVideoThumbnailUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetVideoThumbnailUploadURL)
}

func (h *handler) CreateVideoFileUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetVideoFileUploadURL)
}

func (h *handler) CreateSpotThumbnailURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetSpotThumbnailUploadURL)
}

func (h *handler) getUploadURL(ctx *gin.Context, fn func(context.Context, *media.GenerateUploadURLInput) (*entity.UploadEvent, error)) {
	req := &request.GetUploadURLRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.GenerateUploadURLInput{
		FileType: req.FileType,
	}
	event, err := fn(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.UploadURLResponse{
		Key: event.Key,
		URL: event.UploadURL,
	}
	ctx.JSON(http.StatusOK, res)
}
