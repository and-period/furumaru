package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
	"github.com/and-period/furumaru/api/internal/store"
)

func (h *handler) multiGetProductTags(ctx context.Context, productTagIDs []string) (service.ProductTags, error) {
	if len(productTagIDs) == 0 {
		return service.ProductTags{}, nil
	}
	in := &store.MultiGetProductTagsInput{
		ProductTagIDs: productTagIDs,
	}
	productTags, err := h.store.MultiGetProductTags(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProductTags(productTags), nil
}
