package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

func (h *handler) listShopsByCoordinatorIDs(ctx context.Context, coordinatorIDs []string) (entity.Shops, error) {
	in := &user.ListShopsInput{
		CoordinatorIDs: coordinatorIDs,
		NoLimit:        true,
	}
	shops, _, err := h.user.ListShops(ctx, in)
	return shops, err
}

func (h *handler) getShopByCoordinatorID(ctx context.Context, coordinatorID string) (*entity.Shop, error) {
	in := &user.GetShopByCoordinatorIDInput{
		CoordinatorID: coordinatorID,
	}
	return h.user.GetShopByCoordinatorID(ctx, in)
}
