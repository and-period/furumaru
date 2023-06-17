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
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/gin-gonic/gin"
)

func (h *handler) uploadRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.POST("/coordinators/thumbnail", h.uploadCoordinatorThumbnail)
	arg.POST("/coordinators/header", h.uploadCoordinatorHeader)
	arg.POST("/coordinators/promotion-video", h.uploadCoordinatorPromotionVideo)
	arg.POST("/coordinators/bonus-video", h.uploadCoordinatorBonusVideo)
	arg.POST("/producers/thumbnail", h.uploadProducerThumbnail)
	arg.POST("/producers/header", h.uploadProducerHeader)
	arg.POST("/products/image", h.uploadProductImage)
	arg.POST("/products/video", h.uploadProductVideo)
	arg.POST("/product-types/icon", h.uploadProductTypeIcon)
}

func (h *handler) uploadCoordinatorThumbnail(ctx *gin.Context) {
	const filename = "thumbnail"
	h.uploadFile(ctx, filename, h.media.GenerateCoordinatorThumbnail)
}

func (h *handler) uploadCoordinatorHeader(ctx *gin.Context) {
	const filename = "image"
	h.uploadFile(ctx, filename, h.media.GenerateCoordinatorHeader)
}

func (h *handler) uploadCoordinatorPromotionVideo(ctx *gin.Context) {
	const filename = "video"
	h.uploadFile(ctx, filename, h.media.GenerateCoordinatorPromotionVideo)
}

func (h *handler) uploadCoordinatorBonusVideo(ctx *gin.Context) {
	const filename = "video"
	h.uploadFile(ctx, filename, h.media.GenerateCoordinatorBonusVideo)
}

func (h *handler) uploadProducerThumbnail(ctx *gin.Context) {
	const filename = "thumbnail"
	h.uploadFile(ctx, filename, h.media.GenerateProducerThumbnail)
}

func (h *handler) uploadProducerHeader(ctx *gin.Context) {
	const filename = "image"
	h.uploadFile(ctx, filename, h.media.GenerateProducerHeader)
}

func (h *handler) uploadProductImage(ctx *gin.Context) {
	const filename = "image"
	h.uploadFile(ctx, filename, h.media.GenerateProductMediaImage)
}

func (h *handler) uploadProductVideo(ctx *gin.Context) {
	const filename = "video"
	h.uploadFile(ctx, filename, h.media.GenerateProductMediaVideo)
}

func (h *handler) uploadProductTypeIcon(ctx *gin.Context) {
	const filename = "icon"
	h.uploadFile(ctx, filename, h.media.GenerateProductTypeIcon)
}

func (h *handler) uploadFile(
	ctx *gin.Context,
	filename string,
	generate func(context.Context, *media.GenerateFileInput) (string, error),
) {
	file, header, err := h.parseFile(ctx, filename)
	if err != nil {
		httpError(ctx, err)
		return
	}
	in := &media.GenerateFileInput{
		File:   file,
		Header: header,
	}
	url, err := generate(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.UploadImageResponse{
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
