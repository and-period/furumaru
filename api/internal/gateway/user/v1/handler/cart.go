package handler

import (
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Cart
// @tag.description カート関連
func (h *handler) cartRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/carts")

	r.GET("", h.GetCart)
	r.GET("/:coordinatorId", h.CalcCart)
	r.POST("/-/items", h.AddCartItem)
	r.DELETE("/-/items/:productId", h.RemoveCartItem)
}

// @Summary     買い物かご取得
// @Description 買い物かごの内容を取得します。
// @Tags        Cart
// @Router      /carts [get]
// @Security    cookieauth
// @Produce     json
// @Success     200 {object} response.CartResponse
func (h *handler) GetCart(ctx *gin.Context) {
	sessionID := h.getSessionID(ctx)
	if sessionID == "" {
		// セッションIDがない場合は空のレスポンスを返す
		res := &response.CartResponse{
			Carts:        []*response.Cart{},
			Coordinators: []*response.Coordinator{},
			Products:     []*response.Product{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	in := &store.GetCartInput{
		SessionID: sessionID,
	}
	cart, err := h.store.GetCart(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if cart == nil {
		h.notFound(ctx, errNotFoundCart)
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
		products, err = h.multiGetProducts(ectx, cart.Baskets.ProductIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.CartResponse{
		Carts:        service.NewCarts(cart).Response(),
		Coordinators: coordinators.Response(),
		Products:     products.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     買い物かご計算
// @Description 指定されたコーディネータの買い物かごの送料や合計金額を計算します。
// @Tags        Cart
// @Router      /carts/{coordinatorId} [get]
// @Security    cookieauth
// @Param       coordinatorId path string true "コーディネータID"
// @Param       number query int64 false "箱数"
// @Param       prefecture query int32 false "都道府県コード"
// @Param       promotion query string false "プロモーションコード"
// @Produce     json
// @Success     200 {object} response.CalcCartResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CalcCart(ctx *gin.Context) {
	boxNumber, err := util.GetQueryInt64(ctx, "number", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	prefectureCode, err := util.GetQueryInt32(ctx, "prefecture", 0)
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
			SessionID:      h.getSessionID(ctx),
			CoordinatorID:  coordinatorID,
			BoxNumber:      boxNumber,
			PromotionCode:  promotionCode,
			PrefectureCode: prefectureCode,
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
	products, err := h.multiGetProducts(ctx, items.ProductIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.CalcCartResponse{
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

// @Summary     買い物かごに商品追加
// @Description 買い物かごに商品を追加します。
// @Tags        Cart
// @Router      /carts/-/items [post]
// @Security    cookieauth
// @Accept      json
// @Param       request body request.AddCartItemRequest true "商品追加"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "商品が非公開になっている"
// @Failure     404 {object} util.ErrorResponse "商品が存在しない"
// @Failure     412 {object} util.ErrorResponse "商品在庫が不足している"
func (h *handler) AddCartItem(ctx *gin.Context) {
	req := &request.AddCartItemRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.AddCartItemInput{
		SessionID: h.getSessionID(ctx),
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

// @Summary     買い物かごから商品削除
// @Description 買い物かごから商品を削除します。
// @Tags        Cart
// @Router      /carts/-/items/{productId} [delete]
// @Security    cookieauth
// @Param       productId path string true "商品ID"
// @Param       number query int64 false "箱数"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) RemoveCartItem(ctx *gin.Context) {
	boxNumber, err := util.GetQueryInt64(ctx, "number", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.RemoveCartItemInput{
		SessionID: h.getSessionID(ctx),
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
