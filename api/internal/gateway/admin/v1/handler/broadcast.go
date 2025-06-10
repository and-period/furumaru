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
