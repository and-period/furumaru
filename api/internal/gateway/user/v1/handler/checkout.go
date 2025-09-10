package handler

import (
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @tag.name        Checkout
// @tag.description 決済・注文関連
func (h *handler) checkoutRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/checkouts", h.authentication)

	r.POST("/products", h.CheckoutProduct)
	r.GET("/experiences/:experienceId", h.PreCheckoutExperience)
	r.POST("/experiences/:experienceId", h.CheckoutExperience)
	r.GET("/:transactionId", h.GetCheckoutState)
}

// @Summary     商品決済
// @Description 商品の決済を実行し、注文を作成します。
// @Tags        Checkout
// @Router      /checkouts/products [post]
// @Security    bearerauth
// @Security    cookieauth
// @Accept      json
// @Param       request body types.CheckoutProductRequest true "商品決済"
// @Produce     json
// @Success     200 {object} types.CheckoutResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     403 {object} util.ErrorResponse "決済システムがメンテナンス中 もしくは 店舗が利用停止中"
// @Failure     412 {object} util.ErrorResponse "前提条件エラー(商品在庫が不足、無効なプロモーションなど...)"
func (h *handler) CheckoutProduct(ctx *gin.Context) {
	req := &types.CheckoutProductRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	methodType := service.PaymentMethodType(req.PaymentMethod)
	if err := h.checkPaymentSystem(ctx, methodType); err != nil {
		return
	}
	detail := &store.CheckoutDetail{
		Type:             sentity.OrderTypeProduct,
		UserID:           h.getUserID(ctx),
		SessionID:        h.getSessionID(ctx),
		RequestID:        req.RequestID,
		PromotionCode:    req.PromotionCode,
		BillingAddressID: req.BillingAddressID,
		CallbackURL:      req.CallbackURL,
		Total:            req.Total,
		CheckoutProductDetail: store.CheckoutProductDetail{
			CoordinatorID:     req.CoordinatorID,
			BoxNumber:         req.BoxNumber,
			ShippingAddressID: req.ShippingAddressID,
		},
	}
	params := &checkoutParams{
		methodType: methodType,
		detail:     detail,
		creditCard: req.CreditCard,
	}
	h.checkout(ctx, params)
}

