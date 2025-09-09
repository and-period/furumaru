package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
)

// @tag.name        AuthUser
// @tag.description 認証済みユーザー関連
func (h *handler) authUserRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/users")

	r.POST("", h.CreateAuthUser)
	r.GET("/me", h.authentication, h.GetAuthUser)
	r.PUT("/me", h.authentication, h.UpdateAuthUser)
}

// @Summary     ユーザー情報取得
// @Description ユーザーの詳細情報を取得します。
// @Tags        AuthUser
// @Router      /facilities/{facilityId}/users/me [get]
// @Param       facilityId path string true "施設ID"
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} types.AuthUserResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) GetAuthUser(ctx *gin.Context) {
	in := &user.GetUserInput{
		UserID: h.getUserID(ctx),
	}
	user, err := h.user.GetUser(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.AuthUserResponse{
		AuthUser: service.NewAuthUser(user).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     ユーザー情報登録
// @Description ユーザーの詳細情報を登録します。
// @Tags        AuthUser
// @Router      /facilities/{facilityId}/users [post]
// @Param       facilityId path string true "施設ID"
// @Accept      json
// @Param				request body types.CreateAuthUserRequest true "ユーザー情報"
// @Produce     json
// @Success     200 {object} types.AuthUserResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー。不正なトークン"
// @Failure     409 {object} util.ErrorResponse "ユーザーが既に存在する"
func (h *handler) CreateAuthUser(ctx *gin.Context) {
	req := &types.CreateAuthUserRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	token, err := h.liffVerifier.VerifyIDToken(ctx, req.AuthToken, "")
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	claims, err := h.liffVerifier.GetClaims(token)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	in := &user.CreateFacilityUserInput{
		ProducerID:    h.getProducerID(ctx),
		ProviderType:  entity.UserAuthProviderTypeLINE,
		ProviderID:    claims.Sub,
		Lastname:      req.Lastname,
		Firstname:     req.Firstname,
		LastnameKana:  req.LastnameKana,
		FirstnameKana: req.FirstnameKana,
		Email:         claims.Email,
		PhoneNumber:   req.PhoneNumber,
		LastCheckInAt: jst.ParseFromUnix(req.LastCheckInAt),
	}
	user, err := h.user.CreateFacilityUser(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.AuthUserResponse{
		AuthUser: service.NewAuthUser(user).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     ユーザー情報更新
// @Description ユーザーの詳細情報を更新します。
// @Tags        AuthUser
// @Router      /facilities/{facilityId}/users/me [put]
// @Param       facilityId path string true "施設ID"
// @Security    bearerauth
// @Accept      json
// @Param       request body types.UpdateAuthUserRequest true "ユーザー情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) UpdateAuthUser(ctx *gin.Context) {
	req := &types.UpdateAuthUserRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &user.UpdateFacilityUserInput{
		UserID:        h.getUserID(ctx),
		Lastname:      req.Lastname,
		Firstname:     req.Firstname,
		LastnameKana:  req.LastnameKana,
		FirstnameKana: req.FirstnameKana,
		PhoneNumber:   req.PhoneNumber,
		LastCheckInAt: jst.ParseFromUnix(req.LastCheckInAt),
	}
	if err := h.user.UpdateFacilityUser(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
