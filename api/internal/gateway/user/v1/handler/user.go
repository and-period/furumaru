package handler

import (
	"net/http"

	"github.com/and-period/marche/api/internal/gateway/user/v1/request"
	"github.com/and-period/marche/api/internal/gateway/user/v1/response"
	"github.com/and-period/marche/api/internal/gateway/util"
	"github.com/and-period/marche/api/internal/user"
	"github.com/gin-gonic/gin"
)

func (h *apiV1Handler) userRoutes(rg *gin.RouterGroup) {
	rg.POST("", h.CreateUser)
	rg.POST("/oauth", h.CreateUserWithOAuth)
	rg.POST("/verified", h.VerifyUser)
	rg.GET("/me", h.authentication(), h.GetUserMe)
	rg.DELETE("/me", h.authentication(), h.DeleteUser)
	rg.PATCH("/me/email", h.authentication(), h.UpdateUserEmail)
	rg.POST("/me/email/verified", h.VerifyUserEmail)
	rg.PATCH("/me/password", h.authentication(), h.UpdateUserPassword)
	rg.POST("/me/forgot-password", h.ForgotUserPassword)
	rg.POST("/me/forgot-password/verified", h.ResetUserPassword)
}

func (h *apiV1Handler) GetUserMe(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	in := &user.GetUserInput{
		UserID: getUserID(ctx),
	}
	u, err := h.user.GetUser(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.UserMeResponse{
		ID:           u.ID,
		Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
		ThumbnailURL: u.ThumbnailURL,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) CreateUser(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	req := &request.CreateUserRequest{}
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

	res := &response.CreateUserResponse{
		ID: userID,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) VerifyUser(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	req := &request.VerifyUserRequest{}
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

func (h *apiV1Handler) CreateUserWithOAuth(ctx *gin.Context) {
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

	res := &response.UserMeResponse{
		ID:          u.ID,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) UpdateUserEmail(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err)
		return
	}
	req := &request.UpdateUserEmailRequest{}
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

func (h *apiV1Handler) VerifyUserEmail(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err)
		return
	}
	req := &request.VerifyUserEmailRequest{}
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

func (h *apiV1Handler) InitializeUser(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	req := &request.InitializeUserRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	in := &user.InitializeUserInput{
		UserID:    req.ID,
		AccountID: req.AccountID,
		Username:  req.Username,
	}

	if err := h.user.InitializeUser(c, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) UpdateUserPassword(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err)
		return
	}
	req := &request.UpdateUserPasswordRequest{}
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

func (h *apiV1Handler) ForgotUserPassword(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	req := &request.ForgotUserPasswordRequest{}
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

func (h *apiV1Handler) ResetUserPassword(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	req := &request.ResetUserPasswordRequest{}
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

func (h *apiV1Handler) DeleteUser(ctx *gin.Context) {
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
