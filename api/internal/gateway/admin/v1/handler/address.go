package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/user"
)

func (h *handler) multiGetAddressesByRevision(
	ctx context.Context,
	revisionIDs []int64,
) (service.Addresses, error) {
	if len(revisionIDs) == 0 {
		return service.Addresses{}, nil
	}
	in := &user.MultiGetAddressesByRevisionInput{
		AddressRevisionIDs: revisionIDs,
	}
	addresses, err := h.user.MultiGetAddressesByRevision(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewAddresses(addresses), nil
}
