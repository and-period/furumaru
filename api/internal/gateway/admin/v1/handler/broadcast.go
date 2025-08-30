package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/gin-gonic/gin"
)

// @tag.name        Broadcast
// @tag.description マルシェライブ配信関連
func (h *handler) broadcastRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/schedules/:scheduleId/broadcasts", h.authentication, h.filterAccessSchedule)

	r.GET("", h.GetBroadcast)
	r.POST("", h.UnpauseBroadcast)
	r.DELETE("", h.PauseBroadcast)
	r.POST("/archive-video", h.UploadBroadcastArchive)
	r.POST("/static-image", h.ActivateBroadcastStaticImage)
	r.DELETE("/static-image", h.DeactivateBroadcastStaticImage)
	r.POST("/rtmp", h.ActivateBroadcastRTMP)
	r.POST("/mp4", h.ActivateBroadcastMP4)
	r.POST("/youtube/auth", h.AuthYoutubeBroadcast)
}

func (h *handler) guestBroadcastRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/guests/schedules/-/broadcasts")

	r.GET("", h.GetGuestBroadcast)
	r.POST("/youtube", h.CreateYoutubeBroadcast)
	r.POST("/youtube/auth/complete", h.CallbackAuthYoutubeBroadcast)
}

// @Summary     マルシェライブ配信取得
// @Description 指定されたスケジュールのライブ配信情報を取得します。
// @Tags        Broadcast
// @Router      /v1/schedules/{scheduleId}/broadcasts [get]
// @Security    bearerauth
// @Param       scheduleId path string true "マルシェ開催スケジュールID" example("schedule-id")
// @Produce     json
// @Success     200 {object} response.BroadcastResponse
// @Failure     404 {object} util.ErrorResponse "マルシェライブ配信が存在しない"
func (h *handler) GetBroadcast(ctx *gin.Context) {
	in := &media.GetBroadcastByScheduleIDInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
	}
	broadcast, err := h.media.GetBroadcastByScheduleID(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.BroadcastResponse{
		Broadcast: service.NewBroadcast(broadcast).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     マルシェライブ配信一時停止
// @Description ライブ配信を一時停止します。
// @Tags        Broadcast
// @Router      /v1/schedules/{scheduleId}/broadcasts [delete]
// @Security    bearerauth
// @Param       scheduleId path string true "マルシェ開催スケジュールID" example("schedule-id")
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "マルシェライブ配信が存在しない"
func (h *handler) PauseBroadcast(ctx *gin.Context) {
	in := &media.PauseBroadcastInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
	}
	if err := h.media.PauseBroadcast(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     マルシェライブ配信一時停止解除
// @Description ライブ配信の一時停止を解除します。
// @Tags        Broadcast
// @Router      /v1/schedules/{scheduleId}/broadcasts [post]
// @Security    bearerauth
// @Param       scheduleId path string true "マルシェ開催スケジュールID" example("schedule-id")
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "マルシェライブ配信が存在しない"
func (h *handler) UnpauseBroadcast(ctx *gin.Context) {
	in := &media.UnpauseBroadcastInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
	}
	if err := h.media.UnpauseBroadcast(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     オンデマンド配信用の映像をアップロード
// @Description ライブ配信終了後にオンデマンド配信用の映像をアップロードします。
// @Tags        Broadcast
// @Router      /v1/schedules/{scheduleId}/broadcasts/archive-video [post]
// @Security    bearerauth
// @Param       scheduleId path string true "マルシェ開催スケジュールID" example("schedule-id")
// @Accept      json
// @Param       request body request.UpdateBroadcastArchiveRequest true "アーカイブURL"
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "マルシェライブ配信が存在しない"
// @Failure     412 {object} util.ErrorResponse "マルシェライブ配信が終了していない"
func (h *handler) UploadBroadcastArchive(ctx *gin.Context) {
	req := &request.UpdateBroadcastArchiveRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.UpdateBroadcastArchiveInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
		ArchiveURL: req.ArchiveURL,
	}
	if err := h.media.UpdateBroadcastArchive(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     ライブ配信中の入力をRTMPへ切り替え
// @Description ライブ配信の入力ソースをRTMPに切り替えます。
// @Tags        Broadcast
// @Router      /v1/schedules/{scheduleId}/broadcasts/rtmp [post]
// @Security    bearerauth
// @Param       scheduleId path string true "マルシェ開催スケジュールID" example("schedule-id")
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "マルシェライブ配信が存在しない"
// @Failure     412 {object} util.ErrorResponse "マルシェライブ配信中でない"
func (h *handler) ActivateBroadcastRTMP(ctx *gin.Context) {
	in := &media.ActivateBroadcastRTMPInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
	}
	if err := h.media.ActivateBroadcastRTMP(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     ライブ配信中の入力をMP4へ切り替え
// @Description ライブ配信の入力ソースをMP4ファイルに切り替えます。
// @Tags        Broadcast
// @Router      /v1/schedules/{scheduleId}/broadcasts/mp4 [post]
// @Security    bearerauth
// @Param       scheduleId path string true "マルシェ開催スケジュールID" example("schedule-id")
// @Accept      json
// @Param       request body request.ActivateBroadcastMP4Request true "MP4ファイルURL"
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "マルシェライブ配信が存在しない"
// @Failure     412 {object} util.ErrorResponse "マルシェライブ配信中でない"
func (h *handler) ActivateBroadcastMP4(ctx *gin.Context) {
	req := &request.ActivateBroadcastMP4Request{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.ActivateBroadcastMP4Input{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
		InputURL:   req.InputURL,
	}
	if err := h.media.ActivateBroadcastMP4(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     ライブ配信のふた絵を有効化
// @Description ライブ配信中にふた絵（静止画）を表示します。
// @Tags        Broadcast
// @Router      /v1/schedules/{scheduleId}/broadcasts/static-image [post]
// @Security    bearerauth
// @Param       scheduleId path string true "マルシェ開催スケジュールID" example("schedule-id")
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "マルシェライブ配信が存在しない"
// @Failure     412 {object} util.ErrorResponse "マルシェライブ配信中でない"
func (h *handler) ActivateBroadcastStaticImage(ctx *gin.Context) {
	in := &media.ActivateBroadcastStaticImageInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
	}
	if err := h.media.ActivateBroadcastStaticImage(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     ライブ配信のふた絵を無効化
// @Description ライブ配信中のふた絵（静止画）を無効化して通常配信に戻します。
// @Tags        Broadcast
// @Router      /v1/schedules/{scheduleId}/broadcasts/static-image [delete]
// @Security    bearerauth
// @Param       scheduleId path string true "マルシェ開催スケジュールID" example("schedule-id")
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "マルシェライブ配信が存在しない"
// @Failure     412 {object} util.ErrorResponse "マルシェライブ配信中でない"
func (h *handler) DeactivateBroadcastStaticImage(ctx *gin.Context) {
	in := &media.DeactivateBroadcastStaticImageInput{
		ScheduleID: util.GetParam(ctx, "scheduleId"),
	}
	if err := h.media.DeactivateBroadcastStaticImage(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     ゲスト用ライブ配信情報取得
// @Description セッションIDを使用してゲスト向けのライブ配信情報を取得します。
// @Tags        Guest
// @Tags        Broadcast
// @Router      /v1/guests/schedules/-/broadcasts [get]
// @Security    cookieauth
// @Produce     json
// @Success     200 {object} response.GuestBroadcastResponse
func (h *handler) GetGuestBroadcast(ctx *gin.Context) {
	sessionID, err := h.getSessionID(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	in := &media.GetBroadcastAuthInput{
		SessionID: sessionID,
	}
	auth, err := h.media.GetBroadcastAuth(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	schedule, err := h.getSchedule(ctx, auth.ScheduleID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	shop, err := h.getShop(ctx, schedule.ShopID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	coordinator, err := h.getCoordinator(ctx, schedule.CoordinatorID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.GuestBroadcastResponse{
		Broadcast: service.NewGuestBroadcast(schedule, shop, coordinator).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     マルシェライブ配信のYoutube連携を認証
// @Description Youtube配信用の認証URLを取得します。
// @Tags        Broadcast
// @Router      /v1/schedules/{scheduleId}/broadcasts/youtube/auth [post]
// @Security    bearerauth
// @Param       scheduleId path string true "マルシェ開催スケジュールID" example("schedule-id")
// @Accept      json
// @Param       request body request.AuthYoutubeBroadcastRequest true "Youtubeハンドル"
// @Produce     json
// @Success     200 {object} response.AuthYoutubeBroadcastResponse
// @Failure     404 {object} util.ErrorResponse "マルシェライブ配信が存在しない"
// @Failure     412 {object} util.ErrorResponse "マルシェライブ配信前でない"
func (h *handler) AuthYoutubeBroadcast(ctx *gin.Context) {
	req := &request.AuthYoutubeBroadcastRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.AuthYoutubeBroadcastInput{
		ScheduleID:    util.GetParam(ctx, "scheduleId"),
		YoutubeHandle: req.YoutubeHandle,
	}
	authURL, err := h.media.AuthYoutubeBroadcast(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.AuthYoutubeBroadcastResponse{
		URL: authURL,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     マルシェライブ配信のYoutube認証後処理
// @Description Youtube認証完了後のコールバック処理を行います。
// @Tags        Guest
// @Tags        Broadcast
// @Router      /v1/guests/schedules/-/broadcasts/youtube/auth/complete [post]
// @Accept      json
// @Param       request body request.CallbackAuthYoutubeBroadcastRequest true "認証コールバック"
// @Produce     json
// @Success     200 {object} response.GuestBroadcastResponse
// @Failure     401 {object} util.ErrorResponse "Youtube APIの認証エラー"
// @Failure     403 {object} util.ErrorResponse "Youtube APIの権限エラー"
// @Failure     404 {object} util.ErrorResponse "マルシェライブ配信が存在しない"
func (h *handler) CallbackAuthYoutubeBroadcast(ctx *gin.Context) {
	req := &request.CallbackAuthYoutubeBroadcastRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &media.AuthYoutubeBroadcastEventInput{
		State:    req.State,
		AuthCode: req.AuthCode,
	}
	auth, err := h.media.AuthYoutubeBroadcastEvent(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	h.setSessionID(ctx, auth.SessionID)
	schedule, err := h.getSchedule(ctx, auth.ScheduleID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	shop, err := h.getShop(ctx, schedule.ShopID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	coordinator, err := h.getCoordinator(ctx, schedule.CoordinatorID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.GuestBroadcastResponse{
		Broadcast: service.NewGuestBroadcast(schedule, shop, coordinator).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     マルシェライブ配信のYoutube連携
// @Description Youtube側でライブ配信を作成します。
// @Tags        Guest
// @Tags        Broadcast
// @Router      /v1/guests/schedules/-/broadcasts/youtube [post]
// @Security    cookieauth
// @Accept      json
// @Param       request body request.CreateYoutubeBroadcastRequest true "Youtube配信設定"
// @Produce     json
// @Success     204
// @Failure     401 {object} util.ErrorResponse "Youtube APIの認証エラー"
// @Failure     403 {object} util.ErrorResponse "Youtube APIの権限エラー"
// @Failure     404 {object} util.ErrorResponse "マルシェライブ配信が存在しない"
// @Failure     412 {object} util.ErrorResponse "マルシェライブ配信前でない"
func (h *handler) CreateYoutubeBroadcast(ctx *gin.Context) {
	req := &request.CreateYoutubeBroadcastRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	sessionID, err := h.getSessionID(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}
	in := &media.CreateYoutubeBroadcastInput{
		SessionID:   sessionID,
		Title:       req.Title,
		Description: req.Description,
		Public:      req.Public,
	}
	if err := h.media.CreateYoutubeBroadcast(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
