package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
)

func (h *handler) listShopsByCoordinatorIDs(ctx context.Context, coordinatorIDs []string) (sentity.Shops, error) {
	in := &store.ListShopsInput{
		CoordinatorIDs: coordinatorIDs,
		NoLimit:        true,
	}
	shops, _, err := h.store.ListShops(ctx, in)
	return shops, err
}

func (h *handler) listShopsByProducerIDs(ctx context.Context, producerIDs []string) (sentity.Shops, error) {
	in := &store.ListShopsInput{
		ProducerIDs: producerIDs,
		NoLimit:     true,
	}
	shops, _, err := h.store.ListShops(ctx, in)
	return shops, err
}

func (h *handler) getShop(ctx context.Context, shopID string) (*sentity.Shop, error) {
	in := &store.GetShopInput{
		ShopID: shopID,
	}
	return h.store.GetShop(ctx, in)
}

func (h *handler) getShopByCoordinatorID(ctx context.Context, coordinatorID string) (*sentity.Shop, error) {
	in := &store.GetShopByCoordinatorIDInput{
		CoordinatorID: coordinatorID,
	}
	return h.store.GetShopByCoordinatorID(ctx, in)
}
