package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) videoRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/videos", h.authentication)

	r.GET("", h.ListVideos)
	r.POST("", h.CreateVideo)
	r.GET("/:videoId", h.filterAccessVideo, h.GetVideo)
	r.PATCH("/:videoId", h.filterAccessVideo, h.UpdateVideo)
	r.DELETE("/:videoId", h.filterAccessVideo, h.DeleteVideo)
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
		Name:          util.GetQuery(ctx, "name", ""),
		CoordinatorID: util.GetQuery(ctx, "coordinatorId", ""),
		Limit:         limit,
		Offset:        offset,
		NoLimit:       false,
	}
	videos, total, err := h.media.ListVideos(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(videos) == 0 {
		res := &response.VideosResponse{
			Videos: []*response.Video{},
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

	res := &response.VideosResponse{
		Videos:       service.NewVideos(videos).Response(),
		Coordinators: coordinators.Response(),
		Products:     products.Response(),
		Experiences:  experiences.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetVideo(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.VideoResponse{
		Video:       &response.Video{},
		Coordinator: &response.Coordinator{},
		Products:    []*response.Product{},
		Experiences: []*response.Experience{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateVideo(ctx *gin.Context) {
	req := &request.CreateVideoRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	res := &response.VideoResponse{
		Video:       &response.Video{},
		Coordinator: &response.Coordinator{},
		Products:    []*response.Product{},
		Experiences: []*response.Experience{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateVideo(ctx *gin.Context) {
	req := &request.UpdateVideoRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteVideo(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.Status(http.StatusNoContent)
}

func (h *handler) getVideo(ctx context.Context, videoID string) (*service.Video, error) {
	// TODO: 詳細の実装
	return &service.Video{}, nil
}
