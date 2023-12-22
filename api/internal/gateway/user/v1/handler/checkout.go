package handler

import (
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) checkoutRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/checkouts", h.authentication)

	r.POST("", h.Checkout)
	r.GET("/:transactionId", h.GetCheckoutState)
}

func (h *handler) Checkout(ctx *gin.Context) {
	req := &request.CheckoutRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	methodType := service.PaymentMethodType(req.PaymentMethod)
	system, err := h.getPaymentSystem(ctx, methodType)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if !system.InService() {
		h.forbidden(ctx, errors.New("handler: out of service"))
		return
	}
	detail := &store.CheckoutDetail{
		UserID:            getUserID(ctx),
		SessionID:         h.getSessionID(ctx),
		RequestID:         req.RequestID,
		CoordinatorID:     req.CoordinatorID,
		BoxNumber:         req.BoxNumber,
		PromotionCode:     req.PromotionCode,
		BillingAddressID:  req.BillingAddressID,
		ShippingAddressID: req.ShippingAddressID,
		CallbackURL:       req.CallbackURL,
		Total:             req.Total,
	}
	var redirectURL string
	switch methodType {
	case service.PaymentMethodTypeCreditCard:
		if req.CreditCard == nil {
			h.badRequest(ctx, errors.New("handler: credit card is required"))
			break
		}
		in := &store.CheckoutCreditCardInput{
			CheckoutDetail:    *detail,
			Number:            req.CreditCard.Number,
			Month:             req.CreditCard.Month,
			Year:              req.CreditCard.Year,
			VerificationValue: req.CreditCard.VerificationValue,
		}
		redirectURL, err = h.store.CheckoutCreditCard(ctx, in)
	case service.PaymentMethodTypePayPay:
		in := &store.CheckoutPayPayInput{
			CheckoutDetail: *detail,
		}
		redirectURL, err = h.store.CheckoutPayPay(ctx, in)
	case service.PaymentMethodTypeLinePay:
		in := &store.CheckoutLinePayInput{
			CheckoutDetail: *detail,
		}
		redirectURL, err = h.store.CheckoutLinePay(ctx, in)
	case service.PaymentMethodTypeMerpay:
		in := &store.CheckoutMerpayInput{
			CheckoutDetail: *detail,
		}
		redirectURL, err = h.store.CheckoutMerpay(ctx, in)
	case service.PaymentMethodTypeRakutenPay:
		in := &store.CheckoutRakutenPayInput{
			CheckoutDetail: *detail,
		}
		redirectURL, err = h.store.CheckoutRakutenPay(ctx, in)
	case service.PaymentMethodTypeAUPay:
		in := &store.CheckoutAUPayInput{
			CheckoutDetail: *detail,
		}
		redirectURL, err = h.store.CheckoutAUPay(ctx, in)
	default:
		err := errors.New("handler: not implemented payment method")
		h.httpError(ctx, status.Error(codes.Unimplemented, err.Error()))
		return
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.CheckoutResponse{
		URL: redirectURL,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetCheckoutState(ctx *gin.Context) {
	in := &store.GetCheckoutStateInput{
		UserID:        getUserID(ctx),
		TransactionID: util.GetParam(ctx, "transactionId"),
	}
	orderID, status, err := h.store.GetCheckoutState(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.CheckoutStateResponse{
		OrderID: orderID,
		Status:  service.NewPaymentStatus(status).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
