package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

func (h *handler) producerRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/producers")

	r.GET("", h.ListProducers)
	r.GET("/:producerId", h.GetProducer)
}

func (h *handler) ListProducers(ctx *gin.Context) {
	res := &response.ProducersResponse{
		Producers: []*response.Producer{},
		Total:     0,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetProducer(ctx *gin.Context) {
	res := &response.ProducerResponse{
		Producer: &response.Producer{},
		Products: []*response.Product{},
		Lives:    []*response.LiveSummary{},
		Archives: []*response.ArchiveSummary{},
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
