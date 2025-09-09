package handler

import (
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Cart
// @tag.description カート関連
func (h *handler) cartRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/carts", h.authentication)

	r.GET("", h.GetCart)
	r.GET("/:coordinatorId", h.CalcCart)
	r.POST("/-/items", h.AddCartItem)
	r.DELETE("/-/items/:productId", h.RemoveCartItem)
}

// @Summary     カート取得
// @Description カートの内容を取得します。
// @Tags        Cart
// @Router      /facilities/{facilityId}/carts [get]
// @Param       facilityId path string true "施設ID"
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} types.CartResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) GetCart(ctx *gin.Context) {
	in := &store.GetCartInput{
		SessionID: h.getUserID(ctx),
	}
	cart, err := h.store.GetCart(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		coordinators service.Coordinators
		products     service.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinators, err = h.multiGetCoordinators(ectx, cart.Baskets.CoordinatorID())
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, h.getProducerID(ctx), cart.Baskets.ProductIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.CartResponse{
		Carts:        service.NewCarts(cart).Response(),
		Coordinators: coordinators.Response(),
		Products:     products.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     カート計算
// @Description カートの内容を計算します。
// @Tags        Cart
// @Router      /facilities/{facilityId}/carts/{coordinatorId} [get]
// @Param       facilityId    path string true "施設ID"
// @Param       coordinatorId path string true "コーディネーターID"
// @Param       number query int64 false "箱数"
// @Param       promotion query string false "プロモーションコード"
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} types.CalcCartResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) CalcCart(ctx *gin.Context) {
	boxNumber, err := util.GetQueryInt64(ctx, "number", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	promotionCode := util.GetQuery(ctx, "promotion", "")
	coordinatorID := util.GetParam(ctx, "coordinatorId")

	var (
		cart        *entity.Cart
		summary     *entity.OrderPaymentSummary
		coordinator *service.Coordinator
		promotion   *service.Promotion
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &store.CalcCartInput{
			SessionID:     h.getUserID(ctx),
			CoordinatorID: coordinatorID,
			BoxNumber:     boxNumber,
			PromotionCode: promotionCode,
			Pickup:        true, // 施設での受取に限定する
		}
		cart, summary, err = h.store.CalcCart(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ectx, coordinatorID)
		return
	})
	eg.Go(func() (err error) {
		if promotionCode == "" {
			return
		}
		promotion, err = h.getEnabledPromotion(ectx, promotionCode)
		if errors.Is(err, exception.ErrNotFound) {
			err = nil // エラーは返さず、プロモーション未適用状態で返す
		}
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	items := cart.Baskets.MergeByProductID()
	products, err := h.multiGetProducts(ctx, h.getProducerID(ctx), items.ProductIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.CalcCartResponse{
		RequestID:   h.generateID(),
		Carts:       service.NewCarts(cart).Response(),
		Items:       service.NewCartItems(items).Response(),
		Products:    products.Response(),
		Coordinator: coordinator.Response(),
		Promotion:   promotion.Response(),
		SubTotal:    summary.Subtotal,
		Discount:    summary.Discount,
		ShippingFee: summary.ShippingFee,
		Total:       summary.Total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     カートに追加
// @Description カートに商品を追加します。
// @Tags        Cart
// @Router      /facilities/{facilityId}/carts/-/items [post]
// @Param       facilityId path string true "施設ID"
// @Security    bearerauth
// @Accept      json
// @Param       request body types.AddCartItemRequest true "カートに追加リクエスト"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) AddCartItem(ctx *gin.Context) {
	req := &types.AddCartItemRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	if _, err := h.getProduct(ctx, h.getProducerID(ctx), req.ProductID); err != nil {
		h.httpError(ctx, err)
		return
	}
	in := &store.AddCartItemInput{
		SessionID: h.getUserID(ctx),
		UserID:    h.getUserID(ctx),
		UserAgent: ctx.Request.UserAgent(),
		ClientIP:  ctx.ClientIP(),
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}
	if err := h.store.AddCartItem(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     カートから削除
// @Description カートから商品を削除します。
// @Tags        Cart
// @Router      /facilities/{facilityId}/carts/-/items/{productId} [delete]
// @Param       facilityId path string true "施設ID"
// @Param       productId  path string true "商品ID"
// @Security    bearerauth
// @Produce     json
// @Success     204
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) RemoveCartItem(ctx *gin.Context) {
	boxNumber, err := util.GetQueryInt64(ctx, "number", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.RemoveCartItemInput{
		SessionID: h.getUserID(ctx),
		UserID:    h.getUserID(ctx),
		UserAgent: ctx.Request.UserAgent(),
		ClientIP:  ctx.ClientIP(),
		BoxNumber: boxNumber,
		ProductID: util.GetParam(ctx, "productId"),
	}
	if err := h.store.RemoveCartItem(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
