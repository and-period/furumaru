package handler

import (
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

func (h *handler) authRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/auth")

	r.GET("", h.GetAuth)
	r.POST("", h.SignIn)
	r.DELETE("", h.SignOut)
	r.GET("/google", h.authentication, h.AuthGoogleAccount)
	r.POST("/google", h.authentication, h.ConnectGoogleAccount)
	r.POST("/refresh-token", h.RefreshAuthToken)
	r.POST("/device", h.authentication, h.RegisterDevice)
	r.PATCH("/email", h.authentication, h.UpdateAuthEmail)
	r.POST("/email/verified", h.VerifyAuthEmail)
	r.PATCH("/password", h.authentication, h.UpdateAuthPassword)
	r.POST("/forgot-password", h.ForgotAuthPassword)
	r.POST("/forgot-password/verified", h.ResetAuthPassword)
	r.GET("/user", h.authentication, h.GetAuthUser)
	r.GET("/coordinator", h.authentication, h.GetAuthCoordinator)
	r.PATCH("/coordinator", h.authentication, h.UpdateAuthCoordinator)
	r.GET("/coordinator/shippings", h.authentication, h.GetAuthShipping)
	r.PATCH("/coordinator/shippings", h.authentication, h.UpsertAuthShipping)
}

func (h *handler) GetAuth(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}

	in := &user.GetAdminAuthInput{
		AccessToken: token,
	}
	auth, err := h.user.GetAdminAuth(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

type authUser interface {
	AuthUser() *service.AuthUser
}

func (h *handler) GetAuthUser(ctx *gin.Context) {
	adminID := getAdminID(ctx)
	var (
		auth authUser
		err  error
	)
	switch getAdminType(ctx) {
	case service.AdminTypeAdministrator:
		auth, err = h.getAdministrator(ctx, adminID)
	case service.AdminTypeCoordinator:
		auth, err = h.getCoordinator(ctx, adminID)
	case service.AdminTypeProducer:
		auth, err = h.getProducer(ctx, adminID)
	default:
		h.forbidden(ctx, errors.New("handler: unknown admin role"))
		return
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.AuthUserResponse{
		AuthUser: auth.AuthUser().Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) SignIn(ctx *gin.Context) {
	req := &request.SignInRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.SignInAdminInput{
		Key:      req.Username,
		Password: req.Password,
	}
	auth, err := h.user.SignInAdmin(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) AuthGoogleAccount(ctx *gin.Context) {
	in := &user.InitialGoogleAdminAuthInput{
		AdminID: getAdminID(ctx),
		State:   util.GetQuery(ctx, "state", ""),
	}
	authURL, err := h.user.InitialGoogleAdminAuth(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.AuthGoogleAccountResponse{
		URL: authURL,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) ConnectGoogleAccount(ctx *gin.Context) {
	req := &request.ConnectGoogleAccountRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.ConnectGoogleAdminAuthInput{
		AdminID: getAdminID(ctx),
		Code:    req.Code,
		Nonce:   req.Nonce,
	}
	if err := h.user.ConnectGoogleAdminAuth(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) SignOut(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}

	in := &user.SignOutAdminInput{
		AccessToken: token,
	}
	if err := h.user.SignOutAdmin(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) RefreshAuthToken(ctx *gin.Context) {
	req := &request.RefreshAuthTokenRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.RefreshAdminTokenInput{
		RefreshToken: req.RefreshToken,
	}
	auth, err := h.user.RefreshAdminToken(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) RegisterDevice(ctx *gin.Context) {
	req := &request.RegisterAuthDeviceRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.RegisterAdminDeviceInput{
		AdminID: getAdminID(ctx),
		Device:  req.Device,
	}
	if err := h.user.RegisterAdminDevice(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) UpdateAuthEmail(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	req := &request.UpdateAuthEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.UpdateAdminEmailInput{
		AccessToken: token,
		Email:       req.Email,
	}
	if err := h.user.UpdateAdminEmail(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) VerifyAuthEmail(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	req := &request.VerifyAuthEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.VerifyAdminEmailInput{
		AccessToken: token,
		VerifyCode:  req.VerifyCode,
	}
	if err := h.user.VerifyAdminEmail(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) UpdateAuthPassword(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	req := &request.UpdateAuthPasswordRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.UpdateAdminPasswordInput{
		AccessToken:          token,
		OldPassword:          req.OldPassword,
		NewPassword:          req.NewPassword,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	if err := h.user.UpdateAdminPassword(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) ForgotAuthPassword(ctx *gin.Context) {
	req := &request.ForgotAuthPasswordRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.ForgotAdminPasswordInput{
		Email: req.Email,
	}
	if err := h.user.ForgotAdminPassword(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) ResetAuthPassword(ctx *gin.Context) {
	req := &request.ResetAuthPasswordRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.VerifyAdminPasswordInput{
		Email:                req.Email,
		VerifyCode:           req.VerifyCode,
		NewPassword:          req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	if err := h.user.VerifyAdminPassword(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) GetAuthCoordinator(ctx *gin.Context) {
	if getAdminType(ctx) != service.AdminTypeCoordinator {
		h.forbidden(ctx, errors.New("this user is not coordinator"))
		return
	}

	in := &user.GetCoordinatorInput{
		CoordinatorID: getAdminID(ctx),
	}
	coordinator, err := h.user.GetCoordinator(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	productTypes, err := h.multiGetProductTypes(ctx, coordinator.ProductTypeIDs)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.CoordinatorResponse{
		Coordinator:  service.NewCoordinator(coordinator).Response(),
		ProductTypes: productTypes.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateAuthCoordinator(ctx *gin.Context) {
	if getAdminType(ctx) != service.AdminTypeCoordinator {
		h.forbidden(ctx, errors.New("this user is not coordinator"))
		return
	}

	req := &request.UpdateCoordinatorRequest{}
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

	in := &user.UpdateCoordinatorInput{
		CoordinatorID:     getAdminID(ctx),
		Lastname:          req.Lastname,
		Firstname:         req.Firstname,
		LastnameKana:      req.LastnameKana,
		FirstnameKana:     req.FirstnameKana,
		MarcheName:        req.MarcheName,
		Username:          req.Username,
		Profile:           req.Profile,
		ProductTypeIDs:    req.ProductTypeIDs,
		ThumbnailURL:      req.ThumbnailURL,
		HeaderURL:         req.HeaderURL,
		PromotionVideoURL: req.PromotionVideoURL,
		BonusVideoURL:     req.BonusVideoURL,
		InstagramID:       req.InstagramID,
		FacebookID:        req.FacebookID,
		PhoneNumber:       req.PhoneNumber,
		PostalCode:        req.PostalCode,
		PrefectureCode:    req.PrefectureCode,
		City:              req.City,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
		BusinessDays:      req.BusinessDays,
	}
	if err := h.user.UpdateCoordinator(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) GetAuthShipping(ctx *gin.Context) {
	if getAdminType(ctx) != service.AdminTypeCoordinator {
		h.forbidden(ctx, errors.New("this user is not coordinator"))
		return
	}

	in := &store.GetShippingByCoordinatorIDInput{
		CoordinatorID: getAdminID(ctx),
	}
	shipping, err := h.store.GetShippingByCoordinatorID(ctx, in)
	if errors.Is(err, exception.ErrNotFound) {
		// 配送設定の登録をしていない場合、デフォルト設定を返却する
		in := &store.GetDefaultShippingInput{}
		shipping, err = h.store.GetDefaultShipping(ctx, in)
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.ShippingResponse{
		Shipping: service.NewShipping(shipping).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpsertAuthShipping(ctx *gin.Context) {
	if getAdminType(ctx) != service.AdminTypeCoordinator {
		h.forbidden(ctx, errors.New("this user is not coordinator"))
		return
	}

	req := &request.UpsertShippingRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.UpsertShippingInput{
		CoordinatorID:     getAdminID(ctx),
		Box60Rates:        h.newShippingRatesForUpsert(req.Box60Rates),
		Box60Frozen:       req.Box60Frozen,
		Box80Rates:        h.newShippingRatesForUpsert(req.Box80Rates),
		Box80Frozen:       req.Box80Frozen,
		Box100Rates:       h.newShippingRatesForUpsert(req.Box100Rates),
		Box100Frozen:      req.Box100Frozen,
		HasFreeShipping:   req.HasFreeShipping,
		FreeShippingRates: req.FreeShippingRates,
	}
	if err := h.store.UpsertShipping(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
