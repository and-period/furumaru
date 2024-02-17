package handler

import (
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) guestCheckoutRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/guests/checkouts")

	r.POST("", h.GuestCheckout)
	r.GET("/:transactionId", h.GetGuestCheckoutState)
}

func (h *handler) GuestCheckout(ctx *gin.Context) {
	req := &request.GuestCheckoutRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	if req.BillingAddress == nil {
		h.badRequest(ctx, errors.New("handler: billing address is required"))
		return
	}
	if !req.IsSameAddress && req.ShippingAddress == nil {
		h.badRequest(ctx, errors.New("handler: shipping address is required"))
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
	// ゲストユーザー登録
	guestIn := &user.UpsertGuestInput{
		Lastname:      req.BillingAddress.Lastname,
		Firstname:     req.BillingAddress.Firstname,
		LastnameKana:  req.BillingAddress.LastnameKana,
		FirstnameKana: req.BillingAddress.FirstnameKana,
		Email:         req.Email,
	}
	userID, err := h.user.UpsertGuest(ctx, guestIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	// 請求先住所登録
	baddressIn := &user.CreateAddressInput{
		UserID:         userID,
		Lastname:       req.BillingAddress.Lastname,
		Firstname:      req.BillingAddress.Firstname,
		LastnameKana:   req.BillingAddress.LastnameKana,
		FirstnameKana:  req.BillingAddress.FirstnameKana,
		PostalCode:     req.BillingAddress.PostalCode,
		PrefectureCode: req.BillingAddress.PrefectureCode,
		City:           req.BillingAddress.City,
		AddressLine1:   req.BillingAddress.AddressLine1,
		AddressLine2:   req.BillingAddress.AddressLine2,
		PhoneNumber:    req.BillingAddress.PhoneNumber,
		IsDefault:      true,
	}
	baddress, err := h.user.CreateAddress(ctx, baddressIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	billingAddressID := baddress.ID
	var shippingAddressID string
	if req.IsSameAddress {
		shippingAddressID = billingAddressID
	} else {
		saddressIn := &user.CreateAddressInput{
			UserID:         userID,
			Lastname:       req.BillingAddress.Lastname,
			Firstname:      req.BillingAddress.Firstname,
			LastnameKana:   req.BillingAddress.LastnameKana,
			FirstnameKana:  req.BillingAddress.FirstnameKana,
			PostalCode:     req.BillingAddress.PostalCode,
			PrefectureCode: req.BillingAddress.PrefectureCode,
			City:           req.BillingAddress.City,
			AddressLine1:   req.BillingAddress.AddressLine1,
			AddressLine2:   req.BillingAddress.AddressLine2,
			PhoneNumber:    req.BillingAddress.PhoneNumber,
			IsDefault:      true,
		}
		saddress, err := h.user.CreateAddress(ctx, saddressIn)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		shippingAddressID = saddress.ID
	}
	// 購入処理を進める
	detail := &store.CheckoutDetail{
		UserID:            userID,
		SessionID:         h.getSessionID(ctx),
		RequestID:         req.RequestID,
		CoordinatorID:     req.CoordinatorID,
		BoxNumber:         req.BoxNumber,
		PromotionCode:     req.PromotionCode,
		BillingAddressID:  billingAddressID,
		ShippingAddressID: shippingAddressID,
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
			Name:              req.CreditCard.Name,
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
	res := &response.GuestCheckoutResponse{
		URL: redirectURL,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetGuestCheckoutState(ctx *gin.Context) {
	in := &store.GetCheckoutStateInput{
		SessionID:     h.getSessionID(ctx),
		TransactionID: util.GetParam(ctx, "transactionId"),
	}
	orderID, status, err := h.store.GetCheckoutState(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.GuestCheckoutStateResponse{
		OrderID: orderID,
		Status:  service.NewPaymentStatus(status).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
