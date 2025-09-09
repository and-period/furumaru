package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

// @tag.name        GuestCheckout
// @tag.description ゲストチェックアウト関連
func (h *handler) guestCheckoutRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/guests/checkouts")

	r.POST("/products", h.GuestCheckoutProduct)
	r.GET("/experiences/:experienceId", h.PreCheckoutExperience)
	r.POST("/experiences/:experienceId", h.GuestCheckoutExperience)
	r.GET("/:transactionId", h.GetGuestCheckoutState)
}

// @Summary     ゲスト商品決済
// @Description ゲストユーザーとして商品の決済を実行し、注文を作成します。
// @Tags        GuestCheckout
// @Router      /guests/checkouts/products [post]
// @Security    cookieauth
// @Accept      json
// @Param       request body types.GuestCheckoutProductRequest true "ゲスト商品決済"
// @Produce     json
// @Success     200 {object} types.CheckoutResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "決済システムがメンテナンス中"
// @Failure     412 {object} util.ErrorResponse "前提条件エラー(商品在庫が不足、無効なプロモーションなど...)"
func (h *handler) GuestCheckoutProduct(ctx *gin.Context) {
	req := &types.GuestCheckoutProductRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	methodType := service.PaymentMethodType(req.PaymentMethod)
	if err := h.checkPaymentSystem(ctx, methodType); err != nil {
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
	// ゲストユーザー登録
	userID, billingAddressID, err := h.createGuestForCheckout(ctx, req.Email, req.BillingAddress)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	// 配送先住所登録
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
		Type:             sentity.OrderTypeProduct,
		UserID:           userID,
		SessionID:        h.getSessionID(ctx),
		RequestID:        req.RequestID,
		PromotionCode:    req.PromotionCode,
		BillingAddressID: billingAddressID,
		CallbackURL:      req.CallbackURL,
		Total:            req.Total,
		CheckoutProductDetail: store.CheckoutProductDetail{
			CoordinatorID:     req.CoordinatorID,
			BoxNumber:         req.BoxNumber,
			ShippingAddressID: shippingAddressID,
		},
	}
	params := &checkoutParams{
		methodType: methodType,
		detail:     detail,
		creditCard: req.CreditCard,
	}
	h.checkout(ctx, params)
}

// @Summary     ゲスト体験決済
// @Description ゲストユーザーとして体験の決済を実行し、予約を作成します。
// @Tags        GuestCheckout
// @Router      /guests/checkouts/experiences/{experienceId} [post]
// @Security    cookieauth
// @Param       experienceId path string true "体験ID"
// @Accept      json
// @Param       request body types.GuestCheckoutExperienceRequest true "ゲスト体験決済"
// @Produce     json
// @Success     200 {object} types.CheckoutResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) GuestCheckoutExperience(ctx *gin.Context) {
	req := &types.GuestCheckoutExperienceRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	methodType := service.PaymentMethodType(req.PaymentMethod)
	if err := h.checkPaymentSystem(ctx, methodType); err != nil {
		return
	}
	if req.BillingAddress == nil {
		h.badRequest(ctx, errors.New("handler: billing address is required"))
		return
	}
	// ゲストユーザー登録
	userID, billingAddressID, err := h.createGuestForCheckout(ctx, req.Email, req.BillingAddress)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	// 購入処理を進める
	detail := &store.CheckoutDetail{
		Type:             sentity.OrderTypeExperience,
		UserID:           userID,
		SessionID:        h.getSessionID(ctx),
		RequestID:        req.RequestID,
		PromotionCode:    req.PromotionCode,
		BillingAddressID: billingAddressID,
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

func (h *handler) createGuestForCheckout(ctx context.Context, email string, address *types.GuestCheckoutAddress) (string, string, error) {
	// ゲストユーザー登録
	guestIn := &user.UpsertGuestInput{
		Lastname:      address.Lastname,
		Firstname:     address.Firstname,
		LastnameKana:  address.LastnameKana,
		FirstnameKana: address.FirstnameKana,
		Email:         email,
	}
	userID, err := h.user.UpsertGuest(ctx, guestIn)
	if err != nil {
		return "", "", err
	}
	// 請求先住所登録
	baddressIn := &user.CreateAddressInput{
		UserID:         userID,
		Lastname:       address.Lastname,
		Firstname:      address.Firstname,
		LastnameKana:   address.LastnameKana,
		FirstnameKana:  address.FirstnameKana,
		PostalCode:     address.PostalCode,
		PrefectureCode: address.PrefectureCode,
		City:           address.City,
		AddressLine1:   address.AddressLine1,
		AddressLine2:   address.AddressLine2,
		PhoneNumber:    address.PhoneNumber,
		IsDefault:      true,
	}
	baddress, err := h.user.CreateAddress(ctx, baddressIn)
	if err != nil {
		return "", "", err
	}
	return userID, baddress.ID, nil
}

// @Summary     ゲスト決済状態取得
// @Description ゲストユーザーの決済トランザクション状態を取得します。
// @Tags        GuestCheckout
// @Router      /guests/checkouts/{transactionId} [get]
// @Security    cookieauth
// @Param       transactionId path string true "トランザクションID"
// @Produce     json
// @Success     200 {object} types.CheckoutStateResponse
// @Failure     404 {object} util.ErrorResponse "トランザクションが見つからない"
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
	res := &types.GuestCheckoutStateResponse{
		OrderID: orderID,
		Status:  service.NewPaymentStatus(status).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
