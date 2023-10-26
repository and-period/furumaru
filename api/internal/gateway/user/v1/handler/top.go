package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/media"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) topRoutes(rg *gin.RouterGroup) {
	rg.GET("/common", h.TopCommon)
}

func (h *handler) TopCommon(ctx *gin.Context) {
	const (
		defaultLivesLimit    = 100
		defaultArchivesLimit = 6
	)
	var (
		schedules entity.Schedules
		lives     entity.Lives
		products  entity.Products
		archives  entity.Schedules
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &store.ListSchedulesInput{
			EndAtGte:      h.now(),
			OnlyPublished: true,
			Limit:         defaultLivesLimit,
		}
		schedules, _, err = h.store.ListSchedules(ectx, in)
		if err != nil || len(schedules) == 0 {
			return
		}
		livesIn := &store.ListLivesInput{
			ScheduleIDs: schedules.IDs(),
		}
		lives, _, err = h.store.ListLives(ectx, livesIn)
		if err != nil || len(lives) == 0 {
			return
		}
		productsIn := &store.MultiGetProductsInput{
			ProductIDs: lives.ProductIDs(),
		}
		products, err = h.store.MultiGetProducts(ectx, productsIn)
		return
	})
	eg.Go(func() (err error) {
		broadcastsIn := &media.ListBroadcastsInput{
			OnlyArchived: true,
			Limit:        defaultArchivesLimit,
			Orders: []*media.ListBroadcastsOrder{{
				Key:        mentity.BroadcastOrderByUpdatedAt,
				OrderByASC: true,
			}},
		}
		broadcasts, _, err := h.media.ListBroadcasts(ectx, broadcastsIn)
		if err != nil || len(broadcasts) == 0 {
			return
		}
		schedulesIn := &store.MultiGetSchedulesInput{
			ScheduleIDs: broadcasts.ScheduleIDs(),
		}
		archives, err = h.store.MultiGetSchedules(ectx, schedulesIn)
		return
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.TopCommonResponse{
		Lives:    service.NewTopCommonLives(schedules, lives, products).Response(),
		Archives: service.NewTopCommonArchives(archives).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
