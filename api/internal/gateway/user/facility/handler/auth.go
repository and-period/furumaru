package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
)

// @tag.name        Auth
// @tag.description 認証関連
func (h *handler) authRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/auth")

	r.POST("", h.SignIn)
	r.DELETE("", h.authentication, h.SignOut)
}

// @Summary     サインイン
// @Description LINEの認証トークンを渡すことで、ふるマルへサインインします。
// @Tags        Auth
// @Router      /facilities/{facilityId}/auth [post]
// @Param       facilityId path string true "施設ID"
// @Accept      json
// @Param       request body request.SignInRequest true "サインイン"
// @Produce     json
// @Success     200 {object} response.AuthResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     403 {object} util.ErrorResponse "退会済み"
// @Failure     404 {object} util.ErrorResponse "ユーザーが存在しない"
func (h *handler) SignIn(ctx *gin.Context) {
	req := &request.SignInRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	token, err := h.lineVerifier.VerifyIDToken(ctx, req.AuthToken, "")
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
	res := &response.AuthResponse{
		Auth: service.NewAuth(user, auth).Response(),
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
func (h *handler) SignOut(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.Status(http.StatusNoContent)
}
