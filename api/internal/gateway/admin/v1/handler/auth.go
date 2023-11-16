package handler

import (
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

func (h *handler) authRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/auth")

	r.GET("", h.GetAuth)
	r.POST("", h.SignIn)
	r.DELETE("", h.SignOut)
	r.POST("/refresh-token", h.RefreshAuthToken)
	r.POST("/device", h.authentication, h.RegisterDevice)
	r.PATCH("/email", h.authentication, h.UpdateAuthEmail)
	r.POST("/email/verified", h.VerifyAuthEmail)
	r.PATCH("/password", h.authentication, h.UpdateAuthPassword)
	r.POST("/forgot-password", h.ForgotAuthPassword)
	r.POST("/forgot-password/verified", h.ResetAuthPassword)
	r.GET("/user", h.authentication, h.GetAuthUser)
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
	switch getRole(ctx) {
	case service.AdminRoleAdministrator:
		auth, err = h.getAdministrator(ctx, adminID)
	case service.AdminRoleCoordinator:
		auth, err = h.getCoordinator(ctx, adminID)
	case service.AdminRoleProducer:
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

	ctx.JSON(http.StatusNoContent, gin.H{})
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

	ctx.JSON(http.StatusNoContent, gin.H{})
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

	ctx.JSON(http.StatusNoContent, gin.H{})
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

	ctx.JSON(http.StatusNoContent, gin.H{})
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

	ctx.JSON(http.StatusNoContent, gin.H{})
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

	ctx.JSON(http.StatusNoContent, gin.H{})
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

	ctx.JSON(http.StatusNoContent, gin.H{})
}
