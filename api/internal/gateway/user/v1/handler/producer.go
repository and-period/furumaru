package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) producerRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/producers")

	r.GET("", h.ListProducers)
	r.GET("/:producerId", h.GetProducer)
}

func (h *handler) ListProducers(ctx *gin.Context) {
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

	producersIn := &user.ListProducersInput{
		Limit:  limit,
		Offset: offset,
	}
	producers, total, err := h.user.ListProducers(ctx, producersIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProducersResponse{
		Producers: service.NewProducers(producers).Response(),
		Total:     total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetProducer(ctx *gin.Context) {
	producer, err := h.getProducer(ctx, util.GetParam(ctx, "producerId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		lives    service.LiveSummaries
		archives service.ArchiveSummaries
		products entity.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		params := &listLiveSummariesParams{
			producerID: producer.ID,
		}
		lives, err = h.listLiveSummaries(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		params := &listArchiveSummariesParams{
			coordinatorID: producer.ID,
			noLimit:       true,
		}
		archives, err = h.listArchiveSummaries(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		in := &store.ListProductsInput{
			ProducerID:    producer.ID,
			EndAtGte:      h.now(),
			OnlyPublished: true,
			NoLimit:       true,
		}
		products, _, err = h.store.ListProducts(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProducerResponse{
		Producer: producer.Response(),
		Lives:    lives.Response(),
		Archives: archives.Response(),
		Products: service.NewProducts(products).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) listProducersByCoordinatorID(ctx context.Context, coordinatorID string) (service.Producers, error) {
	if coordinatorID == "" {
		return service.Producers{}, nil
	}
	in := &user.ListProducersInput{
		CoordinatorID: coordinatorID,
	}
	producers, _, err := h.user.ListProducers(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProducers(producers), nil
}

func (h *handler) multiGetProducers(ctx context.Context, producerIDs []string) (service.Producers, error) {
	if len(producerIDs) == 0 {
		return service.Producers{}, nil
	}
	in := &user.MultiGetProducersInput{
		ProducerIDs: producerIDs,
	}
	producers, err := h.user.MultiGetProducers(ctx, in)
	if err != nil {
		return nil, err
	}
	if len(producers) == 0 {
		return service.Producers{}, nil
	}
	return service.NewProducers(producers), nil
}

func (h *handler) getProducer(ctx context.Context, producerID string) (*service.Producer, error) {
	in := &user.GetProducerInput{
		ProducerID: producerID,
	}
	producer, err := h.user.GetProducer(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProducer(producer), nil
}
