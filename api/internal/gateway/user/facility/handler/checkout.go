package handler

import (
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @tag.name        Checkout
// @tag.description チェックアウト関連
func (h *handler) checkoutRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/checkouts", h.authentication)

	r.POST("", h.Checkout)
	r.GET("/:transactionId", h.GetCheckoutState)
}

// @Summary     購入する
// @Description 商品を購入します。
// @Tags        Checkout
// @Router      /facilities/{facilityId}/checkouts [post]
// @Param       facilityId path string true "施設ID"
// @Security    bearerauth
// @Accept      json
// @Param				request body types.CheckoutRequest true "チェックアウト情報"
// @Produce     json
// @Success     200 {object} types.CheckoutResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) Checkout(ctx *gin.Context) {
	req := &types.CheckoutRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	var (
		auth     *uentity.User
		producer *service.Producer
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &user.GetUserInput{
			UserID: h.getUserID(ctx),
		}
		auth, err = h.user.GetUser(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, h.getProducerID(ctx))
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	methodType := service.PaymentMethodType(req.PaymentMethod)
	if err := h.checkPaymentSystem(ctx, methodType); err != nil {
		return
	}
	detail := &store.CheckoutDetail{
		Type:          sentity.OrderTypeProduct,
		UserID:        h.getUserID(ctx),
		SessionID:     h.getUserID(ctx),
		RequestID:     req.RequestID,
		PromotionCode: req.PromotionCode,
		CallbackURL:   req.CallbackURL,
		Total:         req.Total,
		CheckoutProductDetail: store.CheckoutProductDetail{
			CoordinatorID:  req.CoordinatorID,
			BoxNumber:      req.BoxNumber,
			Pickup:         true, // 施設での受取に限定する
			PickupAt:       auth.FacilityUser.LastCheckInAt,
			PickupLocation: producer.Username,
		},
	}
	params := &checkoutParams{
		methodType: methodType,
		detail:     detail,
		creditCard: req.CreditCard,
	}
	h.checkout(ctx, params)
}

// @Summary     支払い状態取得
// @Description 支払い状態を取得します。
// @Tags        Checkout
// @Router      /facilities/{facilityId}/checkouts/{transactionId} [get]
// @Param       facilityId path string true "施設ID"
// @Param       transactionId path string true "取引ID"
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} types.CheckoutStateResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     404 {object} util.ErrorResponse "取引が見つからない"
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
