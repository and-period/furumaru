package handler

import (
	"context"
	"errors"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/store"
)

func (h *handler) getShippingByShopID(ctx context.Context, shopID string) (*service.Shipping, error) {
	in := &store.GetShippingByShopIDInput{
		ShopID: shopID,
	}
	shipping, err := h.store.GetShippingByShopID(ctx, in)
	if errors.Is(err, exception.ErrNotFound) {
		// 配送設定の登録をしていない場合、デフォルト配送設定を返却する
		in := &store.GetDefaultShippingInput{}
		shipping, err = h.store.GetDefaultShipping(ctx, in)
	}
	if err != nil {
		return nil, err
	}
	return service.NewShipping(shipping), nil
}
