package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
	"github.com/and-period/furumaru/api/internal/store"
)

func (h *handler) getPaymentSystem(ctx context.Context, methodType service.PaymentMethodType) (*service.PaymentSystem, error) {
	in := &store.GetPaymentSystemInput{
		MethodType: methodType.StoreEntity(),
	}
	system, err := h.store.GetPaymentSystem(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewPaymentSystem(system), nil
}
