package handler

import (
	"context"
	"errors"
	"net/http"

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

func (h *handler) getLive(ctx context.Context, liveID string) (*service.Live, error) {
	in := &store.GetLiveInput{
		LiveID: liveID,
	}
	slive, err := h.store.GetLive(ctx, in)
	if err != nil {
		return nil, err
	}
	live := service.NewLive(slive)
	var (
		products service.Products
		producer *service.Producer
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		productIDs := make([]string, 0, len(live.Products))
		for i := range live.Products {
			productIDs = append(productIDs, live.Products[i].ID)
		}
		products, err = h.multiGetProducts(ectx, productIDs)
		if err != nil {
			return err
		}
		if len(products) != len(productIDs) {
			return errors.New("error: invalid argument")
		}
		return nil
	})
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, live.ProducerID)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	live.Fill(producer, products.Map())
	return live, nil
}
