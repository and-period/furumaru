package handler

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Shop
// @tag.description ショップ関連
func (h *handler) shopRotues(rg *gin.RouterGroup) {
	r := rg.Group("/shops", h.authentication)

	r.GET("/:shopId", h.filterAccessShop, h.GetShop)
	r.PATCH("/:shopId", h.filterAccessShop, h.UpdateShop)
}

func (h *handler) filterAccessShop(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			shop, err := h.getShop(ctx, util.GetParam(ctx, "shopId"))
			if err != nil {
				return false, err
			}
			return currentAdmin(ctx, shop.CoordinatorID), nil
		},
		producer: func(ctx *gin.Context) (bool, error) {
			shop, err := h.getShop(ctx, util.GetParam(ctx, "shopId"))
			if err != nil {
				return false, err
			}
			return slices.Contains(shop.ProducerIDs, getAdminID(ctx)), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Next()
}

// @Summary     ショップ取得
// @Description 指定されたショップの詳細情報を取得します。コーディネーター、生産者、商品種別情報も含まれます。
// @Tags        Shop
// @Router      /v1/shops/{shopId} [get]
// @Security    bearerauth
// @Param       shopId path string true "ショップID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} response.ShopResponse
// @Failure     403 {object} util.ErrorResponse "ショップの参照権限がない"
// @Failure     404 {object} util.ErrorResponse "ショップが存在しない"
func (h *handler) GetShop(ctx *gin.Context) {
	shop, err := h.getShop(ctx, ctx.Param("shopId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		coordinator  *service.Coordinator
		producers    service.Producers
		productTypes service.ProductTypes
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ectx, shop.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		producers, err = h.multiGetProducers(ectx, shop.ProducerIDs)
		return
	})
	eg.Go(func() (err error) {
		productTypes, err = h.multiGetProductTypes(ectx, shop.ProductTypeIDs)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ShopResponse{
		Shop:         shop.Response(),
		Coordinator:  coordinator.Response(),
		Producers:    producers.Response(),
		ProductTypes: productTypes.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     ショップ更新
// @Description ショップの情報を更新します。ショップ名、商品種別、営業日を変更できます。
// @Tags        Shop
// @Router      /v1/shops/{shopId} [patch]
// @Security    bearerauth
// @Param       shopId path string true "ショップID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body request.UpdateShopRequest true "ショップ情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "ショップの更新権限がない"
// @Failure     404 {object} util.ErrorResponse "ショップが存在しない"
func (h *handler) UpdateShop(ctx *gin.Context) {
	req := &request.UpdateShopRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	productTypes, err := h.multiGetProductTypes(ctx, req.ProductTypeIDs)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(productTypes) != len(req.ProductTypeIDs) {
		h.badRequest(ctx, errors.New("handler: unmatch product types length"))
		return
	}

	in := &store.UpdateShopInput{
		ShopID:         ctx.Param("shopId"),
		Name:           req.Name,
		ProductTypeIDs: req.ProductTypeIDs,
		BusinessDays:   req.BusinessDays,
	}
	if err := h.store.UpdateShop(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) listShopsByCoordinatorIDs(ctx context.Context, coordinatorIDs []string) (service.Shops, error) {
	in := &store.ListShopsInput{
		CoordinatorIDs: coordinatorIDs,
		NoLimit:        true,
	}
	shops, _, err := h.store.ListShops(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewShops(shops), nil
}

func (h *handler) listShopsByProducerIDs(ctx context.Context, producerIDs []string) (service.Shops, error) {
	in := &store.ListShopsInput{
		ProducerIDs: producerIDs,
		NoLimit:     true,
	}
	shops, _, err := h.store.ListShops(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewShops(shops), nil
}

func (h *handler) multiGetShops(ctx context.Context, shopIDs []string) (service.Shops, error) {
	if len(shopIDs) == 0 {
		return service.Shops{}, nil
	}
	in := &store.MultiGetShopsInput{
		ShopIDs: shopIDs,
	}
	shops, err := h.store.MultiGetShops(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewShops(shops), nil
}

func (h *handler) getShop(ctx context.Context, shopID string) (*service.Shop, error) {
	in := &store.GetShopInput{
		ShopID: shopID,
	}
	shop, err := h.store.GetShop(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewShop(shop), nil
}

func (h *handler) getShopByCoordinatorID(ctx context.Context, coordinatorID string) (*service.Shop, error) {
	in := &store.GetShopByCoordinatorIDInput{
		CoordinatorID: coordinatorID,
	}
	shop, err := h.store.GetShopByCoordinatorID(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewShop(shop), nil
}
