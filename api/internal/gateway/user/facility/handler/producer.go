package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
	"github.com/and-period/furumaru/api/internal/user"
)

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
	shopsIn := &user.ListShopsInput{
		ProducerIDs: producerIDs,
		NoLimit:     true,
	}
	shops, _, err := h.user.ListShops(ctx, shopsIn)
	if err != nil {
		return nil, err
	}
	return service.NewProducers(producers, shops.GroupByProducerID()), nil
}

func (h *handler) getProducer(ctx context.Context, producerID string) (*service.Producer, error) {
	in := &user.GetProducerInput{
		ProducerID: producerID,
	}
	producer, err := h.user.GetProducer(ctx, in)
	if err != nil {
		return nil, err
	}
	shopsIn := &user.ListShopsInput{
		ProducerIDs: []string{producerID},
		NoLimit:     true,
	}
	shops, _, err := h.user.ListShops(ctx, shopsIn)
	if err != nil {
		return nil, err
	}
	return service.NewProducer(producer, shops), nil
}
