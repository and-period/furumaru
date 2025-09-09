package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/gin-gonic/gin"
)

// @tag.name        Upload
// @tag.description アップロード関連
func (h *handler) uploadRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/upload")

	r.GET("/state", h.GetUploadState)
	r.POST("/users/thumbnail", h.authentication, h.CreateUserThumbnailURL)
	r.POST("/spots/thumbnail", h.authentication, h.CreateSpotThumbnailURL)
}

// @Summary     アップロード状態取得
// @Description アップロードファイルの処理状態を取得します。
// @Tags        Upload
// @Router      /upload/state [get]
// @Param       key query string true "アップロードキー"
// @Produce     json
// @Success     200 {object} types.UploadStateResponse
// @Failure     404 {object} util.ErrorResponse "アップロードファイルが見つかりません"
func (h *handler) GetUploadState(ctx *gin.Context) {
	in := &media.GetUploadEventInput{
		Key: util.GetQuery(ctx, "key", ""),
	}
	event, err := h.media.GetUploadEvent(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.UploadStateResponse{
		URL:    event.ReferenceURL,
		Status: service.NewUploadStatus(event.Status).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     ユーザーサムネイルアップロードURL取得
// @Description ユーザーサムネイルをアップロードするためのURLを取得します。
// @Tags        Upload
// @Router      /upload/users/thumbnail [post]
// @Security    bearerauth
// @Accept      json
// @Produce     json
// @Param       body body types.GetUploadURLRequest true "ファイル情報"
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) CreateUserThumbnailURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetUserThumbnailUploadURL)
}

// @Summary     スポットサムネイルアップロードURL取得
// @Description スポットサムネイルをアップロードするためのURLを取得します。
// @Tags        Upload
// @Router      /upload/spots/thumbnail [post]
// @Security    bearerauth
// @Accept      json
// @Produce     json
// @Param       body body types.GetUploadURLRequest true "ファイル情報"
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) CreateSpotThumbnailURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetSpotThumbnailUploadURL)
}

func (h *handler) getUploadURL(ctx *gin.Context, fn func(context.Context, *media.GenerateUploadURLInput) (*entity.UploadEvent, error)) {
	req := &types.GetUploadURLRequest{}
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
	res := &types.UploadURLResponse{
		Key: event.Key,
		URL: event.UploadURL,
	}
	ctx.JSON(http.StatusOK, res)
}
