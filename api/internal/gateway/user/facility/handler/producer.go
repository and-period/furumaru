package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
)

func (h *handler) getProducer(ctx context.Context, producerID string) (*service.Producer, error) {
	in := &user.GetProducerInput{
		ProducerID: producerID,
	}
	producer, err := h.user.GetProducer(ctx, in)
	if err != nil {
		return nil, err
	}
	shopsIn := &store.ListShopsInput{
		ProducerIDs: []string{producerID},
		NoLimit:     true,
	}
	shops, _, err := h.store.ListShops(ctx, shopsIn)
	if err != nil {
		return nil, err
	}
	return service.NewProducer(producer, shops), nil
}
