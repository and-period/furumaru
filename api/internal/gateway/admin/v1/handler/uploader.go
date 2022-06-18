package handler

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/pkg/uuid"
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
	arg.POST("/coordinators/thumbnail", h.uploadCoordinatorThumbnail)
	arg.POST("/coordinators/header", h.uploadCoordinatorHeader)
	arg.POST("/producers/thumbnail", h.uploadProducerThumbnail)
	arg.POST("/producers/header", h.uploadProducerHeader)
}

func (h *apiV1Handler) uploadCoordinatorThumbnail(ctx *gin.Context) {
	reg := &uploadRegulation{
		dir:      "coordinators/thumbnail",
		filename: "thumbnail",
		maxSize:  10 << 20, // 10MB
		formats:  []string{"image/png", "image/jpeg"},
	}
	h.upload(ctx, reg)
}

func (h *apiV1Handler) uploadCoordinatorHeader(ctx *gin.Context) {
	reg := &uploadRegulation{
		dir:      "coordinators/header",
		filename: "image",
		maxSize:  10 << 20, // 10MB
		formats:  []string{"image/png", "image/jpeg"},
	}
	h.upload(ctx, reg)
}

func (h *apiV1Handler) uploadProducerThumbnail(ctx *gin.Context) {
	reg := &uploadRegulation{
		dir:      "producers/thumbnail",
		filename: "thumbnail",
		maxSize:  10 << 20, // 10MB
		formats:  []string{"image/png", "image/jpeg"},
	}
	h.upload(ctx, reg)
}

func (h *apiV1Handler) uploadProducerHeader(ctx *gin.Context) {
	reg := &uploadRegulation{
		dir:      "producers/header",
		filename: "image",
		maxSize:  10 << 20, // 10MB
		formats:  []string{"image/png", "image/jpeg"},
	}
	h.upload(ctx, reg)
}

func (h *apiV1Handler) upload(ctx *gin.Context, reg *uploadRegulation) {
	file, header, err := h.parseFile(ctx, reg)
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

func (h *apiV1Handler) parseFile(ctx *gin.Context, reg *uploadRegulation) (io.Reader, *multipart.FileHeader, error) {
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

func (h *apiV1Handler) validateFormat(reg *uploadRegulation, file io.Reader) (bool, error) {
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

func (h *apiV1Handler) generateFilePath(reg *uploadRegulation, header *multipart.FileHeader) string {
	key := uuid.Base58Encode(uuid.New())
	extension := filepath.Ext(header.Filename)
	filename := strings.Join([]string{key, extension}, "")
	return strings.Join([]string{reg.dir, filename}, "/")
}
