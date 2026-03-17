package handler

import (
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/gin-gonic/gin"
)

// @tag.name        FeatureRequest
// @tag.description 要望リクエスト関連
func (h *handler) featureRequestRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/feature-requests", h.authentication)

	r.GET("", h.ListFeatureRequests)
	r.POST("", h.CreateFeatureRequest)
	r.GET("/:featureRequestId", h.GetFeatureRequest)
	r.PATCH("/:featureRequestId", h.UpdateFeatureRequest)
	r.DELETE("/:featureRequestId", h.DeleteFeatureRequest)
}

// @Summary     要望リクエスト一覧取得
// @Description 要望リクエストの一覧を取得します。管理者は全件、コーディネーターは自分の提出のみ。
// @Tags        FeatureRequest
// @Router      /v1/feature-requests [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Produce     json
// @Success     200 {object} types.FeatureRequestsResponse
func (h *handler) ListFeatureRequests(ctx *gin.Context) {
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

	in := &messenger.ListFeatureRequestsInput{
		Limit:  limit,
		Offset: offset,
	}

	// コーディネーターの場合は自分の提出のみ取得
	if err := filterAccess(ctx, &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			in.SubmittedBy = getAdminID(ctx)
			return true, nil
		},
	}); err != nil {
		h.httpError(ctx, err)
		return
	}

	featureRequests, total, err := h.messenger.ListFeatureRequests(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	sfeatureRequests := service.NewFeatureRequests(featureRequests)
	admins, err := h.multiGetAdmins(ctx, sfeatureRequests.SubmittedByIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	sfeatureRequests.Fill(admins.Map())

	res := &types.FeatureRequestsResponse{
		FeatureRequests: sfeatureRequests.Response(),
		Total:           total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     要望リクエスト登録
// @Description 新しい要望リクエストを登録します。管理者・コーディネーター両方可。
// @Tags        FeatureRequest
// @Router      /v1/feature-requests [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.CreateFeatureRequestRequest true "要望リクエスト情報"
// @Produce     json
// @Success     201 {object} types.FeatureRequestResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateFeatureRequest(ctx *gin.Context) {
	req := &types.CreateFeatureRequestRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &messenger.CreateFeatureRequestInput{
		Title:       req.Title,
		Description: req.Description,
		Category:    entity.FeatureRequestCategory(req.Category),
		Priority:    entity.FeatureRequestPriority(req.Priority),
		SubmittedBy: getAdminID(ctx),
	}
	sfeatureRequest, err := h.messenger.CreateFeatureRequest(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	featureRequest := service.NewFeatureRequest(sfeatureRequest)
	admin, err := h.getAdmin(ctx, sfeatureRequest.SubmittedBy)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	featureRequest.Fill(admin.Name())

	res := &types.FeatureRequestResponse{
		FeatureRequest: featureRequest.Response(),
	}
	ctx.JSON(http.StatusCreated, res)
}

// @Summary     要望リクエスト取得
// @Description 指定された要望リクエストの詳細情報を取得します。管理者は全件、コーディネーターは自分の提出のみ。
// @Tags        FeatureRequest
// @Router      /v1/feature-requests/{featureRequestId} [get]
// @Security    bearerauth
// @Param       featureRequestId path string true "要望リクエストID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.FeatureRequestResponse
// @Failure     403 {object} util.ErrorResponse "権限なし"
// @Failure     404 {object} util.ErrorResponse "要望リクエストが存在しない"
func (h *handler) GetFeatureRequest(ctx *gin.Context) {
	featureRequestID := util.GetParam(ctx, "featureRequestId")

	in := &messenger.GetFeatureRequestInput{
		FeatureRequestID: featureRequestID,
	}
	sfeatureRequest, err := h.messenger.GetFeatureRequest(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	featureRequest := service.NewFeatureRequest(sfeatureRequest)
	// コーディネーターは自分の提出のみ参照可（権限なしの場合も404を返し、存在情報を漏洩しない）
	if err := filterAccess(ctx, &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			if sfeatureRequest.SubmittedBy != getAdminID(ctx) {
				return false, fmt.Errorf("handler: feature request not found: %w", exception.ErrNotFound)
			}
			return true, nil
		},
	}); err != nil {
		h.httpError(ctx, err)
		return
	}

	admin, err := h.getAdmin(ctx, sfeatureRequest.SubmittedBy)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	featureRequest.Fill(admin.Name())

	res := &types.FeatureRequestResponse{
		FeatureRequest: featureRequest.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     要望リクエスト更新
// @Description 要望リクエストのステータスとコメントを更新します。管理者のみ。
// @Tags        FeatureRequest
// @Router      /v1/feature-requests/{featureRequestId} [patch]
// @Security    bearerauth
// @Param       featureRequestId path string true "要望リクエストID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.UpdateFeatureRequestRequest true "更新情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "権限なし（管理者のみ）"
// @Failure     404 {object} util.ErrorResponse "要望リクエストが存在しない"
func (h *handler) UpdateFeatureRequest(ctx *gin.Context) {
	// 管理者のみ許可
	if err := filterAccess(ctx, &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			return false, fmt.Errorf("handler: only administrator can update feature request: %w", exception.ErrForbidden)
		},
	}); err != nil {
		h.httpError(ctx, err)
		return
	}

	req := &types.UpdateFeatureRequestRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &messenger.UpdateFeatureRequestInput{
		FeatureRequestID: util.GetParam(ctx, "featureRequestId"),
		Status:           entity.FeatureRequestStatus(req.Status),
		Note:             req.Note,
	}
	if err := h.messenger.UpdateFeatureRequest(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     要望リクエスト削除
// @Description 要望リクエストを削除します。管理者は全件、コーディネーターは受付中（Waiting）の自分の提出のみ。
// @Tags        FeatureRequest
// @Router      /v1/feature-requests/{featureRequestId} [delete]
// @Security    bearerauth
// @Param       featureRequestId path string true "要望リクエストID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "権限なし"
// @Failure     404 {object} util.ErrorResponse "要望リクエストが存在しない"
func (h *handler) DeleteFeatureRequest(ctx *gin.Context) {
	featureRequestID := util.GetParam(ctx, "featureRequestId")

	// コーディネーターは受付中（Waiting=1）の自分の提出のみ削除可
	if err := filterAccess(ctx, &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			in := &messenger.GetFeatureRequestInput{
				FeatureRequestID: featureRequestID,
			}
			sfeatureRequest, err := h.messenger.GetFeatureRequest(ctx, in)
			if err != nil {
				return false, err
			}
			if sfeatureRequest.SubmittedBy != getAdminID(ctx) {
				return false, nil
			}
			if sfeatureRequest.Status != entity.FeatureRequestStatusWaiting {
				return false, fmt.Errorf("handler: can only delete waiting feature request: %w", exception.ErrForbidden)
			}
			return true, nil
		},
	}); err != nil {
		h.httpError(ctx, err)
		return
	}

	in := &messenger.DeleteFeatureRequestInput{
		FeatureRequestID: featureRequestID,
	}
	if err := h.messenger.DeleteFeatureRequest(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
