package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
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
	r.PATCH("/password", h.authentication, h.UpdateAuthPassword)
	r.POST("/forgot-password", h.ForgotAuthPassword)
	r.POST("/forgot-password/verified", h.ResetAuthPassword)
	r.GET("/google", h.AuthGoogleAccount)
	r.GET("/line", h.AuthLINEAccount)
}

func (h *handler) GetAuth(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}

	in := &user.GetUserAuthInput{
		AccessToken: token,
	}
	auth, err := h.user.GetUserAuth(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) SignIn(ctx *gin.Context) {
	req := &request.SignInRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.SignInUserInput{
		Key:      req.Username,
		Password: req.Password,
	}
	auth, err := h.user.SignInUser(ctx, in)
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
	}

	in := &user.SignOutUserInput{
		AccessToken: token,
	}
	if err := h.user.SignOutUser(ctx, in); err != nil {
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

	in := &user.RefreshUserTokenInput{
		RefreshToken: req.RefreshToken,
	}
	auth, err := h.user.RefreshUserToken(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
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

	in := &user.UpdateMemberPasswordInput{
		AccessToken:          token,
		OldPassword:          req.OldPassword,
		NewPassword:          req.NewPassword,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	if err := h.user.UpdateMemberPassword(ctx, in); err != nil {
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

	in := &user.ForgotMemberPasswordInput{
		Email: req.Email,
	}
	if err := h.user.ForgotMemberPassword(ctx, in); err != nil {
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

	in := &user.VerifyMemberPasswordInput{
		Email:                req.Email,
		VerifyCode:           req.VerifyCode,
		NewPassword:          req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	if err := h.user.VerifyMemberPassword(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) AuthGoogleAccount(ctx *gin.Context) {
	in := &user.AuthMemberWithGoogleInput{
		AuthMemberDetailWithOAuth: user.AuthMemberDetailWithOAuth{
			SessionID:   h.getSessionID(ctx),
			State:       util.GetQuery(ctx, "state", ""),
			RedirectURI: util.GetQuery(ctx, "redirectUri", ""),
		},
	}
	authURL, err := h.user.AuthMemberWithGoogle(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.AuthGoogleAccountResponse{
		URL: authURL,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) AuthLINEAccount(ctx *gin.Context) {
	in := &user.AuthMemberWithLINEInput{
		AuthMemberDetailWithOAuth: user.AuthMemberDetailWithOAuth{
			SessionID:   h.getSessionID(ctx),
			State:       util.GetQuery(ctx, "state", ""),
			RedirectURI: util.GetQuery(ctx, "redirectUri", ""),
		},
	}
	authURL, err := h.user.AuthMemberWithLINE(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.AuthLINEAccountResponse{
		URL: authURL,
	}
	ctx.JSON(http.StatusOK, res)
}
