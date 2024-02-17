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
	r.PATCH("/email", h.authentication, h.UpdateAuthEmail)
	r.POST("/email/verified", h.authentication, h.VerifyAuthEmail)
	r.PATCH("/password", h.authentication, h.UpdateAuthPassword)
	r.POST("/forgot-password", h.ForgotAuthPassword)
	r.POST("/forgot-password/verified", h.ResetAuthPassword)
	r.GET("/user", h.authentication, h.GetAuthUser)
	r.POST("/user", h.CreateAuth)
	r.DELETE("/user", h.authentication, h.DeleteAuth)
	r.POST("/user/oauth", h.CreateAuthWithOAuth)
	r.POST("/user/verified", h.VerifyAuth)
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

	in := &user.UpdateMemberEmailInput{
		AccessToken: token,
		Email:       req.Email,
	}
	if err := h.user.UpdateMemberEmail(ctx, in); err != nil {
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

	in := &user.VerifyMemberEmailInput{
		AccessToken: token,
		VerifyCode:  req.VerifyCode,
	}
	if err := h.user.VerifyMemberEmail(ctx, in); err != nil {
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

func (h *handler) GetAuthUser(ctx *gin.Context) {
	in := &user.GetUserInput{
		UserID: h.getUserID(ctx),
	}
	u, err := h.user.GetUser(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.AuthUserResponse{
		AuthUser: service.NewAuthUser(u).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateAuth(ctx *gin.Context) {
	req := &request.CreateAuthRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.CreateMemberInput{
		Username:             req.Username,
		AccountID:            req.AccountID,
		Lastname:             req.Lastname,
		Firstname:            req.Firstname,
		LastnameKana:         req.LastnameKana,
		FirstnameKana:        req.FirstnameKana,
		Email:                req.Email,
		PhoneNumber:          req.PhoneNumber,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	userID, err := h.user.CreateMember(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.CreateAuthResponse{
		ID: userID,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) DeleteAuth(ctx *gin.Context) {
	in := &user.DeleteUserInput{
		UserID: h.getUserID(ctx),
	}
	if err := h.user.DeleteUser(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) CreateAuthWithOAuth(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	req := &request.CreateAuthWithOAuthRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.CreateMemberWithOAuthInput{
		AccessToken:   token,
		Username:      req.Username,
		AccountID:     req.AccountID,
		Lastname:      req.Lastname,
		Firstname:     req.Firstname,
		LastnameKana:  req.LastnameKana,
		FirstnameKana: req.FirstnameKana,
		PhoneNumber:   req.PhoneNumber,
	}
	u, err := h.user.CreateMemberWithOAuth(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.AuthUserResponse{
		AuthUser: service.NewAuthUser(u).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) VerifyAuth(ctx *gin.Context) {
	req := &request.VerifyAuthRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.VerifyMemberInput{
		UserID:     req.ID,
		VerifyCode: req.VerifyCode,
	}
	if err := h.user.VerifyMember(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
