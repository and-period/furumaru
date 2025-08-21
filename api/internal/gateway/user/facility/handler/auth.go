package handler

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/gin-gonic/gin"
)

// @tag.name        Auth
// @tag.description 認証関連
func (h *handler) authRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/auth")

	r.POST("", h.SignIn)
	r.DELETE("", h.SignOut)
}

// @Summary     サインイン
// @Description LINEの認証トークンを渡すことで、ふるマルへサインインします。
// @Tags        Auth
// @Router      /facilities/{facilityId}/auth [post]
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
	// TODO: 詳細の実装
	res := &response.AuthResponse{
		Auth: &response.Auth{},
	}
	ctx.JSON(200, res)
}

// @Summary     サインアウト
// @Description ふるマルからサインアウトします。
// @Tags        Auth
// @Router      /facilities/{facilityId}/auth [delete]
// @Security    bearerauth
// @Produce     json
// @Success     204
func (h *handler) SignOut(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.Status(204)
}
