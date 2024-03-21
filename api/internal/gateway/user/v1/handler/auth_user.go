package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
)

func (h *handler) authUserRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/users/me")

	r.POST("", h.CreateAuthUser)
	r.POST("/oauth", h.CreateAuthUserWithOAuth)
	r.POST("/verified", h.VerifyAuthUser)

	auth := r.Group("", h.authentication)
	auth.GET("", h.GetAuthUser)
	auth.DELETE("", h.DeleteAuthUser)
	auth.PATCH("/email", h.UpdateAuthUserEmail)
	auth.POST("/email/verified", h.VerifyAuthUserEmail)
	auth.PATCH("/username", h.UpdateAuthUserUsername)
	auth.PATCH("/account-id", h.UpdateAuthUserAccountID)
	auth.PATCH("/thumbnail", h.UpdateAuthUserThumbnail)
}

func (h *handler) GetAuthUser(ctx *gin.Context) {
	in := &user.GetUserInput{
		UserID: h.getUserID(ctx),
	}
	user, err := h.user.GetUser(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.AuthUserResponse{
		AuthUser: service.NewAuthUser(user).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateAuthUser(ctx *gin.Context) {
	req := &request.CreateAuthUserRequest{}
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
	res := &response.CreateAuthUserResponse{
		ID: userID,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) VerifyAuthUser(ctx *gin.Context) {
	req := &request.VerifyAuthUserRequest{}
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

func (h *handler) CreateAuthUserWithOAuth(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	req := &request.CreateAuthUserWithOAuthRequest{}
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
	user, err := h.user.CreateMemberWithOAuth(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.AuthUserResponse{
		AuthUser: service.NewAuthUser(user).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateAuthUserEmail(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	req := &request.UpdateAuthUserEmailRequest{}
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

func (h *handler) VerifyAuthUserEmail(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	req := &request.VerifyAuthUserEmailRequest{}
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

func (h *handler) UpdateAuthUserUsername(ctx *gin.Context) {
	req := &request.UpdateAuthUserUsernameRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.UpdateMemberUsernameInput{
		UserID:   h.getUserID(ctx),
		Username: req.Username,
	}
	if err := h.user.UpdateMemberUsername(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) UpdateAuthUserAccountID(ctx *gin.Context) {
	req := &request.UpdateAuthUserAccountIDRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.UpdateMemberAccountIDInput{
		UserID:    h.getUserID(ctx),
		AccountID: req.AccountID,
	}
	if err := h.user.UpdateMemberAccountID(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) UpdateAuthUserThumbnail(ctx *gin.Context) {
	req := &request.UpdateAuthUserThumbnailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.UpdateMemberThumbnailURLInput{
		UserID:       h.getUserID(ctx),
		ThumbnailURL: req.ThumbnailURL,
	}
	if err := h.user.UpdateMemberThumbnailURL(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteAuthUser(ctx *gin.Context) {
	in := &user.DeleteUserInput{
		UserID: h.getUserID(ctx),
	}
	if err := h.user.DeleteUser(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) multiGetUsers(ctx context.Context, userIDs []string) (entity.Users, error) {
	if len(userIDs) == 0 {
		return entity.Users{}, nil
	}
	in := &user.MultiGetUsersInput{
		UserIDs: userIDs,
	}
	return h.user.MultiGetUsers(ctx, in)
}
