package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) liveRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("/:liveId", h.GetLive)
	arg.POST("/:liveId/public", h.UpdateLivePublic)
}

func (h *handler) GetLive(ctx *gin.Context) {
	live, err := h.getLive(ctx, util.GetParam(ctx, "liveId"))
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.LiveResponse{
		Live: live.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateLivePublic(ctx *gin.Context) {
	req := &request.UpdateLivePublicRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &store.UpdateLivePublicInput{
		LiveID:      util.GetParam(ctx, "liveId"),
		Published:   true,
		Canceled:    false,
		ChannelName: req.ChannelName,
	}

	if err := h.store.UpdateLivePublic(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) getLive(ctx context.Context, liveID string) (*service.Live, error) {
	in := &store.GetLiveInput{
		LiveID: liveID,
	}
	slive, err := h.store.GetLive(ctx, in)
	if err != nil {
		return nil, err
	}
	live := service.NewLive(slive)
	if err := h.getLiveDetails(ctx, live); err != nil {
		return nil, err
	}
	return live, nil
}

func (h *handler) getLiveDetails(ctx context.Context, live *service.Live) error {
	var (
		products service.Products
		producer *service.Producer
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, live.ProductIDs)
		if err != nil {
			return err
		}
		if len(products) != len(live.ProductIDs) {
			return errors.New("error: invalid argument")
		}
		return nil
	})
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, live.ProducerID)
		return
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	live.Fill(producer, products.Map())
	return nil
}
