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

// @tag.name        SpotType
// @tag.description スポットタイプ関連
func (h *handler) spotTypeRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/spot-types", h.authentication)

	r.GET("", h.ListSpotTypes)
	r.POST("", h.CreateSpotType)
	r.PATCH("/:spotTypeId", h.UpdateSpotType)
	r.DELETE("/:spotTypeId", h.DeleteSpotType)
}

// @Summary     スポットタイプ一覧取得
// @Description スポットタイプの一覧を取得します。名前でのフィルタリングが可能です。
// @Tags        SpotType
// @Router      /v1/spot-types [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Param       name query string false "スポットタイプ名(あいまい検索)" example("観光地")
// @Produce     json
// @Success     200 {object} types.SpotTypesResponse
func (h *handler) ListSpotTypes(ctx *gin.Context) {
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

	in := &store.ListSpotTypesInput{
		Name:   util.GetQuery(ctx, "name", ""),
		Limit:  limit,
		Offset: offset,
	}
	stypes, total, err := h.store.ListSpotTypes(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.SpotTypesResponse{
		SpotTypes: service.NewSpotTypes(stypes).Response(),
		Total:     total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     スポットタイプ登録
// @Description 新しいスポットタイプを登録します。
// @Tags        SpotType
// @Router      /v1/spot-types [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.CreateSpotTypeRequest true "スポットタイプ情報"
// @Produce     json
// @Success     200 {object} types.SpotTypeResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     409 {object} util.ErrorResponse "すでに存在するスポットタイプ名"
func (h *handler) CreateSpotType(ctx *gin.Context) {
	req := &types.CreateSpotTypeRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.CreateSpotTypeInput{
		Name: req.Name,
	}
	spotType, err := h.store.CreateSpotType(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.SpotTypeResponse{
		SpotType: service.NewSpotType(spotType).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     スポットタイプ更新
// @Description スポットタイプの情報を更新します。
// @Tags        SpotType
// @Router      /v1/spot-types/{spotTypeId} [patch]
// @Security    bearerauth
// @Param       spotTypeId path string true "スポットタイプID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.UpdateSpotTypeRequest true "スポットタイプ情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     404 {object} util.ErrorResponse "スポットタイプが存在しない"
// @Failure     409 {object} util.ErrorResponse "すでに存在するスポットタイプ名"
func (h *handler) UpdateSpotType(ctx *gin.Context) {
	req := &types.UpdateSpotTypeRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.UpdateSpotTypeInput{
		SpotTypeID: ctx.Param("spotTypeId"),
		Name:       req.Name,
	}
	if err := h.store.UpdateSpotType(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     スポットタイプ削除
// @Description スポットタイプを削除します。
// @Tags        SpotType
// @Router      /v1/spot-types/{spotTypeId} [delete]
// @Security    bearerauth
// @Param       spotTypeId path string true "スポットタイプID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "スポットタイプが存在しない"
// @Failure     412 {object} util.ErrorResponse "スポット側で紐づいているため削除不可"
func (h *handler) DeleteSpotType(ctx *gin.Context) {
	in := &store.DeleteSpotTypeInput{
		SpotTypeID: ctx.Param("spotTypeId"),
	}
	if err := h.store.DeleteSpotType(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) multiGetSpotTypes(ctx context.Context, spotTypeIDs []string) (service.SpotTypes, error) {
	if len(spotTypeIDs) == 0 {
		return service.SpotTypes{}, nil
	}
	in := &store.MultiGetSpotTypesInput{
		SpotTypeIDs: spotTypeIDs,
	}
	spotTypes, err := h.store.MultiGetSpotTypes(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewSpotTypes(spotTypes), nil
}

func (h *handler) getSpotType(ctx context.Context, spotTypeID string) (*service.SpotType, error) {
	in := &store.GetSpotTypeInput{
		SpotTypeID: spotTypeID,
	}
	spotType, err := h.store.GetSpotType(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewSpotType(spotType), nil
}
