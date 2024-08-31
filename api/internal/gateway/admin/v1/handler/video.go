package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/pkg/jst"
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
			Videos:       []*response.Video{},
			Coordinators: []*response.Coordinator{},
			Products:     []*response.Product{},
			Experiences:  []*response.Experience{},
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

	res := &response.VideoResponse{
		Video:       video.Response(),
		Coordinator: coordinator.Response(),
		Products:    products.Response(),
		Experiences: experiences.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateVideo(ctx *gin.Context) {
	req := &request.CreateVideoRequest{}
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

	if getRole(ctx) == service.AdminRoleCoordinator {
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

	res := &response.VideoResponse{
		Video:       service.NewVideo(video).Response(),
		Coordinator: coordinator.Response(),
		Products:    products.Response(),
		Experiences: experiences.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateVideo(ctx *gin.Context) {
	req := &request.UpdateVideoRequest{}
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
