package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

// @tag.name        Administrator
// @tag.description システム管理者関連
func (h *handler) administratorRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/administrators", h.authentication)

	r.GET("", h.ListAdministrators)
	r.POST("", h.CreateAdministrator)
	r.GET("/:adminId", h.GetAdministrator)
	r.PATCH("/:adminId", h.UpdateAdministrator)
	r.PATCH("/:adminId/email", h.UpdateAdministratorEmail)
	r.PATCH("/:adminId/password", h.ResetAdministratorPassword)
	r.DELETE("/:adminId", h.DeleteAdministrator)
}

// @Summary     システム管理者一覧取得
// @Description システム管理者の一覧を取得します。
// @Tags        Administrator
// @Router      /v1/administrators [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Produce     json
// @Success     200 {object} types.AdministratorsResponse
func (h *handler) ListAdministrators(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.ListAdministratorsInput{
		Limit:  limit,
		Offset: offset,
	}
	admins, total, err := h.user.ListAdministrators(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.AdministratorsResponse{
		Administrators: service.NewAdministrators(admins).Response(),
		Total:          total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     システム管理者取得
// @Description システム管理者の詳細情報を取得します。
// @Tags        Administrator
// @Router      /v1/administrators/{adminId} [get]
// @Security    bearerauth
// @Param       adminId path string true "管理者ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.AdministratorResponse
// @Failure     404 {object} util.ErrorResponse "システム管理者が存在しない"
func (h *handler) GetAdministrator(ctx *gin.Context) {
	in := &user.GetAdministratorInput{
		AdministratorID: util.GetParam(ctx, "adminId"),
	}
	admin, err := h.user.GetAdministrator(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.AdministratorResponse{
		Administrator: service.NewAdministrator(admin).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     システム管理者登録
// @Description 新しいシステム管理者を登録します。
// @Tags        Administrator
// @Router      /v1/administrators [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.CreateAdministratorRequest true "システム管理者情報"
// @Produce     json
// @Success     200 {object} types.AdministratorResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     409 {object} util.ErrorResponse "すでに存在するメールアドレス"
func (h *handler) CreateAdministrator(ctx *gin.Context) {
	req := &types.CreateAdministratorRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.CreateAdministratorInput{
		Lastname:      req.Lastname,
		Firstname:     req.Firstname,
		LastnameKana:  req.LastnameKana,
		FirstnameKana: req.FirstnameKana,
		Email:         req.Email,
		PhoneNumber:   req.PhoneNumber,
	}
	admin, err := h.user.CreateAdministrator(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.AdministratorResponse{
		Administrator: service.NewAdministrator(admin).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     システム管理者更新
// @Description システム管理者の情報を更新します。
// @Tags        Administrator
// @Router      /v1/administrators/{adminId} [patch]
// @Security    bearerauth
// @Param       adminId path string true "システム管理者ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.UpdateAdministratorRequest true "システム管理者情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     404 {object} util.ErrorResponse "システム管理者が存在しない"
func (h *handler) UpdateAdministrator(ctx *gin.Context) {
	req := &types.UpdateAdministratorRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.UpdateAdministratorInput{
		AdministratorID: util.GetParam(ctx, "adminId"),
		Lastname:        req.Lastname,
		Firstname:       req.Firstname,
		LastnameKana:    req.LastnameKana,
		FirstnameKana:   req.FirstnameKana,
		PhoneNumber:     req.PhoneNumber,
	}
	if err := h.user.UpdateAdministrator(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     システム管理者メールアドレス更新
// @Description システム管理者のメールアドレスを更新します。
// @Tags        Administrator
// @Router      /v1/administrators/{adminId}/email [patch]
// @Security    bearerauth
// @Param       adminId path string true "システム管理者ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.UpdateAdministratorEmailRequest true "メールアドレス"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     404 {object} util.ErrorResponse "存在しないシステム管理者"
// @Failure     409 {object} util.ErrorResponse "すでに存在するメールアドレス"
func (h *handler) UpdateAdministratorEmail(ctx *gin.Context) {
	req := &types.UpdateAdministratorEmailRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &user.UpdateAdministratorEmailInput{
		AdministratorID: util.GetParam(ctx, "adminId"),
		Email:           req.Email,
	}
	if err := h.user.UpdateAdministratorEmail(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     システム管理者パスワード更新(ランダム生成)
// @Description システム管理者のパスワードをランダムに生成して更新します。
// @Tags        Administrator
// @Router      /v1/administrators/{adminId}/password [patch]
// @Security    bearerauth
// @Param       adminId path string true "システム管理者ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "システム管理者が存在しない"
func (h *handler) ResetAdministratorPassword(ctx *gin.Context) {
	in := &user.ResetAdministratorPasswordInput{
		AdministratorID: util.GetParam(ctx, "adminId"),
	}
	if err := h.user.ResetAdministratorPassword(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     システム管理者退会
// @Description システム管理者を削除します。
// @Tags        Administrator
// @Router      /v1/administrators/{adminId} [delete]
// @Security    bearerauth
// @Param       adminId path string true "システム管理者ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     404 {object} util.ErrorResponse "システム管理者が存在しない"
func (h *handler) DeleteAdministrator(ctx *gin.Context) {
	in := &user.DeleteAdministratorInput{
		AdministratorID: util.GetParam(ctx, "adminId"),
	}
	if err := h.user.DeleteAdministrator(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) getAdministrator(ctx context.Context, administratorID string) (*service.Administrator, error) {
	in := &user.GetAdministratorInput{
		AdministratorID: administratorID,
	}
	administrator, err := h.user.GetAdministrator(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewAdministrator(administrator), nil
}
