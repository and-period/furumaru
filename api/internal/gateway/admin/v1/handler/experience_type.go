package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
)

// @tag.name        ExperienceType
// @tag.description 体験タイプ関連
func (h *handler) experienceTypeRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/experience-types", h.authentication)

	r.GET("", h.ListExperienceTypes)
	r.POST("", h.CreateExperienceType)
	r.PATCH("/:experienceTypeId", h.UpdateExperienceType)
	r.DELETE("/:experienceTypeId", h.DeleteExperienceType)
}

// @Summary     体験タイプ一覧取得
// @Description 体験タイプの一覧を取得します。名前でのフィルタリングが可能です。
// @Tags        ExperienceType
// @Router      /v1/experience-types [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Param       name query string false "体験タイプ名(あいまい検索)" example("農業")
// @Produce     json
// @Success     200 {object} types.ExperienceTypesResponse
func (h *handler) ListExperienceTypes(ctx *gin.Context) {
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

	in := &store.ListExperienceTypesInput{
		Name:   util.GetQuery(ctx, "name", ""),
		Limit:  limit,
		Offset: offset,
	}
	etypes, total, err := h.store.ListExperienceTypes(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.ExperienceTypesResponse{
		ExperienceTypes: service.NewExperienceTypes(etypes).Response(),
		Total:           total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     体験タイプ登録
// @Description 新しい体験タイプを登録します。
// @Tags        ExperienceType
// @Router      /v1/experience-types [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.CreateExperienceTypeRequest true "体験タイプ情報"
// @Produce     json
// @Success     200 {object} types.ExperienceTypeResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     409 {object} util.ErrorResponse "すでに存在する体験タイプ名"
func (h *handler) CreateExperienceType(ctx *gin.Context) {
	req := &types.CreateExperienceTypeRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.CreateExperienceTypeInput{
		Name: req.Name,
	}
	experienceType, err := h.store.CreateExperienceType(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.ExperienceTypeResponse{
		ExperienceType: service.NewExperienceType(experienceType).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     体験タイプ更新
// @Description 体験タイプの情報を更新します。
// @Tags        ExperienceType
// @Router      /v1/experience-types/{experienceTypeId} [patch]
// @Security    bearerauth
// @Param       experienceTypeId path string true "体験タイプID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.UpdateExperienceTypeRequest true "体験タイプ情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     404 {object} util.ErrorResponse "体験タイプが存在しない"
// @Failure     409 {object} util.ErrorResponse "すでに存在する体験タイプ名"
func (h *handler) UpdateExperienceType(ctx *gin.Context) {
	req := &types.UpdateExperienceTypeRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.UpdateExperienceTypeInput{
		ExperienceTypeID: ctx.Param("experienceTypeId"),
		Name:             req.Name,
	}
	if err := h.store.UpdateExperienceType(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     体験タイプ削除
// @Description 体験タイプを削除します。
// @Tags        ExperienceType
// @Router      /v1/experience-types/{experienceTypeId} [delete]
// @Security    bearerauth
// @Param       experienceTypeId path string true "体験タイプID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "体験タイプが存在しない"
// @Failure     412 {object} util.ErrorResponse "体験側で紐づいているため削除不可"
func (h *handler) DeleteExperienceType(ctx *gin.Context) {
	in := &store.DeleteExperienceTypeInput{
		ExperienceTypeID: ctx.Param("experienceTypeId"),
	}
	if err := h.store.DeleteExperienceType(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) multiGetExperienceTypes(ctx context.Context, experienceTypeIDs []string) (service.ExperienceTypes, error) {
	if len(experienceTypeIDs) == 0 {
		return service.ExperienceTypes{}, nil
	}
	in := &store.MultiGetExperienceTypesInput{
		ExperienceTypeIDs: experienceTypeIDs,
	}
	experienceTypes, err := h.store.MultiGetExperienceTypes(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewExperienceTypes(experienceTypes), nil
}

func (h *handler) getExperienceType(ctx context.Context, experienceTypeID string) (*service.ExperienceType, error) {
	in := &store.GetExperienceTypeInput{
		ExperienceTypeID: experienceTypeID,
	}
	experienceType, err := h.store.GetExperienceType(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewExperienceType(experienceType), nil
}
