package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Video
// @tag.description 動画関連
func (h *handler) videoRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/videos", h.authentication)

	r.GET("", h.ListVideos)
	r.POST("", h.CreateVideo)
	r.GET("/:videoId", h.filterAccessVideo, h.GetVideo)
	r.PATCH("/:videoId", h.filterAccessVideo, h.UpdateVideo)
	r.DELETE("/:videoId", h.filterAccessVideo, h.DeleteVideo)
	r.GET("/:videoId/analytics", h.filterAccessVideo, h.AnalyzeVideo)
}

func (h *handler) filterAccessVideo(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			video, err := h.getVideo(ctx, util.GetParam(ctx, "videoId"))
			if err != nil {
				return false, err
			}
			return video.CoordinatorID == getAdminID(ctx), nil
		},
		producer: func(ctx *gin.Context) (bool, error) {
			return false, nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Next()
}

// @Summary     動画一覧取得
// @Description 動画の一覧を取得します。ページネーションと名前でのフィルタリングに対応しています。
// @Tags        Video
// @Router      /v1/videos [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Param       name query string false "動画名" example("春の特産品紹介")
// @Produce     json
// @Success     200 {object} types.VideosResponse
func (h *handler) ListVideos(ctx *gin.Context) {
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

	in := &media.ListVideosInput{
		Name:   util.GetQuery(ctx, "name", ""),
		Limit:  limit,
		Offset: offset,
	}
	if getAdminType(ctx).Response() == types.AdminTypeCoordinator {
		in.CoordinatorID = getAdminID(ctx)
	}
	videos, total, err := h.media.ListVideos(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(videos) == 0 {
		res := &types.VideosResponse{
			Videos:       []*types.Video{},
			Coordinators: []*types.Coordinator{},
			Products:     []*types.Product{},
			Experiences:  []*types.Experience{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		coordinators service.Coordinators
		products     service.Products
		experiences  service.Experiences
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinators, err = h.multiGetCoordinators(ectx, videos.CoordinatorIDs())
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, videos.ProductIDs())
		return
	})
	eg.Go(func() (err error) {
		experiences, err = h.multiGetExperiences(ectx, videos.ExperienceIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.VideosResponse{
		Videos:       service.NewVideos(videos).Response(),
		Coordinators: coordinators.Response(),
		Products:     products.Response(),
		Experiences:  experiences.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     動画取得
// @Description 指定された動画の詳細情報を取得します。
// @Tags        Video
// @Router      /v1/videos/{videoId} [get]
// @Security    bearerauth
// @Param       videoId path string true "動画ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.VideoResponse
// @Failure     403 {object} util.ErrorResponse "動画の参照権限がない"
// @Failure     404 {object} util.ErrorResponse "動画が存在しない"
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

// @Summary     動画登録
// @Description 新しい動画を登録します。
// @Tags        Video
// @Router      /v1/videos [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.CreateVideoRequest true "動画情報"
// @Produce     json
// @Success     200 {object} types.VideoResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateVideo(ctx *gin.Context) {
	req := &types.CreateVideoRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	var (
		coordinator *service.Coordinator
		products    service.Products
		experiences service.Experiences
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ectx, req.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, req.ProductIDs)
		if len(products) != len(req.ProductIDs) {
			return fmt.Errorf("handler: unmatch products length: %w", exception.ErrInvalidArgument)
		}
		return
	})
	eg.Go(func() (err error) {
		experiences, err = h.multiGetExperiences(ectx, req.ExperienceIDs)
		if len(experiences) != len(req.ExperienceIDs) {
			return fmt.Errorf("handler: unmatch experiences length: %w", exception.ErrInvalidArgument)
		}
		return
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		h.badRequest(ctx, err)
		return
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	if getAdminType(ctx).Response() == types.AdminTypeCoordinator {
		if !currentAdmin(ctx, req.CoordinatorID) {
			h.httpError(ctx, exception.ErrForbidden)
			return
		}
	}

	in := &media.CreateVideoInput{
		Title:             req.Title,
		Description:       req.Description,
		CoordinatorID:     req.CoordinatorID,
		ProductIDs:        req.ProductIDs,
		ExperienceIDs:     req.ExperienceIDs,
		ThumbnailURL:      req.ThumbnailURL,
		VideoURL:          req.VideoURL,
		Public:            req.Public,
		Limited:           req.Limited,
		DisplayProduct:    req.DisplayProduct,
		DisplayExperience: req.DisplayExperience,
		PublishedAt:       jst.ParseFromUnix(req.PublishedAt),
	}
	video, err := h.media.CreateVideo(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.VideoResponse{
		Video:       service.NewVideo(video).Response(),
		Coordinator: coordinator.Response(),
		Products:    products.Response(),
		Experiences: experiences.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     動画更新
// @Description 動画の情報を更新します。
// @Tags        Video
// @Router      /v1/videos/{videoId} [patch]
// @Security    bearerauth
// @Param       videoId path string true "動画ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.UpdateVideoRequest true "動画情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "動画の更新権限がない"
// @Failure     404 {object} util.ErrorResponse "動画が存在しない"
func (h *handler) UpdateVideo(ctx *gin.Context) {
	req := &types.UpdateVideoRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		products, err := h.multiGetProducts(ectx, req.ProductIDs)
		if err != nil {
			return err
		}
		if len(products) != len(req.ProductIDs) {
			return fmt.Errorf("handler: unmatch products length: %w", exception.ErrInvalidArgument)
		}
		return nil
	})
	eg.Go(func() error {
		experiences, err := h.multiGetExperiences(ectx, req.ExperienceIDs)
		if err != nil {
			return err
		}
		if len(experiences) != len(req.ExperienceIDs) {
			return fmt.Errorf("handler: unmatch experiences length: %w", exception.ErrInvalidArgument)
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	in := &media.UpdateVideoInput{
		VideoID:           util.GetParam(ctx, "videoId"),
		Title:             req.Title,
		Description:       req.Description,
		ProductIDs:        req.ProductIDs,
		ExperienceIDs:     req.ExperienceIDs,
		ThumbnailURL:      req.ThumbnailURL,
		VideoURL:          req.VideoURL,
		Public:            req.Public,
		Limited:           req.Limited,
		DisplayProduct:    req.DisplayProduct,
		DisplayExperience: req.DisplayExperience,
		PublishedAt:       jst.ParseFromUnix(req.PublishedAt),
	}
	if err := h.media.UpdateVideo(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     動画分析情報取得
// @Description 指定された動画の視聴者分析データを取得します。集計期間と集計間隔を指定できます。
// @Tags        Video
// @Router      /v1/videos/{videoId}/analytics [get]
// @Security    bearerauth
// @Param       videoId path string true "動画ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       start query integer false "集計開始日時 (unixtime,未指定の場合は動画公開時間)" example("1640962800")
// @Param       end query integer false "集計終了日時 (unixtime,未指定の場合は現在時刻)" example("1640962800")
// @Param       viewerLogInterval query string false "集計間隔 (未指定の場合は1分間隔)" example("minute")
// @Produce     json
// @Success     200 {object} types.AnalyzeVideoResponse
// @Failure     403 {object} util.ErrorResponse "動画の参照権限がない"
// @Failure     404 {object} util.ErrorResponse "動画が存在しない"
func (h *handler) AnalyzeVideo(ctx *gin.Context) {
	const defaultViewerLogInterval = types.VideoViewerLogIntervalMinute

	video, err := h.getVideo(ctx, util.GetParam(ctx, "videoId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	startAtUnix, err := util.GetQueryInt64(ctx, "start", video.PublishedAt)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	endAtUnix, err := util.GetQueryInt64(ctx, "end", jst.Now().Unix())
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	viewerLogIntervalStr := util.GetQuery(ctx, "viewerLogInterval", string(defaultViewerLogInterval))

	startAt := jst.ParseFromUnix(startAtUnix)
	endAt := jst.ParseFromUnix(endAtUnix)
	viewerLogInterval := service.NewVideoViewerLogIntervalFromRequest(viewerLogIntervalStr)

	viewerLogIn := &media.AggregateVideoViewerLogsInput{
		VideoID:      video.ID,
		Interval:     viewerLogInterval.MediaEntity(),
		CreatedAtGte: startAt,
		CreatedAtLt:  endAt,
	}
	viewerLogs, totalViewers, err := h.media.AggregateVideoViewerLogs(ctx, viewerLogIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.AnalyzeVideoResponse{
		ViewerLogs:   service.NewVideoViewerLogs(viewerLogInterval, startAt, endAt, viewerLogs).Response(),
		TotalViewers: totalViewers,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     動画削除
// @Description 動画を削除します。
// @Tags        Video
// @Router      /v1/videos/{videoId} [delete]
// @Security    bearerauth
// @Param       videoId path string true "動画ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "動画の削除権限がない"
// @Failure     404 {object} util.ErrorResponse "動画が存在しない"
func (h *handler) DeleteVideo(ctx *gin.Context) {
	in := &media.DeleteVideoInput{
		VideoID: util.GetParam(ctx, "videoId"),
	}
	if err := h.media.DeleteVideo(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) getVideo(ctx context.Context, videoID string) (*service.Video, error) {
	in := &media.GetVideoInput{
		VideoID: videoID,
	}
	video, err := h.media.GetVideo(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewVideo(video), nil
}
