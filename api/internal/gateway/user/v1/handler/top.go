package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Top
// @tag.description トップページ関連
func (h *handler) topRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/top")

	r.GET("/common", h.TopCommon)
}

// @Summary     トップページ共通データ取得
// @Description トップページに表示する共通データ（ライブ配信、アーカイブ、動画など）を取得します。
// @Tags        Top
// @Router      /top/common [get]
// @Produce     json
// @Success     200 {object} response.TopCommonResponse
func (h *handler) TopCommon(ctx *gin.Context) {
	const (
		defaultArchivesLimit = 6
		defaultVideosLimit   = 10
	)

	var (
		lives            service.LiveSummaries
		archives         service.ArchiveSummaries
		productVideos    service.VideoSummaries
		experienceVideos service.VideoSummaries
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		params := &listLiveSummariesParams{
			noLimit: true,
		}
		lives, _, err = h.listLiveSummaries(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		params := &listArchiveSummariesParams{
			limit:  defaultArchivesLimit,
			offset: 0,
		}
		archives, _, err = h.listArchiveSummaries(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		params := &listVideoSummariesParams{
			category: videoCategoryProduct,
			limit:    defaultVideosLimit,
			offset:   0,
		}
		productVideos, _, err = h.listVideoSummaries(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		params := &listVideoSummariesParams{
			category: videoCategoryExperience,
			limit:    defaultVideosLimit,
			offset:   0,
		}
		experienceVideos, _, err = h.listVideoSummaries(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	set := set.New(lives.CoordinatorIDs()...)
	set.Add(archives.CoordinatorIDs()...)
	set.Add(productVideos.CoordinatorIDs()...)
	set.Add(experienceVideos.CoordinatorIDs()...)
	coordinators, err := h.multiGetCoordinators(ctx, set.Slice())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.TopCommonResponse{
		Lives:            lives.Response(),
		Archives:         archives.Response(),
		ProductVideos:    productVideos.Response(),
		ExperienceVideos: experienceVideos.Response(),
		Coordinators:     coordinators.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
