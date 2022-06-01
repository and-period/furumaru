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

func (h *apiV1Handler) authRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.GetAuth)
	rg.POST("", h.SignIn)
	rg.DELETE("", h.SignOut)
	rg.POST("/refresh-token", h.RefreshAuthToken)
	rg.PATCH("/initialized", h.InitializeAuth)
	rg.PATCH("/email", h.authentication(), h.UpdateAuthEmail)
	rg.POST("/email/verified", h.authentication(), h.VerifyAuthEmail)
	rg.PATCH("/password", h.authentication(), h.UpdateAuthPassword)
	rg.POST("/forgot-password", h.ForgotAuthPassword)
	rg.POST("/forgot-password/verified", h.ResetAuthPassword)
	rg.GET("/user", h.authentication(), h.GetAuthUser)
	rg.POST("/user", h.CreateAuth)
	rg.DELETE("/user", h.authentication(), h.DeleteAuth)
	rg.POST("/user/oauth", h.CreateAuthWithOAuth)
	rg.POST("/user/verified", h.VerifyAuth)
}

func (h *apiV1Handler) GetAuth(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err)
		return
	}

	in := &user.GetUserAuthInput{
		AccessToken: token,
	}
	auth, err := h.user.GetUserAuth(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) SignIn(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	req := &request.SignInRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.SignInUserInput{
		Key:      req.Username,
		Password: req.Password,
	}
	auth, err := h.user.SignInUser(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) SignOut(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err)
	}

	in := &user.SignOutUserInput{
		AccessToken: token,
	}
	if err := h.user.SignOutUser(c, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) RefreshAuthToken(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	req := &request.RefreshAuthTokenRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.RefreshUserTokenInput{
		RefreshToken: req.RefreshToken,
	}
	auth, err := h.user.RefreshUserToken(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) InitializeAuth(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	req := &request.InitializeAuthRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.InitializeUserInput{
		UserID:    getUserID(ctx),
		AccountID: req.AccountID,
		Username:  req.Username,
	}

	if err := h.user.InitializeUser(c, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) UpdateAuthEmail(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err)
		return
	}
	req := &request.UpdateAuthEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.UpdateUserEmailInput{
		AccessToken: token,
		Email:       req.Email,
	}
	if err := h.user.UpdateUserEmail(c, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) VerifyAuthEmail(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err)
		return
	}
	req := &request.VerifyAuthEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.VerifyUserEmailInput{
		AccessToken: token,
		VerifyCode:  req.VerifyCode,
	}
	if err := h.user.VerifyUserEmail(c, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) UpdateAuthPassword(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err)
		return
	}
	req := &request.UpdateAuthPasswordRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.UpdateUserPasswordInput{
		AccessToken:          token,
		OldPassword:          req.OldPassword,
		NewPassword:          req.NewPassword,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	if err := h.user.UpdateUserPassword(c, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) ForgotAuthPassword(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	req := &request.ForgotAuthPasswordRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.ForgotUserPasswordInput{
		Email: req.Email,
	}
	if err := h.user.ForgotUserPassword(c, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) ResetAuthPassword(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	req := &request.ResetAuthPasswordRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.VerifyUserPasswordInput{
		Email:                req.Email,
		VerifyCode:           req.VerifyCode,
		NewPassword:          req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	if err := h.user.VerifyUserPassword(c, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) GetAuthUser(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	in := &user.GetUserInput{
		UserID: getUserID(ctx),
	}
	u, err := h.user.GetUser(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.AuthUserResponse{
		AuthUser: service.NewAuthUser(u).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) CreateAuth(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	req := &request.CreateAuthRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.CreateUserInput{
		Email:                req.Email,
		PhoneNumber:          req.PhoneNumber,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	userID, err := h.user.CreateUser(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.CreateAuthResponse{
		ID: userID,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) DeleteAuth(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	in := &user.DeleteUserInput{
		UserID: getUserID(ctx),
	}
	if err := h.user.DeleteUser(c, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) CreateAuthWithOAuth(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err)
		return
	}

	in := &user.CreateUserWithOAuthInput{
		AccessToken: token,
	}
	u, err := h.user.CreateUserWithOAuth(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.AuthUserResponse{
		AuthUser: service.NewAuthUser(u).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) VerifyAuth(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	req := &request.VerifyAuthRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.VerifyUserInput{
		UserID:     req.ID,
		VerifyCode: req.VerifyCode,
	}
	if err := h.user.VerifyUser(c, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
