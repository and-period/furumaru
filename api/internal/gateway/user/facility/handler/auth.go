package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
)

// @tag.name        Auth
// @tag.description 認証関連
func (h *handler) authRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/auth")

	r.POST("", h.SignIn)
	r.POST("/refresh-token", h.authentication, h.GetAccessToken)
	r.DELETE("", h.authentication, h.SignOut)
}

// @Summary     サインイン
// @Description LINEの認証トークンを渡すことで、ふるマルへサインインします。
// @Tags        Auth
// @Router      /facilities/{facilityId}/auth [post]
// @Param       facilityId path string true "施設ID"
// @Accept      json
// @Param       request body types.SignInRequest true "サインイン"
// @Produce     json
// @Success     200 {object} types.AuthResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     403 {object} util.ErrorResponse "退会済み"
// @Failure     404 {object} util.ErrorResponse "ユーザーが存在しない"
func (h *handler) SignIn(ctx *gin.Context) {
	req := &types.SignInRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	token, err := h.liffVerifier.VerifyIDToken(ctx, req.AuthToken, "")
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	in := &user.GetFacilityUserInput{
		ProducerID:   h.getProducerID(ctx),
		ProviderType: entity.UserAuthProviderTypeLINE,
		ProviderID:   token.Subject,
	}
	user, err := h.user.GetFacilityUser(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if user.Status == entity.UserStatusDeactivated {
		h.forbidden(ctx, err)
		return
	}
	auth, err := h.jwtGenerator.Generate(ctx, user.ID, h.getProducerID(ctx))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.AuthResponse{
		Auth: service.NewAuth(user.ID, auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     アクセストークンの再発行
// @Description 更新トークンを渡すことで、アクセストークンを再発行します。
// @Tags        Auth
// @Router      /facilities/{facilityId}/auth/refresh-token [post]
// @Param       facilityId path string true "施設ID"
// @Security    bearerauth
// @Accept      json
// @Param       request body types.GetAccessTokenRequest true "アクセストークンの再発行"
// @Produce     json
// @Success     200 {object} types.AuthResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     403 {object} util.ErrorResponse "リフレッシュトークンが無効"
func (h *handler) GetAccessToken(ctx *gin.Context) {
	req := &types.GetAccessTokenRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	token, err := h.jwtVerifier.VerifyRefreshToken(ctx, req.RefreshToken) // 更新トークンの検証
	if err != nil {
		h.forbidden(ctx, err)
		return
	}
	auth, err := h.jwtGenerator.RefreshAccessToken(ctx, req.RefreshToken) // アクセストークンの再発行
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.AuthResponse{
		Auth: service.NewAuth(token.UserID, auth).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     サインアウト
// @Description ふるマルからサインアウトします。
// @Tags        Auth
// @Router      /facilities/{facilityId}/auth [delete]
// @Param       facilityId path string true "施設ID"
// @Security    bearerauth
// @Produce     json
// @Success     204
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) SignOut(ctx *gin.Context) {
	if err := h.jwtGenerator.DeleteRefreshToken(ctx, h.getUserID(ctx)); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
