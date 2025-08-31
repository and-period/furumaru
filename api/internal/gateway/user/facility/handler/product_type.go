package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
	"github.com/and-period/furumaru/api/internal/store"
)

func (h *handler) multiGetProductTypes(ctx context.Context, productTypeIDs []string) (service.ProductTypes, error) {
	if len(productTypeIDs) == 0 {
		return service.ProductTypes{}, nil
	}
	in := &store.MultiGetProductTypesInput{
		ProductTypeIDs: productTypeIDs,
	}
	sproductTypes, err := h.store.MultiGetProductTypes(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProductTypes(sproductTypes), nil
}

func (h *handler) getProductType(ctx context.Context, productTypeID string) (*service.ProductType, error) {
	in := &store.GetProductTypeInput{
		ProductTypeID: productTypeID,
	}
	sproductType, err := h.store.GetProductType(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProductType(sproductType), nil
}
