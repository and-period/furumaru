package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) uploadRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.POST("/coordinators/thumbnail", h.uploadCoordinatorThumbnail)
	arg.POST("/coordinators/header", h.uploadCoordinatorHeader)
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

func (h *handler) uploadProducerThumbnail(ctx *gin.Context) {
	const filename = "thumbnail"
	h.uploadFile(ctx, filename, h.media.GenerateProducerThumbnail)
}

func (h *handler) uploadProducerHeader(ctx *gin.Context) {
	const filename = "image"
	h.uploadFile(ctx, filename, h.media.GenerateProducerHeader)
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

type uploadRegulation struct {
	dir      string
	filename string
	maxSize  int64
	formats  []string
}

func (h *handler) uploadProductImage(ctx *gin.Context) {
	reg := &uploadRegulation{
		dir:      "products/image",
		filename: "image",
		maxSize:  10 << 20, // 10MB
		formats:  []string{"image/png", "image/jpeg"},
	}
	h.uploadFileData(ctx, reg)
}

func (h *handler) uploadProductVideo(ctx *gin.Context) {
	reg := &uploadRegulation{
		dir:      "products/video",
		filename: "video",
		maxSize:  200 << 20, // 200MB
		formats:  []string{"video/mp4"},
	}
	h.uploadFileData(ctx, reg)
}

func (h *handler) uploadProductTypeIcon(ctx *gin.Context) {
	reg := &uploadRegulation{
		dir:      "product-types/icon",
		filename: "icon",
		maxSize:  10 << 20, // 10MB
		formats:  []string{"image/png", "image/jpeg"},
	}
	h.uploadFileData(ctx, reg)
}

// Deprecated: Use to upload
func (h *handler) uploadFileData(ctx *gin.Context, reg *uploadRegulation) {
	file, header, err := h.parseFileData(ctx, reg)
	if err != nil {
		httpError(ctx, err)
		return
	}
	path := h.generateFilePath(reg, header)
	url, err := h.storage.Upload(ctx, path, file)
	if err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.UploadImageResponse{
		URL: url,
	}
	ctx.JSON(http.StatusOK, res)
}

// Deprecated: parseFile
func (h *handler) parseFileData(ctx *gin.Context, reg *uploadRegulation) (io.Reader, *multipart.FileHeader, error) {
	media, _, err := mime.ParseMediaType(ctx.GetHeader("Content-Type"))
	if err != nil {
		return nil, nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if !strings.HasPrefix(media, "multipart/") {
		return nil, nil, status.Error(codes.InvalidArgument, errInvalidFileFormat.Error())
	}

	file, header, err := ctx.Request.FormFile(reg.filename)
	if err != nil {
		return nil, nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if header.Size > reg.maxSize {
		return nil, nil, status.Error(codes.InvalidArgument, errTooLargeFileSize.Error())
	}

	var buf bytes.Buffer
	teeReader := io.TeeReader(file, &buf)
	ok, err := h.validateFormat(reg, teeReader)
	if err != nil {
		return nil, nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if !ok {
		return nil, nil, status.Error(codes.InvalidArgument, errInvalidFileFormat.Error())
	}
	return &buf, header, nil
}

func (h *handler) validateFormat(reg *uploadRegulation, file io.Reader) (bool, error) {
	if len(reg.formats) == 0 {
		return true, nil
	}
	buf, err := io.ReadAll(file)
	if err != nil {
		return false, err
	}
	contentType := http.DetectContentType(buf)
	for _, format := range reg.formats {
		if contentType == format {
			return true, nil
		}
	}
	return false, nil
}

func (h *handler) generateFilePath(reg *uploadRegulation, header *multipart.FileHeader) string {
	key := uuid.Base58Encode(uuid.New())
	extension := filepath.Ext(header.Filename)
	filename := strings.Join([]string{key, extension}, "")
	return strings.Join([]string{reg.dir, filename}, "/")
}
