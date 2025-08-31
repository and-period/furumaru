package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
	"github.com/and-period/furumaru/api/internal/store"
)

func (h *handler) aggregateProductRates(ctx context.Context, productIDs ...string) (service.ProductRates, error) {
	if len(productIDs) == 0 {
		return service.ProductRates{}, nil
	}
	in := &store.AggregateProductReviewsInput{
		ProductIDs: productIDs,
	}
	reviews, err := h.store.AggregateProductReviews(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProductRates(reviews), nil
}
