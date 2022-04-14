package handler

import (
	"net/http"

	"github.com/and-period/marche/api/internal/gateway/entity"
	"github.com/and-period/marche/api/internal/gateway/user/v1/request"
	"github.com/and-period/marche/api/internal/gateway/user/v1/response"
	"github.com/and-period/marche/api/internal/gateway/user/v1/service"
	"github.com/and-period/marche/api/internal/gateway/util"
	"github.com/and-period/marche/api/proto/user"
	"github.com/gin-gonic/gin"
)

func (h *apiV1Handler) authRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.GetAuth)
	rg.POST("", h.SignIn)
	rg.DELETE("", h.SignOut)
	rg.POST("/refresh-token", h.RefreshAuthToken)
}

func (h *apiV1Handler) GetAuth(ctx *gin.Context) {
	c := util.SetMetadata(ctx)

	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err)
		return
	}

	in := &user.GetUserAuthRequest{
		AccessToken: token,
	}
	out, err := h.user.GetUserAuth(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	auth := entity.NewUserAuth(out.Auth)

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

	in := &user.SignInUserRequest{
		Username: req.Username,
		Password: req.Password,
	}
	out, err := h.user.SignInUser(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	auth := entity.NewUserAuth(out.Auth)

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

	in := &user.SignOutUserRequest{
		AccessToken: token,
	}
	_, err = h.user.SignOutUser(c, in)
	if err != nil {
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

	in := &user.RefreshUserTokenRequest{
		RefreshToken: req.RefreshToken,
	}
	out, err := h.user.RefreshUserToken(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	auth := entity.NewUserAuth(out.Auth)

	res := &response.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
