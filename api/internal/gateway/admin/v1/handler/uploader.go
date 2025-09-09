package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/gin-gonic/gin"
)

// @tag.name        Upload
// @tag.description ファイルアップロード関連
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

// @Summary     アップロード状態取得
// @Description 指定されたキーのアップロード状態を取得します。
// @Tags        Upload
// @Router      /v1/upload/state [get]
// @Security    bearerauth
// @Param       key query string true "アップロードキー" example("upload-key-123")
// @Produce     json
// @Success     200 {object} types.UploadStateResponse
// @Failure     404 {object} util.ErrorResponse "アップロード情報が存在しない"
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

// @Summary     アーカイブ動画アップロードURL生成
// @Description 指定されたスケジュールのアーカイブ動画アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/schedules/{scheduleId}/broadcasts/archive [post]
// @Security    bearerauth
// @Param       scheduleId path string true "スケジュールID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateBroadcastArchiveMP4UploadURL(ctx *gin.Context) {
	req := &types.GetUploadURLRequest{}
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
	res := &types.UploadURLResponse{
		Key: event.Key,
		URL: event.UploadURL,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     ライブ動画アップロードURL生成
// @Description ライブ配信動画のアップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/schedules/-/broadcasts/live [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateBroadcastLiveMP4UploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetBroadcastLiveMP4UploadURL)
}

// @Summary     コーディネーターサムネイルアップロードURL生成
// @Description コーディネーターのサムネイル画像アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/coordinators/thumbnail [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateCoordinatorThumbnailUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetCoordinatorThumbnailUploadURL)
}

// @Summary     コーディネーターヘッダーアップロードURL生成
// @Description コーディネーターのヘッダー画像アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/coordinators/header [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateCoordinatorHeaderUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetCoordinatorHeaderUploadURL)
}

// @Summary     コーディネータープロモ動画アップロードURL生成
// @Description コーディネーターのプロモーション動画アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/coordinators/promotion-video [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateCoordinatorPromotionVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetCoordinatorPromotionVideoUploadURL)
}

// @Summary     コーディネーター特典動画アップロードURL生成
// @Description コーディネーターの特典動画アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/coordinators/bonus-video [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateCoordinatorBonusVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetCoordinatorBonusVideoUploadURL)
}

// @Summary     体験画像アップロードURL生成
// @Description 体験の画像アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/experiences/image [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateExperienceImageUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetExperienceMediaImageUploadURL)
}

// @Summary     体験動画アップロードURL生成
// @Description 体験の動画アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/experiences/video [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateExperienceVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetExperienceMediaVideoUploadURL)
}

// @Summary     体験プロモ動画アップロードURL生成
// @Description 体験のプロモーション動画アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/experiences/promotion-video [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateExperiencePromotionVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetExperiencePromotionVideoUploadURL)
}

// @Summary     生産者サムネイルアップロードURL生成
// @Description 生産者のサムネイル画像アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/producers/thumbnail [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateProducerThumbnailUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetProducerThumbnailUploadURL)
}

// @Summary     生産者ヘッダーアップロードURL生成
// @Description 生産者のヘッダー画像アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/producers/header [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateProducerHeaderUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetProducerHeaderUploadURL)
}

// @Summary     生産者プロモ動画アップロードURL生成
// @Description 生産者のプロモーション動画アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/producers/promotion-video [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateProducerPromotionVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetProducerPromotionVideoUploadURL)
}

// @Summary     生産者特典動画アップロードURL生成
// @Description 生産者の特典動画アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/producers/bonus-video [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateProducerBonusVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetProducerBonusVideoUploadURL)
}

// @Summary     商品画像アップロードURL生成
// @Description 商品の画像アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/products/image [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateProductImageUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetProductMediaImageUploadURL)
}

// @Summary     商品動画アップロードURL生成
// @Description 商品の動画アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/products/video [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateProductVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetProductMediaVideoUploadURL)
}

// @Summary     商品種別アイコンアップロードURL生成
// @Description 商品種別のアイコンアップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/product-types/icon [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateProductTypeIconUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetProductTypeIconUploadURL)
}

// @Summary     スケジュールサムネイルアップロードURL生成
// @Description スケジュールのサムネイルアップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/schedules/thumbnail [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateScheduleThumbnailUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetScheduleThumbnailUploadURL)
}

// @Summary     スケジュール画像アップロードURL生成
// @Description スケジュールの画像アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/schedules/image [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateScheduleImageUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetScheduleImageUploadURL)
}

// @Summary     スケジュールオープニング動画アップロードURL生成
// @Description スケジュールのオープニング動画アップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/schedules/opening-video [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateScheduleOpeningVideoUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetScheduleOpeningVideoUploadURL)
}

// @Summary     動画サムネイルアップロードURL生成
// @Description 動画のサムネイルアップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/videos/thumbnail [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateVideoThumbnailUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetVideoThumbnailUploadURL)
}

// @Summary     動画ファイルアップロードURL生成
// @Description 動画ファイルのアップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/videos/file [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateVideoFileUploadURL(ctx *gin.Context) {
	h.getUploadURL(ctx, h.media.GetVideoFileUploadURL)
}

// @Summary     スポットサムネイルアップロードURL生成
// @Description スポットのサムネイルアップロードURLを生成します。
// @Tags        Upload
// @Router      /v1/upload/spots/thumbnail [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetUploadURLRequest true "アップロードファイル情報"
// @Produce     json
// @Success     200 {object} types.UploadURLResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
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
