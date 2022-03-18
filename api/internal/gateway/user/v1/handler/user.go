package handler

import (
	"net/http"

	"github.com/and-period/marche/api/internal/gateway/user/v1/response"
	"github.com/and-period/marche/api/internal/gateway/util"
	"github.com/and-period/marche/api/proto/user"
	"github.com/gin-gonic/gin"
)

func (h *apiV1Handler) GetUserMe(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	in := &user.GetUserRequest{
		UserId: getUserID(ctx),
	}
	out, err := h.user.GetUser(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.UserMeResponse{
		ID:          out.User.Id,
		Email:       out.User.Email,
		PhoneNumber: out.User.PhoneNumber,
	}
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
	// TODO: 詳細の実装
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) ForgotUserPassword(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) ResetUserPassword(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *apiV1Handler) DeleteUser(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.JSON(http.StatusNoContent, gin.H{})
}
