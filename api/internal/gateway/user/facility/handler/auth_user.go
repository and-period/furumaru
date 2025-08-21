package handler

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/gin-gonic/gin"
)

// @tag.name        AuthUser
// @tag.description 認証済みユーザー関連
func (h *handler) authUserRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/users/me", h.authentication)

	r.GET("", h.GetAuthUser)
	r.POST("", h.CreateAuthUser)
	r.PUT("/check-in", h.UpdateAuthUserCheckIn)
}

// @Summary     ユーザー情報取得
// @Description ユーザーの詳細情報を取得します。
// @Tags        AuthUser
// @Router      /facilities/{facilityId}/users/me [get]
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} response.AuthUserResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) GetAuthUser(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.AuthUserResponse{
		AuthUser: &response.AuthUser{},
	}
	ctx.JSON(200, res)
}

// @Summary     ユーザー情報登録
// @Description ユーザーの詳細情報を登録します。
// @Tags        AuthUser
// @Router      /facilities/{facilityId}/users/me [post]
// @Accept      json
// @Param				request body request.CreateAuthUserRequest true "ユーザー情報"
// @Produce     json
// @Success     200 {object} response.AuthUserResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     409 {object} util.ErrorResponse "ユーザーが既に存在する"
func (h *handler) CreateAuthUser(ctx *gin.Context) {
	req := &request.CreateAuthUserRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	res := &response.AuthUserResponse{
		AuthUser: &response.AuthUser{},
	}
	ctx.JSON(200, res)
}

// @Summary     ユーザー情報更新
// @Description ユーザーの最新のチェックイン日時を更新します。
// @Tags        AuthUser
// @Router      /facilities/{facilityId}/users/me/check-in [put]
// @Security    bearerauth
// @Accept      json
// @Param       request body request.UpdateAuthUserCheckInRequest true "最新のチェックイン情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     412 {object} util.ErrorResponse "更新日時が不正"
func (h *handler) UpdateAuthUserCheckIn(ctx *gin.Context) {
	req := &request.UpdateAuthUserCheckInRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	ctx.Status(204)
}
