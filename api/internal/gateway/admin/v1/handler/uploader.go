package handler

import (
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/gin-gonic/gin"
)

func (h *handler) uploadRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/upload", h.authentication)

	r.GET("/state", h.GetUploadState)
	r.POST("/coordinators/thumbnail", h.CreateCoordinatorThumbnailUploadURL)
	r.POST("/coordinators/header", h.CreateCoordinatorHeaderUploadURL)
	r.POST("/coordinators/promotion-video", h.CreateCoordinatorPromotionVideoUploadURL)
	r.POST("/coordinators/bonus-video", h.CreateCoordinatorBonusVideoUploadURL)
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
}

func (h *handler) GetUploadState(ctx *gin.Context) {
	in := &media.GetUploadEventInput{
		UploadURL: util.GetQuery(ctx, "src", ""),
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

func (h *handler) getUploadURL(ctx *gin.Context, fn func(context.Context, *media.GenerateUploadURLInput) (string, error)) {
	req := &request.GetUploadURLRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.GenerateUploadURLInput{
		FileType: req.FileType,
	}
	url, err := fn(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.UploadURLResponse{
		URL: url,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) parseFile(ctx *gin.Context, filename string) (io.Reader, *multipart.FileHeader, error) {
	media, _, err := mime.ParseMediaType(ctx.GetHeader("Content-Type"))
	if err != nil {
		err := fmt.Errorf("handler: failed to parse media type. err=%s: %w", err.Error(), exception.ErrInvalidArgument)
		return nil, nil, err
	}
	if !strings.HasPrefix(media, "multipart/") {
		return nil, nil, fmt.Errorf("%s: %w", errInvalidFileFormat.Error(), exception.ErrInvalidArgument)
	}
	file, header, err := ctx.Request.FormFile(filename)
	if err != nil {
		err := fmt.Errorf("handler: failed to get file. err=%s: %w", err.Error(), exception.ErrInvalidArgument)
		return nil, nil, err
	}
	return file, header, nil
}
