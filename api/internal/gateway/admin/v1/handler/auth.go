package handler

import (
	"net/http"

	"github.com/and-period/marche/api/internal/gateway/admin/v1/request"
	"github.com/and-period/marche/api/internal/gateway/admin/v1/response"
	"github.com/and-period/marche/api/internal/gateway/admin/v1/service"
	"github.com/and-period/marche/api/internal/gateway/util"
	"github.com/and-period/marche/api/internal/user"
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

	in := &user.GetAdminAuthInput{
		AccessToken: token,
	}
	auth, err := h.user.GetAdminAuth(c, in)
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

	in := &user.SignInAdminInput{
		Key:      req.Username,
		Password: req.Password,
	}
	auth, err := h.user.SignInAdmin(c, in)
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

	in := &user.SignOutAdminInput{
		AccessToken: token,
	}
	if err := h.user.SignOutAdmin(c, in); err != nil {
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

	in := &user.RefreshAdminTokenInput{
		RefreshToken: req.RefreshToken,
	}
	auth, err := h.user.RefreshAdminToken(c, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.AuthResponse{
		Auth: service.NewAuth(auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
