package handler

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/and-period/marche/api/internal/gateway/admin/v1/response"
	"github.com/and-period/marche/api/pkg/uuid"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type uploadRegulation struct {
	dir      string
	filename string
	maxSize  int64
	formats  []string
}

func (h *apiV1Handler) uploadRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.POST("/stores/thumbnail", h.uploadStoreThumbnail)
}

func (h *apiV1Handler) uploadStoreThumbnail(ctx *gin.Context) {
	reg := &uploadRegulation{
		dir:      "stores/thumbnail",
		filename: "thumbnail",
		maxSize:  10 << 20, // 10MB
		formats:  []string{"image/png", "image/jpeg"},
	}
	h.upload(ctx, reg)
}

func (h *apiV1Handler) upload(ctx *gin.Context, reg *uploadRegulation) {
	file, err := h.parseFile(ctx, reg)
	if err != nil {
		httpError(ctx, err)
		return
	}
	path := h.generateFilePath(reg)
	url, err := h.storage.Upload(ctx, path, file)
	if err != nil {
		httpError(ctx, err)
		return
	}
	res := &response.UploaderResponse{
		URL: url,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) parseFile(ctx *gin.Context, reg *uploadRegulation) (multipart.File, error) {
	media, _, err := mime.ParseMediaType(ctx.GetHeader("Content-Type"))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	strings.HasPrefix(media, "multipart/")

	file, header, err := ctx.Request.FormFile(reg.filename)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if header.Size > reg.maxSize {
		return nil, status.Error(codes.InvalidArgument, errTooLargeFileSize.Error())
	}
	ok, err := h.validateFormat(file, reg)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if !ok {
		return nil, status.Error(codes.InvalidArgument, errInvalidFileFormat.Error())
	}
	return file, nil
}

func (h *apiV1Handler) validateFormat(file io.Reader, reg *uploadRegulation) (bool, error) {
	if len(reg.formats) == 0 {
		return true, nil
	}
	buf, err := io.ReadAll(file)
	if err != nil {
		return false, nil
	}
	contentType := http.DetectContentType(buf)
	for _, format := range reg.formats {
		if contentType == format {
			return true, nil
		}
	}
	return false, nil
}

func (h *apiV1Handler) generateFilePath(reg *uploadRegulation) string {
	key := uuid.Base58Encode(uuid.New())
	return strings.Join([]string{reg.dir, key}, "/")
}
