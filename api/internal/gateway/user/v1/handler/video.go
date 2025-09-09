package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Video
// @tag.description 動画関連
func (h *handler) videoRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/videos")

	r.GET("", h.ListVideos)
	r.GET("/:videoId", h.createVideoViewerLog, h.GetVideo)
}

// @Summary     動画一覧取得
// @Description 動画の一覧を取得します。
// @Tags        Video
// @Router      /videos [get]
// @Param       limit query int64 false "取得上限数(max:200)" default(20)
// @Param       offset query int64 false "取得開始位置(min:0)" default(0)
// @Param       category query string false "カテゴリ（all, product, experience）" default("all")
// @Param       coordinatorId query string false "コーディネーターID"
// @Param       name query string false "動画名（部分一致検索）"
// @Produce     json
// @Success     200 {object} types.VideosResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) ListVideos(ctx *gin.Context) {
	const (
		defaultLimit    = 20
		defaultOffset   = 0
		defaultCategory = "all"
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
	category := util.GetQuery(ctx, "category", defaultCategory)

	params := &listVideoSummariesParams{
		name:          util.GetQuery(ctx, "name", ""),
		category:      videoCategory(category),
		coordinatorID: util.GetQuery(ctx, "coordinatorId", ""),
		limit:         limit,
		offset:        offset,
		noLimit:       false,
	}
	videos, total, err := h.listVideoSummaries(ctx, params)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(videos) == 0 {
		res := &types.VideosResponse{
			Videos:       []*types.VideoSummary{},
			Coordinators: []*types.Coordinator{},
		}
		ctx.JSON(http.StatusOK, res)
	}

	coordinators, err := h.multiGetCoordinators(ctx, videos.CoordinatorIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.VideosResponse{
		Videos:       videos.Response(),
		Coordinators: coordinators.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     動画詳細取得
// @Description 指定されたIDの動画詳細を取得します。
// @Tags        Video
// @Router      /videos/{videoId} [get]
// @Param       videoId path string true "動画ID"
// @Produce     json
// @Success     200 {object} types.VideoResponse
// @Failure     404 {object} util.ErrorResponse "動画が見つかりません"
func (h *handler) GetVideo(ctx *gin.Context) {
	video, err := h.getVideo(ctx, util.GetParam(ctx, "videoId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		coordinator *service.Coordinator
		products    service.Products
		experiences service.Experiences
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ectx, video.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, video.ProductIDs)
		return
	})
	eg.Go(func() (err error) {
		experiences, err = h.multiGetExperiences(ectx, video.ExperienceIDs)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.VideoResponse{
		Video:       video.Response(),
		Coordinator: coordinator.Response(),
		Products:    products.Response(),
		Experiences: experiences.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) getVideo(ctx context.Context, videoID string) (*service.Video, error) {
	in := &media.GetVideoInput{
		VideoID: videoID,
	}
	video, err := h.media.GetVideo(ctx, in)
	if err != nil {
		return nil, err
	}
	if !video.Published() {
		return nil, fmt.Errorf("handler: video is not published. video=%s: %w", videoID, exception.ErrNotFound)
	}
	return service.NewVideo(video), nil
}

type videoCategory string

const (
	videoCategoryAll        videoCategory = "all"
	videoCategoryProduct    videoCategory = "product"
	videoCategoryExperience videoCategory = "experience"
)

type listVideoSummariesParams struct {
	name          string
	category      videoCategory
	coordinatorID string
	limit         int64
	offset        int64
	noLimit       bool
}

func (h *handler) listVideoSummaries(ctx context.Context, params *listVideoSummariesParams) (service.VideoSummaries, int64, error) {
	in := &media.ListVideosInput{
		Name:           params.name,
		CoordinatorID:  params.coordinatorID,
		OnlyPublished:  true,
		ExcludeLimited: true,
		Limit:          params.limit,
		Offset:         params.offset,
		NoLimit:        params.noLimit,
	}
	switch params.category {
	case videoCategoryAll:
	case videoCategoryProduct:
		in.OnlyDisplayProduct = true
	case videoCategoryExperience:
		in.OnlyDisplayExperience = true
	default:
		return nil, 0, fmt.Errorf("handler: invalid category. category=%s: %w", params.category, exception.ErrInvalidArgument)
	}
	videos, total, err := h.media.ListVideos(ctx, in)
	if err != nil {
		return nil, 0, err
	}
	return service.NewVideoSummaries(videos), total, nil
}