// @Summary     体験事前決済情報取得
// @Description 体験を決済する前に必要な情報を取得します。
// @Tags        Checkout
// @Router      /checkouts/experiences/{experienceId} [get]
// @Security    bearerauth
// @Security    cookieauth
// @Param       experienceId path string true "体験ID"
// @Produce     json
// @Success     200 {object} types.PreCheckoutExperienceResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     404 {object} util.ErrorResponse "体験が見つからない"
func (h *handler) PreCheckoutExperience(ctx *gin.Context) {
	adultCount, err := util.GetQueryInt64(ctx, "adult", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	juniorHighSchoolCount, err := util.GetQueryInt64(ctx, "juniorHighSchool", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	elementarySchoolCount, err := util.GetQueryInt64(ctx, "elementarySchool", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	preschoolCount, err := util.GetQueryInt64(ctx, "preschool", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	seniorCount, err := util.GetQueryInt64(ctx, "senior", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	promotionCode := util.GetQuery(ctx, "promotion", "")

	var (
		experience *service.Experience
		promotion  *service.Promotion
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		experience, err = h.getExperience(ectx, ctx.Param("experienceId"))
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

	params := &service.CalcExperienceParams{
		AdultCount:            adultCount,
		JuniorHighSchoolCount: juniorHighSchoolCount,
		ElementarySchoolCount: elementarySchoolCount,
		PreschoolCount:        preschoolCount,
		SeniorCount:           seniorCount,
		Promotion:             promotion,
	}
	subtotal, discount := experience.Calc(params)

	res := &types.PreCheckoutExperienceResponse{
		RequestID:  h.generateID(),
		Experience: experience.Response(),
		Promotion:  promotion.Response(),
		SubTotal:   subtotal,
		Discount:   discount,
		Total:      subtotal - discount,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     体験決済
// @Description 体験の決済を実行し、予約を作成します。
// @Tags        Checkout
// @Router      /checkouts/experiences/{experienceId} [post]
// @Security    bearerauth
// @Security    cookieauth
// @Param       experienceId path string true "体験ID"
// @Accept      json
// @Param       request body types.CheckoutExperienceRequest true "体験決済"
// @Produce     json
// @Success     200 {object} types.CheckoutResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) CheckoutExperience(ctx *gin.Context) {
	req := &types.CheckoutExperienceRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	methodType := service.PaymentMethodType(req.PaymentMethod)
	if err := h.checkPaymentSystem(ctx, methodType); err != nil {
		return
	}
	detail := &store.CheckoutDetail{
		Type:             sentity.OrderTypeExperience,
		UserID:           h.getUserID(ctx),
		SessionID:        h.getSessionID(ctx),
		RequestID:        req.RequestID,
		PromotionCode:    req.PromotionCode,
		BillingAddressID: req.BillingAddressID,
		CallbackURL:      req.CallbackURL,
		Total:            req.Total,
		CheckoutExperienceDetail: store.CheckoutExperienceDetail{
			ExperienceID:          util.GetParam(ctx, "experienceId"),
			AdultCount:            req.AdultCount,
			JuniorHighSchoolCount: req.JuniorHighSchoolCount,
			ElementarySchoolCount: req.ElementarySchoolCount,
			PreschoolCount:        req.PreschoolCount,
			SeniorCount:           req.SeniorCount,
			Transportation:        req.Transportation,
			RequestedDate:         req.RequestedDate,
			RequestedTime:         req.RequestedTime,
		},
	}
	params := &checkoutParams{
		methodType: methodType,
		detail:     detail,
		creditCard: req.CreditCard,
	}
	h.checkout(ctx, params)
}

func (h *handler) checkPaymentSystem(ctx *gin.Context, methodType service.PaymentMethodType) error {
	system, err := h.getPaymentSystem(ctx, methodType)
	if err != nil {
		h.httpError(ctx, err)
		return err
	}
	if !system.InService() {
		err := errors.New("handler: out of service")
		h.forbidden(ctx, err)
		return err
	}
	return nil
}

type checkoutParams struct {
	methodType service.PaymentMethodType
	detail     *store.CheckoutDetail
	creditCard *types.CheckoutCreditCard
}

func (h *handler) checkout(ctx *gin.Context, params *checkoutParams) {
	var (
		redirectURL string
		err         error
	)
	switch types.PaymentMethodType(params.methodType) {
	case types.PaymentMethodTypeCreditCard:
		if params.creditCard == nil {
			h.badRequest(ctx, errors.New("handler: credit card is required"))
			break
		}
		in := &store.CheckoutCreditCardInput{
			CheckoutDetail:    *params.detail,
			Name:              params.creditCard.Name,
			Number:            params.creditCard.Number,
			Month:             params.creditCard.Month,
			Year:              params.creditCard.Year,
			VerificationValue: params.creditCard.VerificationValue,
		}
		redirectURL, err = h.store.CheckoutCreditCard(ctx, in)
	case types.PaymentMethodTypePayPay:
		in := &store.CheckoutPayPayInput{
			CheckoutDetail: *params.detail,
		}
		redirectURL, err = h.store.CheckoutPayPay(ctx, in)
	case types.PaymentMethodTypeLinePay:
		in := &store.CheckoutLinePayInput{
			CheckoutDetail: *params.detail,
		}
		redirectURL, err = h.store.CheckoutLinePay(ctx, in)
	case types.PaymentMethodTypeMerpay:
		in := &store.CheckoutMerpayInput{
			CheckoutDetail: *params.detail,
		}
		redirectURL, err = h.store.CheckoutMerpay(ctx, in)
	case types.PaymentMethodTypeRakutenPay:
		in := &store.CheckoutRakutenPayInput{
			CheckoutDetail: *params.detail,
		}
		redirectURL, err = h.store.CheckoutRakutenPay(ctx, in)
	case types.PaymentMethodTypeAUPay:
		in := &store.CheckoutAUPayInput{
			CheckoutDetail: *params.detail,
		}
		redirectURL, err = h.store.CheckoutAUPay(ctx, in)
	case types.PaymentMethodTypePaidy:
		in := &store.CheckoutPaidyInput{
			CheckoutDetail: *params.detail,
		}
		redirectURL, err = h.store.CheckoutPaidy(ctx, in)
	case types.PaymentMethodTypeBankTransfer:
		in := &store.CheckoutBankTransferInput{
			CheckoutDetail: *params.detail,
		}
		redirectURL, err = h.store.CheckoutBankTransfer(ctx, in)
	case types.PaymentMethodTypePayEasy:
		in := &store.CheckoutPayEasyInput{
			CheckoutDetail: *params.detail,
		}
		redirectURL, err = h.store.CheckoutPayEasy(ctx, in)
	case types.PaymentMethodTypeFree:
		in := &store.CheckoutFreeInput{
			CheckoutDetail: *params.detail,
		}
		redirectURL, err = h.store.CheckoutFree(ctx, in)
	default:
		err := errors.New("handler: not implemented payment method")
		h.httpError(ctx, status.Error(codes.Unimplemented, err.Error()))
		return
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.CheckoutResponse{
		URL: redirectURL,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     決済状態取得
// @Description 決済トランザクションの状態を取得します。
// @Tags        Checkout
// @Router      /checkouts/{transactionId} [get]
// @Security    bearerauth
// @Security    cookieauth
// @Param       transactionId path string true "トランザクションID"
// @Produce     json
// @Success     200 {object} types.CheckoutStateResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     404 {object} util.ErrorResponse "トランザクションが見つからない"
func (h *handler) GetCheckoutState(ctx *gin.Context) {
	in := &store.GetCheckoutStateInput{
		UserID:        h.getUserID(ctx),
		TransactionID: util.GetParam(ctx, "transactionId"),
	}
	orderID, status, err := h.store.GetCheckoutState(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.CheckoutStateResponse{
		OrderID: orderID,
		Status:  service.NewPaymentStatus(status).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
