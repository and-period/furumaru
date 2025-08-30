package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
)

// @tag.name        Schedule
// @tag.description マルシェ開催スケジュール関連
func (h *handler) scheduleRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/schedules", h.authentication)

	r.GET("", h.ListSchedules)
	r.POST("", h.CreateSchedule)
	r.GET("/:scheduleId", h.filterAccessSchedule, h.GetSchedule)
	r.PATCH("/:scheduleId", h.filterAccessSchedule, h.UpdateSchedule)
	r.DELETE("/:scheduleId", h.filterAccessSchedule, h.DeleteSchedule)
	r.PATCH("/:scheduleId/approval", h.filterAccessSchedule, h.ApproveSchedule)
	r.PATCH("/:scheduleId/publish", h.filterAccessSchedule, h.PublishSchedule)
	r.GET("/:scheduleId/analytics", h.filterAccessSchedule, h.AnalyzeSchedule)
}

func (h *handler) filterAccessSchedule(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			schedule, err := h.getSchedule(ctx, util.GetParam(ctx, "scheduleId"))
			if err != nil {
				return false, err
			}
			return currentAdmin(ctx, schedule.CoordinatorID), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Next()
}

// @Summary     マルシェ開催スケジュール一覧取得
// @Description マルシェ開催スケジュールの一覧を取得します。ページネーションに対応しています。
// @Tags        Schedule
// @Router      /v1/schedules [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Produce     json
// @Success     200 {object} response.SchedulesResponse
func (h *handler) ListSchedules(ctx *gin.Context) {
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

	in := &store.ListSchedulesInput{
		ShopID: getShopID(ctx),
		Limit:  limit,
		Offset: offset,
	}
	schedules, total, err := h.store.ListSchedules(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(schedules) == 0 {
		res := &response.SchedulesResponse{
			Schedules:    []*response.Schedule{},
			Coordinators: []*response.Coordinator{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	coordinators, err := h.multiGetCoordinators(ctx, schedules.CoordinatorIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.SchedulesResponse{
		Schedules:    service.NewSchedules(schedules).Response(),
		Coordinators: coordinators.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     マルシェ開催スケジュール取得
// @Description 指定されたマルシェ開催スケジュールの詳細情報を取得します。
// @Tags        Schedule
// @Router      /v1/schedules/{scheduleId} [get]
// @Security    bearerauth
// @Param       scheduleId path string true "スケジュールID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} response.ScheduleResponse
// @Failure     403 {object} util.ErrorResponse "スケジュールの参照権限がない"
// @Failure     404 {object} util.ErrorResponse "スケジュールが存在しない"
func (h *handler) GetSchedule(ctx *gin.Context) {
	scheduleID := util.GetParam(ctx, "scheduleId")
	schedule, err := h.getSchedule(ctx, scheduleID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	coordinator, err := h.getCoordinator(ctx, schedule.CoordinatorID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := response.ScheduleResponse{
		Schedule:    schedule.Response(),
		Coordinator: coordinator.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     マルシェ開催スケジュール登録
// @Description 新しいマルシェ開催スケジュールを登録します。
// @Tags        Schedule
// @Router      /v1/schedules [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body request.CreateScheduleRequest true "スケジュール情報"
// @Produce     json
// @Success     200 {object} response.ScheduleResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateSchedule(ctx *gin.Context) {
	req := &request.CreateScheduleRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	if getAdminType(ctx).IsCoordinator() {
		if !currentAdmin(ctx, req.CoordinatorID) {
			h.forbidden(ctx, errors.New("handler: invalid coordinator id"))
			return
		}
	}

	coordinator, err := h.getCoordinator(ctx, req.CoordinatorID)
	if errors.Is(err, exception.ErrNotFound) {
		h.badRequest(ctx, err)
		return
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	shop, err := h.getShopByCoordinatorID(ctx, req.CoordinatorID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	in := &store.CreateScheduleInput{
		ShopID:          shop.ID,
		CoordinatorID:   req.CoordinatorID,
		Title:           req.Title,
		Description:     req.Description,
		ThumbnailURL:    req.ThumbnailURL,
		ImageURL:        req.ImageURL,
		OpeningVideoURL: req.OpeningVideoURL,
		Public:          req.Public,
		StartAt:         jst.ParseFromUnix(req.StartAt),
		EndAt:           jst.ParseFromUnix(req.EndAt),
	}
	schedule, err := h.store.CreateSchedule(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	sschedule := service.NewSchedule(schedule)

	res := &response.ScheduleResponse{
		Schedule:    sschedule.Response(),
		Coordinator: coordinator.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     マルシェ開催スケジュール更新
// @Description マルシェ開催スケジュールの情報を更新します。
// @Tags        Schedule
// @Router      /v1/schedules/{scheduleId} [patch]
// @Security    bearerauth
// @Param       scheduleId path string true "スケジュールID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body request.UpdateScheduleRequest true "スケジュール情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "スケジュールの更新権限がない"
// @Failure     404 {object} util.ErrorResponse "スケジュールが存在しない"
func (h *handler) UpdateSchedule(ctx *gin.Context) {
	req := &request.UpdateScheduleRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.UpdateScheduleInput{
		ScheduleID:      util.GetParam(ctx, "scheduleId"),
		Title:           req.Title,
		Description:     req.Description,
		ThumbnailURL:    req.ThumbnailURL,
		ImageURL:        req.ImageURL,
		OpeningVideoURL: req.OpeningVideoURL,
		StartAt:         jst.ParseFromUnix(req.StartAt),
		EndAt:           jst.ParseFromUnix(req.EndAt),
	}
	if err := h.store.UpdateSchedule(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     マルシェ開催スケジュール削除
// @Description マルシェ開催スケジュールを削除します。
// @Tags        Schedule
// @Router      /v1/schedules/{scheduleId} [delete]
// @Security    bearerauth
// @Param       scheduleId path string true "スケジュールID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     403 {object} util.ErrorResponse "スケジュールの削除権限がない"
// @Failure     404 {object} util.ErrorResponse "スケジュールが存在しない"
func (h *handler) DeleteSchedule(ctx *gin.Context) {
	in := &store.DeleteScheduleInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
	}
	if err := h.store.DeleteSchedule(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     マルシェ開催スケジュール承認
// @Description マルシェ開催スケジュールの承認状態を更新します。
// @Tags        Schedule
// @Router      /v1/schedules/{scheduleId}/approval [patch]
// @Security    bearerauth
// @Accept      json
// @Param       scheduleId path string true "スケジュールID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       request body request.ApproveScheduleRequest true "承認情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "スケジュールの承認権限がない"
// @Failure     404 {object} util.ErrorResponse "スケジュールが存在しない"
func (h *handler) ApproveSchedule(ctx *gin.Context) {
	req := &request.ApproveScheduleRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.ApproveScheduleInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
		AdminID:    getAdminID(ctx),
		Approved:   req.Approved,
	}
	if err := h.store.ApproveSchedule(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     マルシェ開催スケジュール公開
// @Description マルシェ開催スケジュールの公開状態を更新します。
// @Tags        Schedule
// @Router      /v1/schedules/{scheduleId}/publish [patch]
// @Security    bearerauth
// @Accept      json
// @Param       scheduleId path string true "スケジュールID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       request body request.PublishScheduleRequest true "公開設定情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     403 {object} util.ErrorResponse "スケジュールの公開権限がない"
// @Failure     404 {object} util.ErrorResponse "スケジュールが存在しない"
func (h *handler) PublishSchedule(ctx *gin.Context) {
	req := &request.PublishScheduleRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.PublishScheduleInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
		Public:     req.Public,
	}
	if err := h.store.PublishSchedule(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     マルシェ分析情報取得
// @Description 指定されたマルシェ開催スケジュールの視聴者分析データを取得します。集計期間と集計間隔を指定できます。
// @Tags        Schedule
// @Router      /v1/schedules/{scheduleId}/analytics [get]
// @Security    bearerauth
// @Param       scheduleId path string true "スケジュールID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Param       startAt query integer false "集計開始日時 (unixtime,未指定の場合はスケジュール開始時間)" example("1640962800")
// @Param       endAt query integer false "集計終了日時 (unixtime,未指定の場合はスケジュール終了時間)" example("1640962800")
// @Param       viewerLogInterval query string false "集計間隔 (未指定の場合は1分間隔)" example("minute")
// @Produce     json
// @Success     200 {object} response.AnalyzeScheduleResponse
// @Failure     403 {object} util.ErrorResponse "スケジュールの参照権限がない"
// @Failure     404 {object} util.ErrorResponse "スケジュールが存在しない"
func (h *handler) AnalyzeSchedule(ctx *gin.Context) {
	const defaultViewerLogInterval = service.BroadcastViewerLogIntervalMinute

	schedule, err := h.getSchedule(ctx, util.GetParam(ctx, "scheduleId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	startAtUnix, err := util.GetQueryInt64(ctx, "startAt", schedule.StartAt)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	endAtUnix, err := util.GetQueryInt64(ctx, "endAt", schedule.EndAt)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	viewerLogIntervalStr := util.GetQuery(ctx, "viewerLogInterval", string(defaultViewerLogInterval))

	startAt := jst.ParseFromUnix(startAtUnix)
	endAt := jst.ParseFromUnix(endAtUnix)
	viewerLogInterval := service.NewBroadcastViewerLogIntervalFromRequest(viewerLogIntervalStr)

	viewerLogsIn := &media.AggregateBroadcastViewerLogsInput{
		ScheduleID:   schedule.ID,
		Interval:     viewerLogInterval.MediaEntity(),
		CreatedAtGte: startAt,
		CreatedAtLt:  endAt,
	}
	viewerLogs, totalViewers, err := h.media.AggregateBroadcastViewerLogs(ctx, viewerLogsIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.AnalyzeScheduleResponse{
		ViewerLogs:   service.NewBroadcastViewerLogs(viewerLogInterval, startAt, endAt, viewerLogs).Response(),
		TotalViewers: totalViewers,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) getSchedule(ctx context.Context, scheduleID string) (*service.Schedule, error) {
	in := &store.GetScheduleInput{
		ScheduleID: scheduleID,
	}
	sschedule, err := h.store.GetSchedule(ctx, in)
	if err != nil {
		return nil, err
	}
	schedule := service.NewSchedule(sschedule)
	return schedule, nil
}
