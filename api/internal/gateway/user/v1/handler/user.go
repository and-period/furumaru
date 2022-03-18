package handler

import (
	"net/http"

	"github.com/and-period/marche/api/internal/gateway/user/v1/request"
	"github.com/and-period/marche/api/internal/gateway/user/v1/response"
	"github.com/and-period/marche/api/internal/gateway/util"
	"github.com/and-period/marche/api/proto/user"
	"github.com/gin-gonic/gin"
)

func (h *apiV1Handler) GetUserMe(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.UserMeResponse{}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) CreateUser(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.CreateUserResponse{}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) VerifyUser(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) CreateUserWithOAuth(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.UserMeResponse{}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) UpdateUserEmail(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) VerifyUserEmail(ctx *gin.Context) {
	// TODO: 詳細の実装
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

	in := &user.UpdateUserPasswordRequest{
		AccessToken:          token,
		OldPassword:          req.OldPassword,
		NewPassword:          req.NewPassword,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	_, err = h.user.UpdateUserPassword(c, in)
	if err != nil {
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

	in := &user.ForgotUserPasswordRequest{
		Email: req.Email,
	}
	_, err := h.user.ForgotUserPassword(c, in)
	if err != nil {
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

	in := &user.VerifyUserPasswordRequest{
		Email:                req.Email,
		VerifyCode:           req.VerifyCode,
		NewPassword:          req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}
	_, err := h.user.VerifyUserPassword(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) DeleteUser(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.JSON(http.StatusNoContent, gin.H{})
}
